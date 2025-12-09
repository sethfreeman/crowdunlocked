# =============================================================================
# VPC Outputs
# =============================================================================

output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "private_subnet_ids" {
  description = "IDs of private subnets"
  value       = aws_subnet.private[*].id
}

output "public_subnet_ids" {
  description = "IDs of public subnets"
  value       = aws_subnet.public[*].id
}

# =============================================================================
# EKS Outputs
# =============================================================================

output "eks_cluster_name" {
  description = "Name of the EKS cluster"
  value       = aws_eks_cluster.main.name
}

output "eks_cluster_endpoint" {
  description = "Endpoint for EKS cluster API server"
  value       = aws_eks_cluster.main.endpoint
}

output "eks_cluster_arn" {
  description = "ARN of the EKS cluster"
  value       = aws_eks_cluster.main.arn
}

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
# ECR Outputs
# =============================================================================

output "ecr_repository_urls" {
  description = "Map of ECR repository URLs"
  value = merge(
    { for k, v in aws_ecr_repository.services : k => v.repository_url },
    { web = aws_ecr_repository.web.repository_url }
  )
}
