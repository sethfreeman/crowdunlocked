output "eks_cluster_name" {
  value = aws_eks_cluster.main.name
}

output "eks_cluster_endpoint" {
  value = aws_eks_cluster.main.endpoint
}

output "cloudfront_distribution_id" {
  value = length(aws_cloudfront_distribution.main) > 0 ? aws_cloudfront_distribution.main[0].id : null
}

output "cloudfront_domain_name" {
  value = length(aws_cloudfront_distribution.main) > 0 ? aws_cloudfront_distribution.main[0].domain_name : null
}

output "alb_dns_name" {
  value = aws_lb.main.dns_name
}

output "github_actions_role_arn" {
  description = "ARN of the GitHub Actions IAM role for prod account"
  value       = module.github_oidc.role_arn
}
