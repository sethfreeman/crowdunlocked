terraform {
  required_version = ">= 1.10"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.25"
    }
  }
}

# GitHub OIDC Provider
resource "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = [
    "sts.amazonaws.com",
  ]

  thumbprint_list = [
    "6938fd4d98bab03faadb97b34396831e3780aea1",
    "1c58a3a8518e8759bf075b76b750d4f2df264fcd",
  ]

  tags = {
    Name        = "github-actions-oidc"
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

# IAM Role for GitHub Actions
resource "aws_iam_role" "github_actions" {
  name = "github-actions-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = aws_iam_openid_connect_provider.github.arn
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
          }
          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:${var.github_org}/${var.github_repo}:*"
          }
        }
      }
    ]
  })

  tags = {
    Name        = "github-actions-${var.environment}"
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

# Policy for terraform operations
resource "aws_iam_role_policy" "github_actions_terraform" {
  name = "terraform-permissions"
  role = aws_iam_role.github_actions.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          # S3 for terraform state
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket",
          # DynamoDB for state locking
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:DeleteItem",
          # EKS
          "eks:*",
          # EC2 for VPC/networking
          "ec2:*",
          # IAM for roles/policies
          "iam:*",
          # DynamoDB for application tables
          "dynamodb:*",
          # CloudWatch Logs
          "logs:*",
          # CloudFront
          "cloudfront:*",
          # ACM
          "acm:*",
          # Route53
          "route53:*",
          # Secrets Manager
          "secretsmanager:*",
          # KMS
          "kms:*",
          # Auto Scaling
          "autoscaling:*",
          # Elastic Load Balancing
          "elasticloadbalancing:*",
        ]
        Resource = "*"
      }
    ]
  })
}
