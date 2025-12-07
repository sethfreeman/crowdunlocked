# Terraform CI/CD Setup Guide

This guide walks you through setting up automated terraform deployments via GitHub Actions.

## Architecture

- **Committed tfvars**: Non-sensitive config (region, environment) committed to git
- **GitHub Secrets**: Sensitive data (emails, personal info) stored as secrets
- **Environment Variables**: Terraform reads `TF_VAR_*` from GitHub Actions
- **OIDC Authentication**: GitHub Actions assumes AWS IAM roles (no long-lived credentials)

## Prerequisites

1. AWS account (will become mgmt account)
2. GitHub repository with Actions enabled
3. Domain name decided (doesn't need to be purchased yet)

## Setup Steps

### 1. Create S3 Backend (One-Time Manual Setup)

Before running terraform via CI/CD, create the state backend:

```bash
# Set your AWS credentials locally
export AWS_PROFILE=your-profile

# Create S3 bucket
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

# Create DynamoDB table
aws dynamodb create-table \
  --table-name terraform-state-lock \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1
```

### 2. Configure GitHub OIDC with AWS

Create an IAM OIDC provider and roles for GitHub Actions:

```bash
# Get your AWS account ID
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
GITHUB_ORG="your-github-org"
GITHUB_REPO="your-repo-name"

# Create OIDC provider
aws iam create-open-id-connect-provider \
  --url https://token.actions.githubusercontent.com \
  --client-id-list sts.amazonaws.com \
  --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1

# Create trust policy
cat > trust-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
        },
        "StringLike": {
          "token.actions.githubusercontent.com:sub": "repo:${GITHUB_ORG}/${GITHUB_REPO}:*"
        }
      }
    }
  ]
}
EOF

# Create IAM role for GitHub Actions
aws iam create-role \
  --role-name GitHubActionsRole \
  --assume-role-policy-document file://trust-policy.json

# Attach admin policy (adjust permissions as needed)
aws iam attach-role-policy \
  --role-name GitHubActionsRole \
  --policy-arn arn:aws:iam::aws:policy/AdministratorAccess
```

### 3. Configure GitHub Secrets

Go to your GitHub repository → Settings → Secrets and variables → Actions

Add these secrets:

**AWS Authentication:**
- `AWS_ROLE_mgmt`: `arn:aws:iam::YOUR_ACCOUNT_ID:role/GitHubActionsRole`
- `AWS_ROLE_dev`: (same for now, will update after org is created)
- `AWS_ROLE_prod`: (same for now, will update after org is created)

**Domain Configuration:**
- `PROD_DOMAIN_NAME`: `crowdunlocked.com`
- `DEV_DOMAIN_NAME`: `crowdunlockedbeta.com`

**AWS Account Emails:**
- `DEV_ACCOUNT_EMAIL`: `yourname+crowdunlocked-dev@gmail.com`
- `PROD_ACCOUNT_EMAIL`: `yourname+crowdunlocked-prod@gmail.com`

**Domain Registration Contact Info:**
- `DOMAIN_CONTACT_EMAIL`: `admin@crowdunlocked.com`
- `DOMAIN_CONTACT_FIRST_NAME`: `John`
- `DOMAIN_CONTACT_LAST_NAME`: `Doe`
- `DOMAIN_CONTACT_PHONE`: `+1.5551234567`
- `DOMAIN_CONTACT_ADDRESS`: `123 Main Street`
- `DOMAIN_CONTACT_CITY`: `San Francisco`
- `DOMAIN_CONTACT_STATE`: `CA`
- `DOMAIN_CONTACT_ZIP`: `94102`

### 4. Update Domain Name in Committed Config

Edit `infra/terraform/mgmt/management.tfvars`:
```hcl
prod_domain_name = "crowdunlocked.com"
dev_domain_name  = "crowdunlockedbeta.com"
```

Commit and push this change.

### 5. Configure GitHub Environments (Optional but Recommended)

Go to Settings → Environments and create:
- `management` - Require manual approval for applies
- `dev` - Auto-deploy on main branch
- `prod` - Require manual approval for applies

### 6. First Deployment

1. Create a PR with your terraform changes
2. GitHub Actions will run `terraform plan` for all environments
3. Review the plans in the PR
4. Merge to `main` branch
5. GitHub Actions will apply changes sequentially (management → dev → prod)
6. If you configured environments, approve the deployments

### 7. Post-Deployment

After mgmt account is created:
1. Note the dev and prod account IDs from terraform outputs
2. Create IAM roles in those accounts for GitHub Actions
3. Update `AWS_ROLE_dev` and `AWS_ROLE_prod` secrets with the new role ARNs

## Workflow Behavior

**On Pull Request (any branch):**
- Runs `terraform plan` for all environments
- Posts plan output as PR comment (if configured)
- No changes applied
- Must pass before merge

**On Push to develop:**
- Runs `terraform plan` for all environments
- Applies changes to management + dev environments
- Auto-approves (no manual step)
- Flux deploys to dev EKS cluster automatically

**On Push to main:**
- Runs `terraform plan` for all environments
- Applies changes to management + prod environments
- Auto-approves (no manual step)
- Flux deploys to prod EKS cluster automatically

**Manual Trigger:**
- Can deploy specific environment via workflow_dispatch
- Useful for testing or emergency deployments

## Branching Strategy

```
feature/new-feature
        ↓ (PR + tests pass)
    develop ──────────────→ Dev Environment (crowdunlockedbeta.com)
        ↓ (PR + tests pass)
     main ────────────────→ Prod Environment (crowdunlocked.com)
```

**Deployment Flow:**
1. Create feature branch from `develop`
2. Write tests first (TDD)
3. Implement feature
4. Create PR to `develop`
5. Tests run automatically
6. Self-review and merge (solo dev)
7. Auto-deploys to dev environment
8. Test on dev
9. Create PR from `develop` to `main`
10. Auto-deploys to prod environment

**No manual approvals** - Just PR review required

## Local Development

You can still run terraform locally for testing:

```bash
cd infra/terraform/mgmt

# Option 1: Use environment variables
export TF_VAR_prod_domain_name="crowdunlocked.com"
export TF_VAR_dev_domain_name="crowdunlockedbeta.com"
export TF_VAR_dev_account_email="yourname+dev@gmail.com"
# ... export other vars ...
terraform plan -var-file=management.tfvars

# Option 2: Create local terraform.tfvars (gitignored)
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
terraform plan -var-file=management.tfvars
```

## Domain Strategy

This setup registers **two separate domains**:
- **Production**: `crowdunlocked.com` - Your main public domain
- **Development**: `crowdunlockedbeta.com` - Separate domain for testing

**Why two domains instead of subdomains?**
- Cleaner separation between environments
- No risk of dev affecting prod DNS
- Can use different registrars or configurations
- Better for testing DNS changes

**Cost**: ~$24-100/year total (depends on TLDs chosen)

## Security Notes

- No sensitive data in git repository
- All secrets stored in GitHub Secrets (encrypted at rest)
- OIDC authentication (no long-lived AWS credentials)
- Environment protection rules prevent accidental deployments
- Terraform state stored in encrypted S3 bucket

## Troubleshooting

**"Error: No valid credential sources found"**
- Check that OIDC provider is created in AWS
- Verify GitHub role ARN in secrets
- Ensure trust policy allows your repo

**"Error: Backend initialization required"**
- Ensure S3 bucket and DynamoDB table exist
- Check bucket name matches in backend config

**"Error: Required variable not set"**
- Verify all required secrets are set in GitHub
- Check secret names match `TF_VAR_*` format in workflow
