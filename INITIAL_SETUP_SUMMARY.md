# Initial Setup Summary

This document summarizes the infrastructure and workflow setup for Crowd Unlocked.

## Infrastructure Tool

This project uses **OpenTofu** (open-source Terraform fork) instead of HashiCorp Terraform. OpenTofu is:
- Fully open source (MPL 2.0 license)
- Drop-in compatible with Terraform
- Community-driven and part of the Linux Foundation
- No licensing concerns for commercial use

All commands use `tofu` instead of `terraform`, but the syntax is identical.

## What Was Configured

### 1. Domain Strategy
- **Production**: `crowdunlocked.com`
- **Development**: `crowdunlockedbeta.com`
- Both domains registered via Route 53 in mgmt account
- Separate hosted zones and ACM certificates for each

### 2. OpenTofu/Terraform Structure

```
infra/terraform/
├── management/
│   ├── main.tf                    # AWS Org, domains, certificates
│   ├── variables.tf               # Variable definitions
│   ├── outputs.tf                 # Outputs for other environments
│   ├── management.tfvars          # Non-sensitive config (committed)
│   ├── terraform.tfvars.example   # Template for sensitive data
│   └── README.md                  # Setup instructions
├── dev/
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   └── dev.tfvars                 # Non-sensitive config (committed)
└── prod/
    ├── main.tf
    ├── variables.tf
    ├── outputs.tf
    └── prod.tfvars                # Non-sensitive config (committed)
```

**Sensitive data** (emails, personal info) → GitHub Secrets as `TF_VAR_*`

### 3. Branching Strategy

```
feature/new-feature
        ↓ (PR + tests)
    develop ──────────→ Dev Environment (crowdunlockedbeta.com)
        ↓ (PR + tests)
     main ────────────→ Prod Environment (crowdunlocked.com)
```

### 4. CI/CD Pipelines

#### OpenTofu Workflow (`.github/workflows/terraform.yaml`)
- **PR**: Plan for all environments
- **Push to develop**: Apply to management + dev
- **Push to main**: Apply to management + prod
- No manual approvals

#### Services Workflow (`.github/workflows/services.yaml`)
- **PR**: Test all services, lint, format check
- **Push to develop**: Build images, tag as `dev-{sha}`, update manifests
- **Push to main**: Build images, tag as `prod-{sha}`, update manifests
- Flux auto-deploys from git

### 5. GitOps with Flux

- **Dev cluster**: Watches `develop` branch → `k8s/overlays/dev`
- **Prod cluster**: Watches `main` branch → `k8s/overlays/prod`
- Auto-reconciles every 10 minutes
- No manual approval needed

### 6. Test-Driven Development

- All new code requires tests first
- Local testing with `make test` or `make test-watch`
- CI runs tests on every PR
- Coverage tracking enabled

### 7. Documentation Created

- `DEVELOPMENT_WORKFLOW.md` - Daily development guide
- `docs/TDD_GUIDE.md` - Test-driven development practices
- `infra/terraform/CICD_SETUP.md` - CI/CD setup instructions
- `infra/terraform/mgmt/README.md` - Mgmt account setup
- Updated `Makefile` with test and quality targets
- Updated `README.md` with workflow overview

## Next Steps

### 1. Initial Bootstrap (One-Time)

```bash
# Create S3 backend for OpenTofu state
aws s3api create-bucket --bucket crowdunlocked-terraform-state --region us-east-1
aws dynamodb create-table --table-name terraform-state-lock ...

# Set up GitHub OIDC provider in AWS
# See infra/terraform/CICD_SETUP.md for commands
```

### 2. Configure GitHub Secrets

Add these secrets to your GitHub repository:

**AWS Authentication:**
- `AWS_ROLE_mgmt`
- `AWS_ROLE_dev`
- `AWS_ROLE_prod`

**Domain Configuration:**
- `PROD_DOMAIN_NAME`: `crowdunlocked.com`
- `DEV_DOMAIN_NAME`: `crowdunlockedbeta.com`

**AWS Account Emails:**
- `DEV_ACCOUNT_EMAIL`
- `PROD_ACCOUNT_EMAIL`

**Domain Contact Info:**
- `DOMAIN_CONTACT_EMAIL`
- `DOMAIN_CONTACT_FIRST_NAME`
- `DOMAIN_CONTACT_LAST_NAME`
- `DOMAIN_CONTACT_PHONE`
- `DOMAIN_CONTACT_ADDRESS`
- `DOMAIN_CONTACT_CITY`
- `DOMAIN_CONTACT_STATE`
- `DOMAIN_CONTACT_ZIP`

### 3. Create Feature Branch

```bash
git checkout develop
git checkout -b feature/initial-infrastructure
git add .
git commit -m "feat(infra): initial terraform and CI/CD setup"
git push origin feature/initial-infrastructure
```

### 4. Create Pull Request

- Base: `develop`
- Title: `feat(infra): initial terraform and CI/CD setup`
- Review the terraform plans in PR
- Merge when ready

### 5. Deploy to Dev

After merging to `develop`:
- OpenTofu creates AWS Organization, accounts, and registers domains
- This is a one-time bootstrap
- Takes ~15-30 minutes (domain registration is slow)

### 6. Verify Dev Deployment

```bash
# Check OpenTofu outputs
cd infra/terraform/mgmt
tofu output

# Verify domains are registered
aws route53domains list-domains

# Check hosted zones
aws route53 list-hosted-zones
```

### 7. Promote to Production

When ready:
```bash
git checkout main
git merge develop
git push origin main
```

## Cost Estimate

### Monthly Costs
- Route 53 hosted zones: $1.00/month (2 zones)
- EKS Auto Mode: Pay per pod (starts at ~$0)
- DynamoDB: Pay per request (starts at ~$0)
- CloudFront: Pay per request (starts at ~$0)

### Annual Costs
- Domain registrations: $24-100/year (2 domains)

### Total First Year
- ~$36-112 (mostly domain registration)
- Scales with usage

## Workflow Summary

### Daily Development

1. Create feature branch from `develop`
2. Write tests first (TDD)
3. Implement feature
4. Run `make test` locally
5. Push and create PR
6. CI runs tests automatically
7. Merge to `develop`
8. Auto-deploys to dev
9. Test on `crowdunlockedbeta.com`
10. Merge `develop` to `main`
11. Auto-deploys to prod
12. Live on `crowdunlocked.com`

### No Manual Steps
- No manual terraform applies
- No manual deployments
- No manual approvals
- Just PR reviews (self-review OK for solo dev)

## Key Features

✅ **Two separate domains** for clean environment separation
✅ **Committed tfvars** for non-sensitive config
✅ **GitHub Secrets** for sensitive data
✅ **OIDC authentication** (no long-lived AWS credentials)
✅ **Auto-deploy on merge** (no manual steps)
✅ **Test-driven development** enforced by CI
✅ **GitOps with Flux** for Kubernetes deployments
✅ **Comprehensive documentation** for solo developer

## Security

- No sensitive data in git
- All secrets in GitHub Secrets (encrypted)
- OIDC for AWS authentication
- Terraform state in encrypted S3
- Privacy protection on domain registrations
- TLS everywhere (ACM certificates)

## Solo Developer Optimizations

- Self-review allowed
- No manual approvals
- Fast feedback with `make test-watch`
- Comprehensive docs for reference
- Automated everything
- Clear workflow to follow

## Troubleshooting

See individual documentation files:
- OpenTofu/Infrastructure issues: `infra/terraform/CICD_SETUP.md`
- Development issues: `DEVELOPMENT_WORKFLOW.md`
- Testing issues: `docs/TDD_GUIDE.md`

## Questions?

All documentation is in the repo:
- Start with `DEVELOPMENT_WORKFLOW.md`
- Reference `docs/TDD_GUIDE.md` for testing
- Check `infra/terraform/CICD_SETUP.md` for CI/CD
- Review `docs/ARCHITECTURE.md` for system design
