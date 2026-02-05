variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "crowdunlocked"
}

variable "environment" {
  description = "Environment name (dev, prod)"
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

variable "vercel_project_ids" {
  description = "List of Vercel project identifiers for OIDC subject conditions"
  type        = list(string)
  # Format: "team_<team_id>:project_<project_id>:environment_<environment>"
  # Example: ["team_abc123:project_xyz789:environment_production"]
}

variable "venues_table_arn" {
  description = "ARN of the venues DynamoDB table"
  type        = string
}

variable "bookings_table_arn" {
  description = "ARN of the bookings DynamoDB table"
  type        = string
}