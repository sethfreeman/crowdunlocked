#!/bin/bash
set -e

# Script to manually create GitHub OIDC provider and role
# Usage: ./scripts/create-oidc-role.sh <environment> <github-org> <github-repo>

ENVIRONMENT=$1
GITHUB_ORG=$2
GITHUB_REPO=$3

if [ -z "$ENVIRONMENT" ] || [ -z "$GITHUB_ORG" ] || [ -z "$GITHUB_REPO" ]; then
  echo "Usage: $0 <environment> <github-org> <github-repo>"
  echo "Example: $0 mgmt sethfreeman crowdunlocked"
  exit 1
fi

echo "Creating OIDC provider and role for $ENVIRONMENT environment..."

# Create OIDC provider (if it doesn't exist)
OIDC_PROVIDER_ARN=$(aws iam list-open-id-connect-providers --query "OpenIDConnectProviderList[?contains(Arn, 'token.actions.githubusercontent.com')].Arn" --output text)

if [ -z "$OIDC_PROVIDER_ARN" ]; then
  echo "Creating OIDC provider..."
  OIDC_PROVIDER_ARN=$(aws iam create-open-id-connect-provider \
    --url https://token.actions.githubusercontent.com \
    --client-id-list sts.amazonaws.com \
    --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1 1c58a3a8518e8759bf075b76b750d4f2df264fcd \
    --tags Key=Name,Value=github-actions-oidc Key=Environment,Value=$ENVIRONMENT Key=ManagedBy,Value=manual \
    --query 'OpenIDConnectProviderArn' --output text)
  echo "Created OIDC provider: $OIDC_PROVIDER_ARN"
else
  echo "OIDC provider already exists: $OIDC_PROVIDER_ARN"
fi

# Create trust policy
cat > /tmp/trust-policy-$ENVIRONMENT.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "$OIDC_PROVIDER_ARN"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
        },
        "StringLike": {
          "token.actions.githubusercontent.com:sub": "repo:$GITHUB_ORG/$GITHUB_REPO:*"
        }
      }
    }
  ]
}
EOF

# Create IAM role
ROLE_NAME="github-actions-$ENVIRONMENT"
echo "Creating IAM role: $ROLE_NAME..."

ROLE_ARN=$(aws iam create-role \
  --role-name $ROLE_NAME \
  --assume-role-policy-document file:///tmp/trust-policy-$ENVIRONMENT.json \
  --tags Key=Name,Value=$ROLE_NAME Key=Environment,Value=$ENVIRONMENT Key=ManagedBy,Value=manual \
  --query 'Role.Arn' --output text 2>/dev/null || \
  aws iam get-role --role-name $ROLE_NAME --query 'Role.Arn' --output text)

echo "Role ARN: $ROLE_ARN"

# Create inline policy for terraform permissions
cat > /tmp/role-policy-$ENVIRONMENT.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject",
        "s3:ListBucket",
        "dynamodb:GetItem",
        "dynamodb:PutItem",
        "dynamodb:DeleteItem",
        "eks:*",
        "ec2:*",
        "iam:*",
        "dynamodb:*",
        "logs:*",
        "cloudfront:*",
        "acm:*",
        "route53:*",
        "secretsmanager:*",
        "kms:*",
        "autoscaling:*",
        "elasticloadbalancing:*"
      ],
      "Resource": "*"
    }
  ]
}
EOF

echo "Attaching policy to role..."
aws iam put-role-policy \
  --role-name $ROLE_NAME \
  --policy-name terraform-permissions \
  --policy-document file:///tmp/role-policy-$ENVIRONMENT.json

echo ""
echo "✅ OIDC setup complete for $ENVIRONMENT!"
echo ""
echo "Add this to GitHub Secrets:"
echo "AWS_ROLE_$ENVIRONMENT=$ROLE_ARN"
echo ""

# Cleanup
rm /tmp/trust-policy-$ENVIRONMENT.json /tmp/role-policy-$ENVIRONMENT.json
