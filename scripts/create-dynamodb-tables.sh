#!/bin/bash
set -e

ENDPOINT=${1:-http://localhost:8000}
REGION=${2:-us-east-1}

echo "Creating DynamoDB tables at $ENDPOINT..."

TABLES=("bookings-local" "releases-local" "publicity-local" "social-local" "money-local")

for table in "${TABLES[@]}"; do
    echo "Creating table: $table"
    aws dynamodb create-table \
        --table-name "$table" \
        --attribute-definitions AttributeName=id,AttributeType=S \
        --key-schema AttributeName=id,KeyType=HASH \
        --billing-mode PAY_PER_REQUEST \
        --endpoint-url "$ENDPOINT" \
        --region "$REGION" \
        2>/dev/null || echo "  Table $table already exists"
done

echo "✅ All tables created"
