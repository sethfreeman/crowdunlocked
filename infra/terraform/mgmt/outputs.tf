output "organization_id" {
  value = aws_organizations_organization.main.id
}

output "dev_account_id" {
  value = aws_organizations_account.dev.id
}

output "prod_account_id" {
  value = aws_organizations_account.prod.id
}

output "prod_route53_zone_id" {
  value = aws_route53_zone.prod.zone_id
}

output "dev_route53_zone_id" {
  value = aws_route53_zone.dev.zone_id
}

output "prod_acm_certificate_arn" {
  value = aws_acm_certificate.prod.arn
}

output "dev_acm_certificate_arn" {
  value = aws_acm_certificate.dev.arn
}

output "prod_domain_name" {
  value = var.prod_domain_name
}

output "dev_domain_name" {
  value = var.dev_domain_name
}
