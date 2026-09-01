variable "project_name" {
  description = "Name of the AWS-side project (used for resource naming)"
  type        = string
  default     = "crowdunlocked"
}

variable "environment" {
  description = "Environment name (dev, prod) - used for AWS resource naming"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "aws_account_id" {
  description = "AWS account ID"
  type        = string
}

# ---------------------------------------------------------------------------
# Vercel OIDC identity
#
# Vercel OIDC tokens carry these claims (team issuer mode):
#   iss = https://oidc.vercel.com/<vercel_team_slug>
#   aud = https://vercel.com/<vercel_team_slug>
#   sub = owner:<vercel_team_slug>:project:<vercel_project_name>:environment:<env>
#
# NOTE: sub uses the team SLUG and project NAME (not the team_/prj_ IDs).
# ---------------------------------------------------------------------------

variable "vercel_team_slug" {
  description = "Vercel team slug (e.g. seth-freemans-projects). Appears in the dashboard URL and in the OIDC iss/aud/sub claims."
  type        = string
}

variable "vercel_project_name" {
  description = "Vercel project name (e.g. crowdunlocked). Used in the OIDC sub claim."
  type        = string
}

variable "vercel_environments" {
  description = "Vercel deployment environments allowed to assume the role."
  type        = list(string)
  default     = ["production", "preview"]
}

variable "venues_table_arn" {
  description = "ARN of the venues DynamoDB table"
  type        = string
}

variable "bookings_table_arn" {
  description = "ARN of the bookings DynamoDB table"
  type        = string
}
