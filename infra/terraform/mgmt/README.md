# Mgmt Account Setup

This terraform configuration sets up the AWS Organization mgmt account with:
- AWS Organizations with dev and prod member accounts
- Two domain registrations via Route 53:
  - Production domain (e.g., crowdunlocked.com)
  - Development domain (e.g., crowdunlockedbeta.com)
- Route 53 hosted zones for DNS management
- ACM certificates for TLS (one per domain)

## Prerequisites

1. AWS account (this will become your mgmt account)
2. AWS CLI configured with credentials
3. Terraform >= 1.6
4. S3 bucket for terraform state: `crowdunlocked-terraform-state`
5. DynamoDB table for state locking: `terraform-state-lock`

## Setup Steps

### 1. Create State Backend Resources

Before running terraform, create the S3 bucket and DynamoDB table:

```bash
# Create S3 bucket for terraform state
aws s3api create-bucket \
  --bucket crowdunlocked-terraform-state \
  --region us-east-1

aws s3api put-bucket-versioning \
  --bucket crowdunlocked-terraform-state \
  --versioning-configuration Status=Enabled

aws s3api put-bucket-encryption \
  --bucket crowdunlocked-terraform-state \
  --server-side-encryption-configuration '{
    "Rules": [{
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }]
  }'

# Create DynamoDB table for state locking
aws dynamodb create-table \
  --table-name terraform-state-lock \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1
```

### 2. Configure Variables

Copy the example file and fill in your values:

```bash
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your actual values
```

**Important Notes:**
- `dev_account_email` and `prod_account_email` must be unique email addresses
- Phone number format: `+[country code].[number]` (e.g., `+1.5551234567`)
- Domain registration costs vary by TLD (typically $12-50/year per domain)
- You'll be registering TWO domains (prod and dev), so double the cost
- Privacy protection is enabled by default to hide contact info from WHOIS

### 3. Initialize and Apply

```bash
terraform init
terraform plan
terraform apply
```

### 4. Domain Registration

Both domain registrations happen automatically when you apply. Note:
- Registration can take 5-15 minutes per domain
- You'll receive verification emails at the contact email for each domain
- Both domains will auto-renew annually
- Privacy protection hides your contact info from public WHOIS
- Total cost: ~$24-100/year for both domains (depends on TLDs)

### 5. Post-Setup

After successful apply:
1. Check your email for domain verification (if required)
2. Note the outputs (organization ID, account IDs, zone ID)
3. Configure AWS SSO in the AWS Console
4. Set up cross-account access for dev/prod accounts

## Outputs

- `organization_id`: AWS Organization ID
- `dev_account_id`: Development account ID
- `prod_account_id`: Production account ID
- `prod_route53_zone_id`: Production hosted zone ID
- `dev_route53_zone_id`: Development hosted zone ID
- `prod_acm_certificate_arn`: Production certificate ARN
- `dev_acm_certificate_arn`: Development certificate ARN
- `prod_domain_name`: Production domain name
- `dev_domain_name`: Development domain name

## Cost Estimate

- Domain registrations: $24-100/year (2 domains, depends on TLDs)
- Route 53 hosted zones: $1.00/month (2 zones @ $0.50 each)
- AWS Organizations: Free
- ACM certificates: Free

## Troubleshooting

### Domain Already Registered
If the domain is already registered elsewhere, you'll need to transfer it to Route 53 instead. Use `aws_route53domains_domain_transfer` resource.

### Email Verification
Some TLDs require email verification. Check your inbox and spam folder for verification emails from AWS.

### State Backend Issues
If you get state backend errors, ensure the S3 bucket and DynamoDB table exist in us-east-1.
