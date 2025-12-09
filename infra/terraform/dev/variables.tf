# =============================================================================
# Core Variables
# =============================================================================


variable "aws_region" {
  description = "AWS region for infrastructure deployment"
  type        = string
  default     = "us-west-2"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

# =============================================================================
# GitHub Integration
# =============================================================================

variable "github_org" {
  description = "GitHub organization name for CI/CD integration"
  type        = string
}

variable "github_repo" {
  description = "GitHub repository name for CI/CD integration"
  type        = string
}
