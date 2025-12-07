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

# AWS Organizations
resource "aws_organizations_organization" "main" {
  aws_service_access_principals = [
    "sso.amazonaws.com",
    "cloudtrail.amazonaws.com",
    "config.amazonaws.com"
  ]
  feature_set = "ALL"
}

resource "aws_organizations_account" "dev" {
  name      = "crowdunlocked-dev"
  email     = var.dev_account_email
  role_name = "OrganizationAccountAccessRole"
}

resource "aws_organizations_account" "prod" {
  name      = "crowdunlocked-prod"
  email     = var.prod_account_email
  role_name = "OrganizationAccountAccessRole"
}

# Identity Center (SSO)
resource "aws_ssoadmin_permission_set" "admin" {
  name             = "AdministratorAccess"
  instance_arn     = tolist(data.aws_ssoadmin_instances.main.arns)[0]
  session_duration = "PT8H"
}

data "aws_ssoadmin_instances" "main" {}

resource "aws_ssoadmin_managed_policy_attachment" "admin" {
  instance_arn       = tolist(data.aws_ssoadmin_instances.main.arns)[0]
  permission_set_arn = aws_ssoadmin_permission_set.admin.arn
  managed_policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}

# Production Domain Registration
resource "aws_route53domains_registered_domain" "prod" {
  domain_name = var.prod_domain_name

  dynamic "name_server" {
    for_each = aws_route53_zone.prod.name_servers
    content {
      name = name_server.value
    }
  }

  auto_renew = true

  admin_contact {
    contact_type   = var.domain_contact_type
    email          = var.domain_contact_email
    first_name     = var.domain_contact_first_name
    last_name      = var.domain_contact_last_name
    phone_number   = var.domain_contact_phone
    address_line_1 = var.domain_contact_address
    city           = var.domain_contact_city
    state          = var.domain_contact_state
    zip_code       = var.domain_contact_zip
    country_code   = var.domain_contact_country
  }

  registrant_contact {
    contact_type   = var.domain_contact_type
    email          = var.domain_contact_email
    first_name     = var.domain_contact_first_name
    last_name      = var.domain_contact_last_name
    phone_number   = var.domain_contact_phone
    address_line_1 = var.domain_contact_address
    city           = var.domain_contact_city
    state          = var.domain_contact_state
    zip_code       = var.domain_contact_zip
    country_code   = var.domain_contact_country
  }

  tech_contact {
    contact_type   = var.domain_contact_type
    email          = var.domain_contact_email
    first_name     = var.domain_contact_first_name
    last_name      = var.domain_contact_last_name
    phone_number   = var.domain_contact_phone
    address_line_1 = var.domain_contact_address
    city           = var.domain_contact_city
    state          = var.domain_contact_state
    zip_code       = var.domain_contact_zip
    country_code   = var.domain_contact_country
  }

  admin_privacy      = true
  registrant_privacy = true
  tech_privacy       = true

  lifecycle {
    prevent_destroy = true
  }
}

# Development Domain Registration
resource "aws_route53domains_registered_domain" "dev" {
  domain_name = var.dev_domain_name

  dynamic "name_server" {
    for_each = aws_route53_zone.dev.name_servers
    content {
      name = name_server.value
    }
  }

  auto_renew = true

  admin_contact {
    contact_type   = var.domain_contact_type
    email          = var.domain_contact_email
    first_name     = var.domain_contact_first_name
    last_name      = var.domain_contact_last_name
    phone_number   = var.domain_contact_phone
    address_line_1 = var.domain_contact_address
    city           = var.domain_contact_city
    state          = var.domain_contact_state
    zip_code       = var.domain_contact_zip
    country_code   = var.domain_contact_country
  }

  registrant_contact {
    contact_type   = var.domain_contact_type
    email          = var.domain_contact_email
    first_name     = var.domain_contact_first_name
    last_name      = var.domain_contact_last_name
    phone_number   = var.domain_contact_phone
    address_line_1 = var.domain_contact_address
    city           = var.domain_contact_city
    state          = var.domain_contact_state
    zip_code       = var.domain_contact_zip
    country_code   = var.domain_contact_country
  }

  tech_contact {
    contact_type   = var.domain_contact_type
    email          = var.domain_contact_email
    first_name     = var.domain_contact_first_name
    last_name      = var.domain_contact_last_name
    phone_number   = var.domain_contact_phone
    address_line_1 = var.domain_contact_address
    city           = var.domain_contact_city
    state          = var.domain_contact_state
    zip_code       = var.domain_contact_zip
    country_code   = var.domain_contact_country
  }

  admin_privacy      = true
  registrant_privacy = true
  tech_privacy       = true

  lifecycle {
    prevent_destroy = true
  }
}

# Route 53 Hosted Zones
resource "aws_route53_zone" "prod" {
  name = var.prod_domain_name
}

resource "aws_route53_zone" "dev" {
  name = var.dev_domain_name
}

# ACM Certificates
resource "aws_acm_certificate" "prod" {
  domain_name               = var.prod_domain_name
  subject_alternative_names = ["*.${var.prod_domain_name}"]
  validation_method         = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate" "dev" {
  domain_name               = var.dev_domain_name
  subject_alternative_names = ["*.${var.dev_domain_name}"]
  validation_method         = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

# Certificate Validation Records
resource "aws_route53_record" "prod_cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.prod.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = aws_route53_zone.prod.zone_id
}

resource "aws_route53_record" "dev_cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.dev.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = aws_route53_zone.dev.zone_id
}

resource "aws_acm_certificate_validation" "prod" {
  certificate_arn         = aws_acm_certificate.prod.arn
  validation_record_fqdns = [for record in aws_route53_record.prod_cert_validation : record.fqdn]
}

resource "aws_acm_certificate_validation" "dev" {
  certificate_arn         = aws_acm_certificate.dev.arn
  validation_record_fqdns = [for record in aws_route53_record.dev_cert_validation : record.fqdn]
}
