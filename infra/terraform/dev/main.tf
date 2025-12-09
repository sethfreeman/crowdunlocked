# =============================================================================
# Terraform Configuration
# =============================================================================

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
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-state-lock"
  }
}

provider "aws" {
  region = var.aws_region
}

# =============================================================================
# Data Sources
# =============================================================================

data "aws_availability_zones" "available" {
  state = "available"
}

data "terraform_remote_state" "mgmt" {
  backend = "s3"
  config = {
    bucket = "crowdunlocked-terraform-state"
    key    = "mgmt/terraform.tfstate"
    region = "us-east-1"
  }
}

# =============================================================================
# Local Variables
# =============================================================================

locals {
  services            = ["bookings", "releases", "publicity", "social", "money"]
  acm_certificate_arn = data.terraform_remote_state.mgmt.outputs.dev_acm_certificate_arn
  has_certificate     = local.acm_certificate_arn != null
}

# =============================================================================
# VPC and Networking
# =============================================================================

resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "crowdunlocked-dev"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "crowdunlocked-dev"
  }
}

# Public Subnets
resource "aws_subnet" "public" {
  count                   = 3
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.${count.index + 101}.0/24"
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name                                        = "crowdunlocked-dev-public-${count.index + 1}"
    "kubernetes.io/role/elb"                    = "1"
    "kubernetes.io/cluster/crowdunlocked-dev"   = "shared"
  }
}

# Private Subnets
resource "aws_subnet" "private" {
  count             = 3
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.${count.index + 1}.0/24"
  availability_zone = data.aws_availability_zones.available.names[count.index]

  tags = {
    Name                                        = "crowdunlocked-dev-private-${count.index + 1}"
    "kubernetes.io/role/internal-elb"           = "1"
    "kubernetes.io/cluster/crowdunlocked-dev"   = "shared"
  }
}

# Elastic IPs for NAT Gateways
resource "aws_eip" "nat" {
  count  = 3
  domain = "vpc"

  tags = {
    Name = "crowdunlocked-dev-nat-${count.index + 1}"
  }

  depends_on = [aws_internet_gateway.main]
}

# NAT Gateways
resource "aws_nat_gateway" "main" {
  count         = 3
  allocation_id = aws_eip.nat[count.index].id
  subnet_id     = aws_subnet.public[count.index].id

  tags = {
    Name = "crowdunlocked-dev-${count.index + 1}"
  }

  depends_on = [aws_internet_gateway.main]
}

# Public Route Table
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "crowdunlocked-dev-public"
  }
}

resource "aws_route_table_association" "public" {
  count          = 3
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# Private Route Tables
resource "aws_route_table" "private" {
  count  = 3
  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.main[count.index].id
  }

  tags = {
    Name = "crowdunlocked-dev-private-${count.index + 1}"
  }
}

resource "aws_route_table_association" "private" {
  count          = 3
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[count.index].id
}

# =============================================================================
# IAM Roles and Policies
# =============================================================================

# EKS Cluster Role
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

# EKS Node Role
resource "aws_iam_role" "eks_node" {
  name = "crowdunlocked-dev-eks-node"

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

resource "aws_iam_role_policy_attachment" "eks_node_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.eks_node.name
}

resource "aws_iam_role_policy_attachment" "eks_cni_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.eks_node.name
}

resource "aws_iam_role_policy_attachment" "eks_container_registry_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.eks_node.name
}

# =============================================================================
# EKS Cluster
# =============================================================================

resource "aws_eks_cluster" "main" {
  name     = "crowdunlocked-dev"
  role_arn = aws_iam_role.eks_cluster.arn
  version  = "1.34"

  vpc_config {
    subnet_ids              = aws_subnet.private[*].id
    endpoint_private_access = true
    endpoint_public_access  = true
  }

  access_config {
    authentication_mode = "API_AND_CONFIG_MAP"
  }

  compute_config {
    enabled                      = true
    node_pools                   = ["general-purpose"]
    node_role_arn                = aws_iam_role.eks_node.arn
  }

  bootstrap_self_managed_addons = false

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
    aws_iam_role_policy_attachment.eks_cluster_policy,
    aws_nat_gateway.main
  ]
}

# =============================================================================
# DynamoDB Tables
# =============================================================================

resource "aws_dynamodb_table" "bookings" {
  name             = "bookings-dev"
  billing_mode     = "PAY_PER_REQUEST"
  hash_key         = "id"
  stream_enabled   = true
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

  attribute {
    name = "geohash"
    type = "S"
  }

  attribute {
    name = "geohash_sort"
    type = "S"
  }

  attribute {
    name = "city_state"
    type = "S"
  }

  attribute {
    name = "name"
    type = "S"
  }

  attribute {
    name = "venue_type"
    type = "S"
  }

  attribute {
    name = "rating_id"
    type = "S"
  }

  attribute {
    name = "external_source_id"
    type = "S"
  }

  global_secondary_index {
    name            = "GeohashIndex"
    hash_key        = "geohash"
    range_key       = "geohash_sort"
    projection_type = "ALL"
  }

  global_secondary_index {
    name            = "CityIndex"
    hash_key        = "city_state"
    range_key       = "name"
    projection_type = "ALL"
  }

  global_secondary_index {
    name            = "VenueTypeIndex"
    hash_key        = "venue_type"
    range_key       = "rating_id"
    projection_type = "ALL"
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

# =============================================================================
# ECR Repositories
# =============================================================================

resource "aws_ecr_repository" "services" {
  for_each = toset(local.services)
  name     = "crowdunlocked/${each.key}"

  image_scanning_configuration {
    scan_on_push = true
  }

  image_tag_mutability = "MUTABLE"

  tags = {
    Name        = "crowdunlocked-${each.key}"
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

resource "aws_ecr_repository" "web" {
  name = "crowdunlocked/web"

  image_scanning_configuration {
    scan_on_push = true
  }

  image_tag_mutability = "MUTABLE"

  tags = {
    Name        = "crowdunlocked-web"
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

resource "aws_ecr_lifecycle_policy" "services" {
  for_each   = toset(local.services)
  repository = aws_ecr_repository.services[each.key].name

  policy = jsonencode({
    rules = [{
      rulePriority = 1
      description  = "Keep last 10 images"
      selection = {
        tagStatus   = "any"
        countType   = "imageCountMoreThan"
        countNumber = 10
      }
      action = {
        type = "expire"
      }
    }]
  })
}

resource "aws_ecr_lifecycle_policy" "web" {
  repository = aws_ecr_repository.web.name

  policy = jsonencode({
    rules = [{
      rulePriority = 1
      description  = "Keep last 10 images"
      selection = {
        tagStatus   = "any"
        countType   = "imageCountMoreThan"
        countNumber = 10
      }
      action = {
        type = "expire"
      }
    }]
  })
}

# =============================================================================
# CloudWatch
# =============================================================================

resource "aws_cloudwatch_log_group" "services" {
  for_each = toset(local.services)

  name              = "/aws/eks/crowdunlocked-dev/${each.key}"
  retention_in_days = 7

  tags = {
    Environment = "dev"
    Service     = each.key
  }
}

resource "aws_cloudwatch_metric_alarm" "high_error_rate" {
  for_each = toset(local.services)

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

  tags = {
    Environment = "dev"
    Service     = each.key
  }
}
