# ===================================================================
# Vercel OIDC Provider Configuration
#
# Team issuer mode: the issuer is scoped to the team slug.
#   iss = https://oidc.vercel.com/<team_slug>
# The OIDC provider's condition keys are prefixed with the issuer host+path
# WITHOUT the scheme, i.e. "oidc.vercel.com/<team_slug>:aud" / ":sub".
# ===================================================================

locals {
  vercel_issuer_host = "oidc.vercel.com/${var.vercel_team_slug}"
  vercel_issuer_url  = "https://${local.vercel_issuer_host}"
  vercel_audience    = "https://vercel.com/${var.vercel_team_slug}"

  # sub = owner:<team_slug>:project:<project_name>:environment:<env>
  vercel_subjects = [
    for env in var.vercel_environments :
    "owner:${var.vercel_team_slug}:project:${var.vercel_project_name}:environment:${env}"
  ]
}

# Create OIDC provider for Vercel (team issuer mode)
resource "aws_iam_openid_connect_provider" "vercel" {
  url = local.vercel_issuer_url

  client_id_list = [
    local.vercel_audience,
  ]

  # Vercel's OIDC thumbprint. AWS ignores this for recognized shared providers,
  # but the argument is still required.
  thumbprint_list = [
    "6938fd4d98bab03faadb97b34396831e3780aea1",
  ]

  tags = {
    Name        = "${var.project_name}-${var.environment}-vercel-oidc"
    Environment = var.environment
    ManagedBy   = "terraform"
    Service     = "vercel"
  }
}

# ===================================================================
# IAM Role for Vercel Deployments
# ===================================================================

# Trust policy for Vercel OIDC
data "aws_iam_policy_document" "vercel_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.vercel.arn]
    }

    actions = ["sts:AssumeRoleWithWebIdentity"]

    condition {
      test     = "StringEquals"
      variable = "${local.vercel_issuer_host}:aud"
      values   = [local.vercel_audience]
    }

    condition {
      test     = "StringEquals"
      variable = "${local.vercel_issuer_host}:sub"
      values   = local.vercel_subjects
    }
  }
}

# Create IAM role for Vercel
resource "aws_iam_role" "vercel" {
  name               = "${var.project_name}-${var.environment}-vercel-role"
  assume_role_policy = data.aws_iam_policy_document.vercel_assume_role.json

  tags = {
    Name        = "${var.project_name}-${var.environment}-vercel-role"
    Environment = var.environment
    ManagedBy   = "terraform"
    Service     = "vercel"
  }
}

# ===================================================================
# IAM Policies for DynamoDB Access
# ===================================================================

# Policy for DynamoDB access
data "aws_iam_policy_document" "vercel_dynamodb" {
  # Venues table access
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem",
      "dynamodb:Query",
      "dynamodb:Scan",
    ]
    resources = [
      var.venues_table_arn,
      "${var.venues_table_arn}/index/*",
    ]
  }

  # Bookings table access
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem",
      "dynamodb:Query",
      "dynamodb:Scan",
    ]
    resources = [
      var.bookings_table_arn,
      "${var.bookings_table_arn}/index/*",
    ]
  }
}

# Create and attach DynamoDB policy
resource "aws_iam_policy" "vercel_dynamodb" {
  name        = "${var.project_name}-${var.environment}-vercel-dynamodb"
  description = "DynamoDB access policy for Vercel deployments"
  policy      = data.aws_iam_policy_document.vercel_dynamodb.json

  tags = {
    Name        = "${var.project_name}-${var.environment}-vercel-dynamodb"
    Environment = var.environment
    ManagedBy   = "terraform"
    Service     = "vercel"
  }
}

resource "aws_iam_role_policy_attachment" "vercel_dynamodb" {
  role       = aws_iam_role.vercel.name
  policy_arn = aws_iam_policy.vercel_dynamodb.arn
}

# ===================================================================
# Optional: CloudWatch Logs Access
# ===================================================================

data "aws_iam_policy_document" "vercel_logs" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = [
      "arn:aws:logs:${var.aws_region}:${var.aws_account_id}:log-group:/vercel/${var.project_name}/*",
    ]
  }
}

resource "aws_iam_policy" "vercel_logs" {
  name        = "${var.project_name}-${var.environment}-vercel-logs"
  description = "CloudWatch Logs access policy for Vercel deployments"
  policy      = data.aws_iam_policy_document.vercel_logs.json

  tags = {
    Name        = "${var.project_name}-${var.environment}-vercel-logs"
    Environment = var.environment
    ManagedBy   = "terraform"
    Service     = "vercel"
  }
}

resource "aws_iam_role_policy_attachment" "vercel_logs" {
  role       = aws_iam_role.vercel.name
  policy_arn = aws_iam_policy.vercel_logs.arn
}