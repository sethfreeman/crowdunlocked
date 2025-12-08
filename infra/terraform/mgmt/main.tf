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
    key            = "mgmt/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-state-lock"
  }
}

provider "aws" {
  region = var.aws_region
}

# NOTE: AWS Organization, accounts, domains, and certificates are created manually
# See docs/AWS_ORGANIZATION_SETUP.md for setup instructions
#
# This terraform configuration imports and manages existing resources

# Data sources for manually created resources
data "aws_route53_zone" "prod" {
  name = var.prod_domain_name
}

data "aws_route53_zone" "dev" {
  name = var.dev_domain_name
}

data "aws_acm_certificate" "prod" {
  domain   = var.prod_domain_name
  statuses = ["ISSUED"]
  most_recent = true
}

data "aws_acm_certificate" "dev" {
  domain   = var.dev_domain_name
  statuses = ["ISSUED"]
  most_recent = true
}

# GitHub OIDC for CI/CD
module "github_oidc" {
  source = "../modules/github-oidc"

  environment = "mgmt"
  github_org  = var.github_org
  github_repo = var.github_repo
}
