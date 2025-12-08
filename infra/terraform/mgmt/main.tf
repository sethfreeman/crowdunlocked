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
# This terraform configuration uses data sources to reference existing resources

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

# NOTE: GitHub OIDC provider and role created manually via scripts/create-oidc-role.sh
# Not managed by terraform to avoid chicken-and-egg problem

# Terraform state bucket (created manually during bootstrap)
data "aws_s3_bucket" "terraform_state" {
  bucket = "crowdunlocked-terraform-state"
}

# Allow dev and prod OIDC roles to access state bucket
resource "aws_s3_bucket_policy" "terraform_state" {
  bucket = data.aws_s3_bucket.terraform_state.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowDevProdOIDCAccess"
        Effect = "Allow"
        Principal = {
          AWS = [
            "arn:aws:iam::${var.dev_account_id}:role/github-actions-dev",
            "arn:aws:iam::${var.dev_account_id}:role/OrganizationAccountAccessRole",
            "arn:aws:iam::${var.prod_account_id}:role/github-actions-prod",
            "arn:aws:iam::987470856210:role/github-actions-dev"
          ]
        }
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = [
          data.aws_s3_bucket.terraform_state.arn,
          "${data.aws_s3_bucket.terraform_state.arn}/*"
        ]
      }
    ]
  })
}
