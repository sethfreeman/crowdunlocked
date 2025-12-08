output "eks_cluster_name" {
  value = aws_eks_cluster.main.name
}

output "eks_cluster_endpoint" {
  value = aws_eks_cluster.main.endpoint
}

output "dynamodb_tables" {
  value = {
    bookings      = aws_dynamodb_table.bookings.name
    releases      = aws_dynamodb_table.releases.name
    publicity     = aws_dynamodb_table.publicity.name
    social        = aws_dynamodb_table.social.name
    money         = aws_dynamodb_table.money.name
  }
}
