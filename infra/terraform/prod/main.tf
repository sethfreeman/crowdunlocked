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
    key            = "prod/terraform.tfstate"
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
  name     = "crowdunlocked-prod"
  role_arn = aws_iam_role.eks_cluster.arn
  version  = "1.34"

  vpc_config {
    subnet_ids              = aws_subnet.private[*].id
    endpoint_private_access = true
    endpoint_public_access  = false
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
  cidr_block           = "10.1.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "crowdunlocked-prod"
  }
}

resource "aws_subnet" "private" {
  count             = 3
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.1.${count.index + 1}.0/24"
  availability_zone = data.aws_availability_zones.available.names[count.index]

  tags = {
    Name                              = "crowdunlocked-prod-private-${count.index + 1}"
    "kubernetes.io/role/internal-elb" = "1"
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}

# IAM Roles
resource "aws_iam_role" "eks_cluster" {
  name = "crowdunlocked-prod-eks-cluster"

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
  name = "crowdunlocked-prod-eks-node"

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

# DynamoDB Tables (Production)
resource "aws_dynamodb_table" "bookings" {
  name           = "bookings-prod"
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
    Environment = "prod"
    Service     = "bookings"
  }
}

resource "aws_dynamodb_table" "venues" {
  name         = "venues-prod"
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
    Environment = "prod"
    Service     = "bookings"
  }
}

# Read mgmt state for ACM certificate ARN
data "terraform_remote_state" "mgmt" {
  backend = "s3"
  config = {
    bucket = "crowdunlocked-terraform-state"
    key    = "mgmt/terraform.tfstate"
    region = "us-east-1"
  }
}

locals {
  acm_certificate_arn = data.terraform_remote_state.mgmt.outputs.prod_acm_certificate_arn
  has_certificate     = local.acm_certificate_arn != null
}

# CloudFront Distribution (only created after mgmt account deploys certificates)
resource "aws_cloudfront_distribution" "main" {
  count = local.has_certificate ? 1 : 0

  enabled             = true
  is_ipv6_enabled     = true
  comment             = "Crowd Unlocked Production CDN"
  default_root_object = "index.html"
  price_class         = "PriceClass_All"

  origin {
    domain_name = aws_lb.main.dns_name
    origin_id   = "eks-alb"

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  default_cache_behavior {
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "eks-alb"

    forwarded_values {
      query_string = true
      headers      = ["Host", "Authorization"]

      cookies {
        forward = "all"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = local.acm_certificate_arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }
}

resource "aws_lb" "main" {
  name               = "crowdunlocked-prod-alb"
  internal           = false
  load_balancer_type = "application"
  subnets            = aws_subnet.private[*].id

  enable_deletion_protection = true
  enable_http2               = true
  enable_cross_zone_load_balancing = true
}

# CloudWatch Alarms (Production)
resource "aws_cloudwatch_metric_alarm" "high_error_rate" {
  for_each = toset(["bookings", "releases", "publicity", "social", "money"])

  alarm_name          = "${each.key}-prod-high-error-rate"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "5XXError"
  namespace           = "AWS/ApiGateway"
  period              = 300
  statistic           = "Sum"
  threshold           = 5
  alarm_description   = "High error rate for ${each.key} service"
  treat_missing_data  = "notBreaching"
}

# GitHub OIDC for CI/CD
module "github_oidc" {
  source = "../modules/github-oidc"

  environment = "prod"
  github_org  = var.github_org
  github_repo = var.github_repo
}
