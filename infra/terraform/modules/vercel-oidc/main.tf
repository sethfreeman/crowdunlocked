# ===================================================================
# Vercel OIDC Provider Configuration
# ===================================================================

# Create OIDC provider for Vercel
resource "aws_iam_openid_connect_provider" "vercel" {
  url = "https://oidc.vercel.com"

  client_id_list = [
    "https://vercel.com",
  ]

  # Vercel's OIDC thumbprint
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
      variable = "oidc.vercel.com:aud"
      values   = ["https://vercel.com"]
    }

    condition {
      test     = "StringLike"
      variable = "oidc.vercel.com:sub"
      # Format: team_<team_id>:project_<project_id>:environment_<environment>
      values = var.vercel_project_ids
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