terraform {
  required_version = ">= 1.10"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.25"
    }
  }
  backend "s3" {
    bucket         = "crowdunlocked-terraform-state"
    key            = "dev/terraform.tfstate"
    region         = "us-east-1"  # State bucket stays in us-east-1
    encrypt        = true
    dynamodb_table = "terraform-state-lock"
  }
}

provider "aws" {
  region = var.aws_region
}

# EKS Auto Mode Cluster
resource "aws_eks_cluster" "main" {
  name     = "crowdunlocked-dev"
  role_arn = aws_iam_role.eks_cluster.arn
  version  = "1.34"

  vpc_config {
    subnet_ids              = aws_subnet.private[*].id
    endpoint_private_access = true
    endpoint_public_access  = true
  }

  compute_config {
    enabled       = true
    node_pools    = ["general-purpose"]
    node_role_arn = aws_iam_role.eks_node.arn
  }

  kubernetes_network_config {
    elastic_load_balancing {
      enabled = true
    }
  }

  storage_config {
    block_storage {
      enabled = true
    }
  }

  depends_on = [
    aws_iam_role_policy_attachment.eks_cluster_policy
  ]
}

# VPC
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "crowdunlocked-dev"
  }
}

resource "aws_subnet" "private" {
  count             = 3
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.${count.index + 1}.0/24"
  availability_zone = data.aws_availability_zones.available.names[count.index]

  tags = {
    Name                              = "crowdunlocked-dev-private-${count.index + 1}"
    "kubernetes.io/role/internal-elb" = "1"
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}

# IAM Roles
resource "aws_iam_role" "eks_cluster" {
  name = "crowdunlocked-dev-eks-cluster"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "eks.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "eks_cluster_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.eks_cluster.name
}

resource "aws_iam_role" "eks_node" {
  name = "crowdunlocked-dev-eks-node"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "eks_node_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.eks_node.name
}

resource "aws_iam_role_policy_attachment" "eks_cni_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.eks_node.name
}

# Check if mgmt state exists by trying to read it
# This will return null if the state doesn't exist yet
data "external" "mgmt_state_check" {
  program = ["sh", "-c", <<-EOT
    if aws s3api head-object --bucket crowdunlocked-terraform-state --key mgmt/terraform.tfstate --region us-east-1 2>/dev/null; then
      echo '{"exists":"true"}'
    else
      echo '{"exists":"false"}'
    fi
  EOT
  ]
}

# Only read mgmt state if it exists
data "terraform_remote_state" "mgmt" {
  count   = data.external.mgmt_state_check.result.exists == "true" ? 1 : 0
  backend = "s3"
  config = {
    bucket = "crowdunlocked-terraform-state"
    key    = "mgmt/terraform.tfstate"
    region = "us-east-1"
  }
}

locals {
  acm_certificate_arn = length(data.terraform_remote_state.mgmt) > 0 ? data.terraform_remote_state.mgmt[0].outputs.dev_acm_certificate_arn : null
  has_certificate     = local.acm_certificate_arn != null
}

# DynamoDB Tables
resource "aws_dynamodb_table" "bookings" {
  name           = "bookings-dev"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "id"
  stream_enabled = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  attribute {
    name = "id"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }

  tags = {
    Environment = "dev"
    Service     = "bookings"
  }
}

resource "aws_dynamodb_table" "venues" {
  name         = "venues-dev"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  # GSI1: Geohash index for spatial queries
  attribute {
    name = "geohash"
    type = "S"
  }

  attribute {
    name = "geohash_sort"
    type = "S"
  }

  global_secondary_index {
    name            = "GeohashIndex"
    hash_key        = "geohash"
    range_key       = "geohash_sort"
    projection_type = "ALL"
  }

  # GSI2: City index for location-based queries
  attribute {
    name = "city_state"
    type = "S"
  }

  attribute {
    name = "name"
    type = "S"
  }

  global_secondary_index {
    name            = "CityIndex"
    hash_key        = "city_state"
    range_key       = "name"
    projection_type = "ALL"
  }

  # GSI3: Venue type index for filtering
  attribute {
    name = "venue_type"
    type = "S"
  }

  attribute {
    name = "rating_id"
    type = "S"
  }

  global_secondary_index {
    name            = "VenueTypeIndex"
    hash_key        = "venue_type"
    range_key       = "rating_id"
    projection_type = "ALL"
  }

  # GSI4: External ID index for deduplication
  attribute {
    name = "external_source_id"
    type = "S"
  }

  global_secondary_index {
    name            = "ExternalIdIndex"
    hash_key        = "external_source_id"
    range_key       = "id"
    projection_type = "ALL"
  }

  point_in_time_recovery {
    enabled = true
  }

  tags = {
    Environment = "dev"
    Service     = "bookings"
  }
}

resource "aws_dynamodb_table" "releases" {
  name         = "releases-dev"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }

  tags = {
    Environment = "dev"
    Service     = "releases"
  }
}

resource "aws_dynamodb_table" "publicity" {
  name         = "publicity-dev"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }

  tags = {
    Environment = "dev"
    Service     = "publicity"
  }
}

resource "aws_dynamodb_table" "social" {
  name         = "social-dev"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }

  tags = {
    Environment = "dev"
    Service     = "social"
  }
}

resource "aws_dynamodb_table" "money" {
  name         = "money-dev"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }

  tags = {
    Environment = "dev"
    Service     = "money"
  }
}

# CloudWatch Log Groups
resource "aws_cloudwatch_log_group" "services" {
  for_each = toset(["bookings", "releases", "publicity", "social", "money"])

  name              = "/aws/eks/crowdunlocked-dev/${each.key}"
  retention_in_days = 7
}

# CloudWatch Alarms
resource "aws_cloudwatch_metric_alarm" "high_error_rate" {
  for_each = toset(["bookings", "releases", "publicity", "social", "money"])

  alarm_name          = "${each.key}-dev-high-error-rate"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "5XXError"
  namespace           = "AWS/ApiGateway"
  period              = 300
  statistic           = "Sum"
  threshold           = 10
  alarm_description   = "High error rate for ${each.key} service"
  treat_missing_data  = "notBreaching"
}
