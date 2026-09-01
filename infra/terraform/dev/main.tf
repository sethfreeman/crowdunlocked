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

data "aws_caller_identity" "current" {}

# =============================================================================
# Local Variables
# =============================================================================

locals {
  aws_account_id = data.aws_caller_identity.current.account_id
}

# =============================================================================
# DynamoDB Tables
#
# These back the Next.js API routes (venues, bookings). Retained after the
# EKS -> Vercel migration. Pay-per-request billing scales to ~zero when idle.
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
# Vercel OIDC Integration
#
# Grants Vercel serverless functions keyless (OIDC) access to the venues and
# bookings DynamoDB tables. This is the backend auth path after migrating off
# EKS. Populate vercel_project_ids after creating the Vercel project.
# =============================================================================

module "vercel_oidc" {
  source = "../modules/vercel-oidc"

  project_name   = "crowdunlocked"
  environment    = "dev"
  aws_region     = var.aws_region
  aws_account_id = local.aws_account_id

  # Vercel identity (team issuer mode). sub uses the team slug + project name.
  vercel_team_slug    = "seth-freemans-projects"
  vercel_project_name = "crowdunlocked"
  vercel_environments = ["production", "preview"]

  venues_table_arn   = aws_dynamodb_table.venues.arn
  bookings_table_arn = aws_dynamodb_table.bookings.arn
}
