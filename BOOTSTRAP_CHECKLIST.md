# Bootstrap Checklist

Use this checklist to set up Crowd Unlocked infrastructure from scratch.

## Prerequisites

- [ ] AWS account (will become mgmt account)
- [ ] GitHub repository created
- [ ] AWS CLI installed and configured
- [ ] OpenTofu >= 1.10 installed (or Terraform >= 1.6)
- [ ] Domain names decided:
  - Production: `crowdunlocked.com`
  - Development: `crowdunlockedbeta.com`

## Phase 1: AWS Backend Setup (One-Time)

### Create S3 Bucket for Terraform State

```bash
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
```

- [ ] S3 bucket created
- [ ] Versioning enabled
- [ ] Encryption enabled

### Create DynamoDB Table for State Locking

```bash
aws dynamodb create-table \
  --table-name terraform-state-lock \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1
```

- [ ] DynamoDB table created

## Phase 2: GitHub OIDC Setup

### Get AWS Account ID

```bash
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
echo "AWS Account ID: $AWS_ACCOUNT_ID"
```

- [ ] Account ID noted: _______________

### Create OIDC Provider

```bash
aws iam create-open-id-connect-provider \
  --url https://token.actions.githubusercontent.com \
  --client-id-list sts.amazonaws.com \
  --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1
```

- [ ] OIDC provider created

### Create IAM Role for GitHub Actions

```bash
# Set your GitHub org and repo
GITHUB_ORG="your-github-org"
GITHUB_REPO="crowdunlocked"

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

# Create role
aws iam create-role \
  --role-name GitHubActionsRole \
  --assume-role-policy-document file://trust-policy.json

# Attach admin policy (adjust as needed)
aws iam attach-role-policy \
  --role-name GitHubActionsRole \
  --policy-arn arn:aws:iam::aws:policy/AdministratorAccess
```

- [ ] Trust policy created
- [ ] IAM role created
- [ ] Admin policy attached
- [ ] Role ARN noted: _______________

## Phase 3: GitHub Secrets Configuration

Go to: `https://github.com/YOUR_ORG/crowdunlocked/settings/secrets/actions`

### AWS Authentication

- [ ] `AWS_ROLE_mgmt` = `arn:aws:iam::ACCOUNT_ID:role/GitHubActionsRole`
- [ ] `AWS_ROLE_dev` = (same as above for now)
- [ ] `AWS_ROLE_prod` = (same as above for now)

### Domain Configuration

- [ ] `PROD_DOMAIN_NAME` = `crowdunlocked.com`
- [ ] `DEV_DOMAIN_NAME` = `crowdunlockedbeta.com`

### AWS Account Emails (Must be unique!)

- [ ] `DEV_ACCOUNT_EMAIL` = `yourname+crowdunlocked-dev@gmail.com`
- [ ] `PROD_ACCOUNT_EMAIL` = `yourname+crowdunlocked-prod@gmail.com`

### Domain Registration Contact Info

- [ ] `DOMAIN_CONTACT_EMAIL` = `admin@crowdunlocked.com`
- [ ] `DOMAIN_CONTACT_FIRST_NAME` = `John`
- [ ] `DOMAIN_CONTACT_LAST_NAME` = `Doe`
- [ ] `DOMAIN_CONTACT_PHONE` = `+1.5551234567` (format: +[country].[number])
- [ ] `DOMAIN_CONTACT_ADDRESS` = `123 Main Street`
- [ ] `DOMAIN_CONTACT_CITY` = `San Francisco`
- [ ] `DOMAIN_CONTACT_STATE` = `CA`
- [ ] `DOMAIN_CONTACT_ZIP` = `94102`

## Phase 4: Update Configuration Files

### Update Domain Names

Edit `infra/terraform/mgmt/management.tfvars`:

```hcl
prod_domain_name = "crowdunlocked.com"      # Your actual domain
dev_domain_name  = "crowdunlockedbeta.com"  # Your actual domain
```

- [ ] Domain names updated in `management.tfvars`

### Update GitHub Repository URL

Edit `flux/clusters/dev/flux-system/gotk-sync.yaml`:

```yaml
url: ssh://git@github.com/YOUR_ORG/crowdunlocked
```

Edit `flux/clusters/prod/flux-system/gotk-sync.yaml`:

```yaml
url: ssh://git@github.com/YOUR_ORG/crowdunlocked
```

- [ ] GitHub URLs updated in Flux configs

## Phase 5: Initial Commit

### Create Feature Branch

```bash
git checkout -b feature/initial-infrastructure
git add .
git commit -m "feat(infra): initial terraform and CI/CD setup"
git push origin feature/initial-infrastructure
```

- [ ] Feature branch created
- [ ] Changes committed
- [ ] Pushed to GitHub

### Create Pull Request

1. Go to GitHub
2. Create PR from `feature/initial-infrastructure` to `develop`
3. Title: `feat(infra): initial terraform and CI/CD setup`
4. Review terraform plans in PR comments

- [ ] PR created
- [ ] Terraform plans reviewed
- [ ] All checks passing

### Merge to Develop

- [ ] PR merged to `develop`

## Phase 6: Monitor Initial Deployment

### Watch GitHub Actions

Go to: `https://github.com/YOUR_ORG/crowdunlocked/actions`

- [ ] Terraform workflow started
- [ ] Tests passing
- [ ] Terraform plan completed
- [ ] Terraform apply started

### Wait for Domain Registration

This takes 15-30 minutes. Check email for verification.

- [ ] Domain registration started
- [ ] Verification email received (check spam)
- [ ] Domains verified (if required)
- [ ] Terraform apply completed

### Verify AWS Resources

```bash
# Check organization
aws organizations describe-organization

# Check accounts
aws organizations list-accounts

# Check domains
aws route53domains list-domains

# Check hosted zones
aws route53 list-hosted-zones

# Check certificates
aws acm list-certificates --region us-east-1
```

- [ ] AWS Organization created
- [ ] Dev account created
- [ ] Prod account created
- [ ] Both domains registered
- [ ] Both hosted zones created
- [ ] Both certificates issued

## Phase 7: Update Cross-Account Roles

After dev and prod accounts are created, you need to create IAM roles in those accounts for GitHub Actions.

### Get Account IDs

```bash
cd infra/terraform/mgmt
terraform output dev_account_id
terraform output prod_account_id
```

- [ ] Dev account ID: _______________
- [ ] Prod account ID: _______________

### Create Roles in Dev and Prod Accounts

(This will be done via terraform in those accounts later)

For now, update GitHub Secrets:

- [ ] `AWS_ROLE_dev` = `arn:aws:iam::DEV_ACCOUNT_ID:role/GitHubActionsRole`
- [ ] `AWS_ROLE_prod` = `arn:aws:iam::PROD_ACCOUNT_ID:role/GitHubActionsRole`

## Phase 8: Test the Workflow

### Create Test Feature

```bash
git checkout develop
git pull origin develop
git checkout -b feature/test-workflow
echo "# Test" > TEST.md
git add TEST.md
git commit -m "test: verify CI/CD workflow"
git push origin feature/test-workflow
```

- [ ] Test branch created
- [ ] PR created
- [ ] CI checks passing
- [ ] PR merged

### Verify Auto-Deploy

- [ ] GitHub Actions deployed to dev
- [ ] No errors in workflow

## Phase 9: Promote to Production

When ready to set up production:

```bash
git checkout main
git pull origin main
git merge develop
git push origin main
```

- [ ] Merged to main
- [ ] Terraform applied to prod
- [ ] Production environment ready

## Phase 10: Set Up Flux

### Install Flux CLI

```bash
# macOS
brew install fluxcd/tap/flux

# Or download from https://fluxcd.io/flux/installation/
```

- [ ] Flux CLI installed

### Bootstrap Flux on Dev Cluster

```bash
# After EKS cluster is created
flux bootstrap github \
  --owner=YOUR_ORG \
  --repository=crowdunlocked \
  --branch=develop \
  --path=flux/clusters/dev \
  --personal
```

- [ ] Flux bootstrapped on dev cluster

### Bootstrap Flux on Prod Cluster

```bash
flux bootstrap github \
  --owner=YOUR_ORG \
  --repository=crowdunlocked \
  --branch=main \
  --path=flux/clusters/prod \
  --personal
```

- [ ] Flux bootstrapped on prod cluster

## Completion Checklist

- [ ] AWS Organization created
- [ ] Dev and prod accounts created
- [ ] Both domains registered and verified
- [ ] Terraform state backend working
- [ ] GitHub Actions workflows running
- [ ] Secrets configured correctly
- [ ] Flux deployed to both clusters
- [ ] Dev environment accessible
- [ ] Prod environment accessible
- [ ] Documentation reviewed

## Next Steps

Now you're ready for daily development! See:
- [DEVELOPMENT_WORKFLOW.md](DEVELOPMENT_WORKFLOW.md) for daily workflow
- [docs/TDD_GUIDE.md](docs/TDD_GUIDE.md) for testing practices

## Troubleshooting

### Terraform State Lock Issues

```bash
# If state is locked
aws dynamodb delete-item \
  --table-name terraform-state-lock \
  --key '{"LockID":{"S":"crowdunlocked-terraform-state/management/terraform.tfstate"}}'
```

### Domain Registration Stuck

- Check email for verification (including spam)
- Domain registration can take up to 30 minutes
- Check AWS Console → Route 53 → Registered domains

### GitHub Actions Failing

- Verify all secrets are set correctly
- Check IAM role trust policy
- Ensure OIDC provider is created
- Review workflow logs for specific errors

### Need Help?

Review the documentation:
- [INITIAL_SETUP_SUMMARY.md](INITIAL_SETUP_SUMMARY.md)
- [infra/terraform/CICD_SETUP.md](infra/terraform/CICD_SETUP.md)
