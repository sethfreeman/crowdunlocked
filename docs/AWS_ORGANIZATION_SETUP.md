# AWS Organization Setup Guide

This is a **one-time manual setup** for creating the AWS Organization and member accounts. After this, terraform will manage infrastructure within each account.

## Prerequisites
- AWS root account access
- Email addresses for dev and prod accounts (must be unique)

---

## Step 1: Create AWS Organization

1. Log into AWS Console with **root account**
2. Go to **AWS Organizations** service
3. Click **"Create organization"**
4. Choose **"Enable all features"**
5. Verify your email when prompted
6. Wait for confirmation email

**Result:** Your account is now the management account.

---

## Step 2: Create Dev Account

1. In AWS Organizations console
2. Click **"Add an AWS account"**
3. Choose **"Create an AWS account"**
4. Fill in:
   - **Account name**: `crowdunlocked-dev`
   - **Email**: Your dev email (e.g., `aws-dev@yourdomain.com`)
   - **IAM role name**: `OrganizationAccountAccessRole`
5. Click **"Create AWS account"**
6. Wait 5-10 minutes for account creation
7. Check email for verification

**Result:** Dev account created with ID like `123456789012`

---

## Step 3: Create Prod Account

1. Repeat Step 2 with:
   - **Account name**: `crowdunlocked-prod`
   - **Email**: Your prod email (e.g., `aws-prod@yourdomain.com`)
   - **IAM role name**: `OrganizationAccountAccessRole`

**Result:** Prod account created with ID like `987654321098`

---

## Step 4: Set Up Cross-Account Access

The `OrganizationAccountAccessRole` is automatically created. Test access:

### Access Dev Account
```bash
aws sts assume-role \
  --role-arn arn:aws:iam::DEV_ACCOUNT_ID:role/OrganizationAccountAccessRole \
  --role-session-name dev-session
```

### Access Prod Account
```bash
aws sts assume-role \
  --role-arn arn:aws:iam::PROD_ACCOUNT_ID:role/OrganizationAccountAccessRole \
  --role-session-name prod-session
```

---

## Step 5: Register Domains (Management Account)

### Option A: Via AWS Console (Easier)

1. Go to **Route 53** in management account
2. Click **"Registered domains"** → **"Register domain"**
3. Search for `crowdunlocked.com`
4. Add to cart, proceed to checkout
5. Fill in contact information
6. Enable privacy protection
7. Complete purchase (~$12/year)
8. Repeat for `crowdunlockedbeta.com`

### Option B: Via AWS CLI

```bash
# Register production domain
aws route53domains register-domain \
  --region us-east-1 \
  --domain-name crowdunlocked.com \
  --duration-in-years 1 \
  --auto-renew \
  --admin-contact file://contact.json \
  --registrant-contact file://contact.json \
  --tech-contact file://contact.json \
  --privacy-protect-admin-contact \
  --privacy-protect-registrant-contact \
  --privacy-protect-tech-contact

# Register dev domain
aws route53domains register-domain \
  --region us-east-1 \
  --domain-name crowdunlockedbeta.com \
  --duration-in-years 1 \
  --auto-renew \
  --admin-contact file://contact.json \
  --registrant-contact file://contact.json \
  --tech-contact file://contact.json \
  --privacy-protect-admin-contact \
  --privacy-protect-registrant-contact \
  --privacy-protect-tech-contact
```

**contact.json:**
```json
{
  "FirstName": "Your",
  "LastName": "Name",
  "ContactType": "PERSON",
  "AddressLine1": "123 Main St",
  "City": "San Francisco",
  "State": "CA",
  "CountryCode": "US",
  "ZipCode": "94102",
  "PhoneNumber": "+1.5551234567",
  "Email": "your-email@example.com"
}
```

**Note:** Domain registration can take 10-15 minutes. You'll receive verification emails.

---

## Step 6: Create Hosted Zones (Management Account)

After domains are registered:

```bash
# Create hosted zone for prod domain
aws route53 create-hosted-zone \
  --name crowdunlocked.com \
  --caller-reference $(date +%s) \
  --region us-east-1

# Create hosted zone for dev domain
aws route53 create-hosted-zone \
  --name crowdunlockedbeta.com \
  --caller-reference $(date +%s) \
  --region us-east-1
```

**Save the hosted zone IDs** - you'll need them for terraform.

---

## Step 7: Create ACM Certificates (Management Account)

```bash
# Request certificate for prod domain
aws acm request-certificate \
  --domain-name crowdunlocked.com \
  --subject-alternative-names "*.crowdunlocked.com" \
  --validation-method DNS \
  --region us-east-1

# Request certificate for dev domain
aws acm request-certificate \
  --domain-name crowdunlockedbeta.com \
  --subject-alternative-names "*.crowdunlockedbeta.com" \
  --validation-method DNS \
  --region us-east-1
```

**Save the certificate ARNs** - you'll need them for terraform.

---

## Step 8: Validate Certificates

1. Go to **ACM** in AWS Console (us-east-1)
2. Click on each certificate
3. Click **"Create records in Route 53"**
4. Click **"Create records"**
5. Wait 5-30 minutes for validation
6. Status will change from "Pending validation" to "Issued"

---

## Step 9: Update Terraform Variables

Create `infra/terraform/mgmt/terraform.tfvars`:

```hcl
# Account IDs (from Step 2 & 3)
dev_account_id  = "123456789012"
prod_account_id = "987654321098"

# Domain names
prod_domain_name = "crowdunlocked.com"
dev_domain_name  = "crowdunlockedbeta.com"

# Hosted Zone IDs (from Step 6)
prod_route53_zone_id = "Z1234567890ABC"
dev_route53_zone_id  = "Z0987654321XYZ"

# Certificate ARNs (from Step 7)
prod_acm_certificate_arn = "arn:aws:acm:us-east-1:ACCOUNT:certificate/abc123..."
dev_acm_certificate_arn  = "arn:aws:acm:us-east-1:ACCOUNT:certificate/xyz789..."
```

---

## Step 10: Configure AWS Profiles

Add to `~/.aws/config`:

```ini
[profile crowdunlocked-mgmt]
region = us-east-1

[profile crowdunlocked-dev]
region = us-west-2
role_arn = arn:aws:iam::DEV_ACCOUNT_ID:role/OrganizationAccountAccessRole
source_profile = crowdunlocked-mgmt

[profile crowdunlocked-prod]
region = us-west-2
role_arn = arn:aws:iam::PROD_ACCOUNT_ID:role/OrganizationAccountAccessRole
source_profile = crowdunlocked-mgmt
```

Test access:
```bash
aws sts get-caller-identity --profile crowdunlocked-mgmt
aws sts get-caller-identity --profile crowdunlocked-dev
aws sts get-caller-identity --profile crowdunlocked-prod
```

---

## Step 11: Update GitHub Secrets

Add these to your GitHub repository secrets:

```
AWS_ACCOUNT_ID_MGMT=111111111111
AWS_ACCOUNT_ID_DEV=222222222222
AWS_ACCOUNT_ID_PROD=333333333333

PROD_DOMAIN_NAME=crowdunlocked.com
DEV_DOMAIN_NAME=crowdunlockedbeta.com

PROD_ROUTE53_ZONE_ID=Z1234567890ABC
DEV_ROUTE53_ZONE_ID=Z0987654321XYZ

PROD_ACM_CERTIFICATE_ARN=arn:aws:acm:...
DEV_ACM_CERTIFICATE_ARN=arn:aws:acm:...
```

---

## What's Next?

Now that the organization and accounts are set up, terraform will manage:
- ✅ EKS clusters (dev and prod)
- ✅ DynamoDB tables
- ✅ VPCs and networking
- ✅ IAM roles and policies
- ✅ CloudFront distributions (using the certificates you created)

Run terraform:
```bash
cd infra/terraform/dev
terraform init
terraform plan
terraform apply

cd ../prod
terraform init
terraform plan
terraform apply
```

---

## Troubleshooting

### Can't access member accounts
- Verify `OrganizationAccountAccessRole` exists in member accounts
- Check you're using management account credentials as source

### Domain registration pending
- Can take 10-15 minutes
- Check email for verification
- Some TLDs require additional verification

### Certificate validation stuck
- Verify DNS records were created in Route 53
- Check nameservers are correct
- Can take up to 30 minutes

### Organization creation failed
- Only one organization per account
- Must use root account credentials
- Check email verification

---

## Cost Estimate

**One-time:**
- Domain registration: ~$12/year each = $24/year

**Monthly:**
- AWS Organizations: Free
- Route 53 hosted zones: $0.50/zone/month = $1/month
- ACM certificates: Free

**Total:** ~$24/year + $12/month for infrastructure
