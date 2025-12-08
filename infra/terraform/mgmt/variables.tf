variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "prod_domain_name" {
  description = "Production domain name (e.g., crowdunlocked.com)"
  type        = string
}

variable "dev_domain_name" {
  description = "Development domain name (e.g., crowdunlockedbeta.com)"
  type        = string
}

variable "dev_account_email" {
  description = "Email for dev AWS account"
  type        = string
}

variable "prod_account_email" {
  description = "Email for prod AWS account"
  type        = string
}

# Domain registration contact information
variable "domain_contact_type" {
  description = "Contact type (PERSON, COMPANY, ASSOCIATION, PUBLIC_BODY, RESELLER)"
  type        = string
  default     = "COMPANY"
}

variable "domain_contact_email" {
  description = "Contact email for domain registration"
  type        = string
}

variable "domain_contact_first_name" {
  description = "Contact first name"
  type        = string
}

variable "domain_contact_last_name" {
  description = "Contact last name"
  type        = string
}

variable "domain_contact_phone" {
  description = "Contact phone number (format: +1.1234567890)"
  type        = string
}

variable "domain_contact_address" {
  description = "Contact address line 1"
  type        = string
}

variable "domain_contact_city" {
  description = "Contact city"
  type        = string
}

variable "domain_contact_state" {
  description = "Contact state/province"
  type        = string
}

variable "domain_contact_zip" {
  description = "Contact ZIP/postal code"
  type        = string
}

variable "domain_contact_country" {
  description = "Contact country code (e.g., US, CA, GB)"
  type        = string
  default     = "US"
}

variable "github_org" {
  description = "GitHub organization name"
  type        = string
}

variable "github_repo" {
  description = "GitHub repository name"
  type        = string
}
