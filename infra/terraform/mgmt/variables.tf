variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "prod_domain_name" {
  description = "Production domain name (must be manually registered)"
  type        = string
}

variable "dev_domain_name" {
  description = "Development domain name (must be manually registered)"
  type        = string
}

variable "github_org" {
  description = "GitHub organization name"
  type        = string
}

variable "github_repo" {
  description = "GitHub repository name"
  type        = string
}

variable "dev_account_id" {
  description = "AWS account ID for dev environment"
  type        = string
}

variable "prod_account_id" {
  description = "AWS account ID for prod environment"
  type        = string
}
