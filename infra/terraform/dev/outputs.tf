# =============================================================================
# DynamoDB Outputs
# =============================================================================

output "dynamodb_tables" {
  description = "Map of DynamoDB table names"
  value = {
    bookings  = aws_dynamodb_table.bookings.name
    venues    = aws_dynamodb_table.venues.name
    releases  = aws_dynamodb_table.releases.name
    publicity = aws_dynamodb_table.publicity.name
    social    = aws_dynamodb_table.social.name
    money     = aws_dynamodb_table.money.name
  }
}

# =============================================================================
# Vercel OIDC Outputs
# =============================================================================

output "vercel_role_arn" {
  description = "ARN of the IAM role for Vercel OIDC"
  value       = module.vercel_oidc.vercel_role_arn
}

output "vercel_role_name" {
  description = "Name of the IAM role for Vercel OIDC"
  value       = module.vercel_oidc.vercel_role_name
}

output "vercel_oidc_provider_arn" {
  description = "ARN of the Vercel OIDC provider"
  value       = module.vercel_oidc.oidc_provider_arn
}
