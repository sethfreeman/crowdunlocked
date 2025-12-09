#!/bin/bash
set -e

# Script to generate AWS Console login URL for member accounts
# Usage: ./scripts/get-console-url.sh <environment>

ENVIRONMENT=$1

if [ -z "$ENVIRONMENT" ]; then
  echo "Usage: $0 <environment>"
  echo "Example: $0 dev"
  exit 1
fi

echo "Generating console URL for $ENVIRONMENT environment..."

# Get temporary credentials by assuming the role from mgmt account
CREDS=$(aws sts assume-role \
  --role-arn arn:aws:iam::$(AWS_PROFILE=crowdunlocked-$ENVIRONMENT aws sts get-caller-identity --query Account --output text):role/OrganizationAccountAccessRole \
  --role-session-name console-session \
  --duration-seconds 3600 \
  --output json)

# Extract credentials
ACCESS_KEY=$(echo $CREDS | jq -r '.Credentials.AccessKeyId')
SECRET_KEY=$(echo $CREDS | jq -r '.Credentials.SecretAccessKey')
SESSION_TOKEN=$(echo $CREDS | jq -r '.Credentials.SessionToken')

# Create session JSON
SESSION_JSON=$(cat <<EOF
{
  "sessionId": "$ACCESS_KEY",
  "sessionKey": "$SECRET_KEY",
  "sessionToken": "$SESSION_TOKEN"
}
EOF
)

# URL encode the session
SESSION_ENCODED=$(echo -n "$SESSION_JSON" | jq -sRr @uri)

# Generate signin token
SIGNIN_TOKEN=$(curl -s "https://signin.aws.amazon.com/federation?Action=getSigninToken&SessionDuration=3600&Session=$SESSION_ENCODED" | jq -r '.SigninToken')

# Generate console URL
CONSOLE_URL="https://signin.aws.amazon.com/federation?Action=login&Issuer=crowdunlocked&Destination=https%3A%2F%2Fconsole.aws.amazon.com%2F&SigninToken=$SIGNIN_TOKEN"

echo ""
echo "✅ Console URL generated!"
echo ""
echo "Click this URL to access the $ENVIRONMENT account console:"
echo ""
echo "$CONSOLE_URL"
echo ""
echo "Note: This URL is valid for 1 hour."
echo ""
