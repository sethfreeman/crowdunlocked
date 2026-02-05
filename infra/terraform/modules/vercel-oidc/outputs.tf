output "oidc_provider_arn" {
  description = "ARN of the Vercel OIDC provider"
  value       = aws_iam_openid_connect_provider.vercel.arn
}

output "vercel_role_arn" {
  description = "ARN of the IAM role for Vercel"
  value       = aws_iam_role.vercel.arn
}

output "vercel_role_name" {
  description = "Name of the IAM role for Vercel"
  value       = aws_iam_role.vercel.name
}

output "dynamodb_policy_arn" {
  description = "ARN of the DynamoDB policy for Vercel"
  value       = aws_iam_policy.vercel_dynamodb.arn
}

output "logs_policy_arn" {
  description = "ARN of the CloudWatch Logs policy for Vercel"
  value       = aws_iam_policy.vercel_logs.arn
}