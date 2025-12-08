output "prod_route53_zone_id" {
  value = data.aws_route53_zone.prod.zone_id
}

output "dev_route53_zone_id" {
  value = data.aws_route53_zone.dev.zone_id
}

output "prod_acm_certificate_arn" {
  value = data.aws_acm_certificate.prod.arn
}

output "dev_acm_certificate_arn" {
  value = data.aws_acm_certificate.dev.arn
}

output "prod_domain_name" {
  value = var.prod_domain_name
}

output "dev_domain_name" {
  value = var.dev_domain_name
}

output "github_actions_role_arn" {
  description = "ARN of the GitHub Actions IAM role for mgmt account"
  value       = module.github_oidc.role_arn
}
