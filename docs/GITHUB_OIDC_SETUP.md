# GitHub OIDC Setup for CI/CD

This guide explains how to set up GitHub Actions to authenticate with AWS using OIDC (OpenID Connect) for secure, keyless authentication.

## Overview

The terraform configurations create OIDC providers and IAM roles in each AWS account (mgmt, dev, prod). GitHub Actions will assume these roles to deploy infrastructure without needing long-lived AWS credentials.

## Prerequisites

- AWS accounts created (mgmt, dev, prod)
- Terraform applied in all three environments
- GitHub repository admin access

---

## Step 1: Apply Terraform to Create OIDC Roles

The OIDC module is already included in the terraform configurations. Apply terraform in each environment:

### Management Account
```bash
cd infra/terraform/mgmt
tofu init
tofu plan -var-file=mgmt.tfvars
tofu apply -var-file=mgmt.tfvars
```

### Dev Account
```bash
cd infra/terraform/dev
tofu init
tofu plan -var-file=dev.tfvars
tofu apply -var-file=dev.tfvars
```

### Prod Account
```bash
cd infra/terraform/prod
tofu init
tofu plan -var-file=prod.tfvars
tofu apply -var-file=prod.tfvars
```

---

## Step 2: Get the Role ARNs

After applying terraform, get the role ARNs from the outputs:

```bash
# Management account
cd infra/terraform/mgmt
tofu output github_actions_role_arn

# Dev account
cd infra/terraform/dev
tofu output github_actions_role_arn

# Prod account
cd infra/terraform/prod
tofu output github_actions_role_arn
```

The ARNs will look like:
- `arn:aws:iam::111111111111:role/github-actions-mgmt`
- `arn:aws:iam::222222222222:role/github-actions-dev`
- `arn:aws:iam::333333333333:role/github-actions-prod`

---

## Step 3: Add GitHub Secrets

Add these secrets to your GitHub repository:

### Via GitHub UI

1. Go to your repository on GitHub
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add each secret:

| Secret Name | Value | Description |
|------------|-------|-------------|
| `AWS_ROLE_mgmt` | `arn:aws:iam::MGMT_ACCOUNT_ID:role/github-actions-mgmt` | OIDC role for mgmt account |
| `AWS_ROLE_dev` | `arn:aws:iam::DEV_ACCOUNT_ID:role/github-actions-dev` | OIDC role for dev account |
| `AWS_ROLE_prod` | `arn:aws:iam::PROD_ACCOUNT_ID:role/github-actions-prod` | OIDC role for prod account |

### Via GitHub CLI

```bash
# Set the role ARNs (replace with your actual ARNs)
MGMT_ROLE_ARN="arn:aws:iam::111111111111:role/github-actions-mgmt"
DEV_ROLE_ARN="arn:aws:iam::222222222222:role/github-actions-dev"
PROD_ROLE_ARN="arn:aws:iam::333333333333:role/github-actions-prod"

# Add secrets
gh secret set AWS_ROLE_mgmt --body "$MGMT_ROLE_ARN"
gh secret set AWS_ROLE_dev --body "$DEV_ROLE_ARN"
gh secret set AWS_ROLE_prod --body "$PROD_ROLE_ARN"
```

---

## Step 4: Verify Setup

Push a change to trigger the workflow:

```bash
git add .
git commit -m "Test OIDC authentication"
git push
```

Check the GitHub Actions workflow:
1. Go to **Actions** tab in your repository
2. Click on the running workflow
3. Verify the "Configure AWS Credentials" step succeeds

You should see output like:
```
Assuming role with OIDC
Role assumed successfully
```

---

## How It Works

1. **GitHub Actions requests a token** from GitHub's OIDC provider
2. **Token includes claims** about the repository, branch, and workflow
3. **AWS validates the token** against the OIDC provider
4. **AWS checks the trust policy** to ensure the request is from your repository
5. **AWS issues temporary credentials** (valid for 1 hour)
6. **Workflow uses credentials** to run terraform commands

---

## Security Benefits

- ✅ **No long-lived credentials** stored in GitHub
- ✅ **Automatic credential rotation** (expires after 1 hour)
- ✅ **Scoped to specific repository** (can't be used elsewhere)
- ✅ **Audit trail** in AWS CloudTrail
- ✅ **Can't be leaked** (tokens are short-lived and scoped)

---

## Troubleshooting

### Error: "Credentials could not be loaded"

**Cause:** GitHub secret not set or incorrect role ARN

**Fix:**
1. Verify secrets are set: `gh secret list`
2. Check role ARN format: `arn:aws:iam::ACCOUNT_ID:role/github-actions-ENV`
3. Ensure role exists in AWS: `aws iam get-role --role-name github-actions-ENV`

### Error: "Not authorized to perform sts:AssumeRoleWithWebIdentity"

**Cause:** OIDC provider not created or trust policy incorrect

**Fix:**
1. Verify OIDC provider exists:
   ```bash
   aws iam list-open-id-connect-providers
   ```
2. Check trust policy on role:
   ```bash
   aws iam get-role --role-name github-actions-ENV
   ```
3. Ensure repository name matches in trust policy

### Error: "Access denied" during terraform operations

**Cause:** IAM role doesn't have sufficient permissions

**Fix:**
1. Check the role's policies:
   ```bash
   aws iam list-role-policies --role-name github-actions-ENV
   ```
2. Update the module's policy in `infra/terraform/modules/github-oidc/main.tf`
3. Re-apply terraform

---

## Manual Deployment (Without OIDC)

If you need to deploy manually without OIDC:

```bash
# Configure AWS credentials
export AWS_PROFILE=crowdunlocked-dev

# Run terraform
cd infra/terraform/dev
tofu init
tofu plan -var-file=dev.tfvars
tofu apply -var-file=dev.tfvars
```

---

## Next Steps

After OIDC is set up:
1. ✅ Push changes to trigger automated deployments
2. ✅ Merge PRs to deploy to dev (develop branch) or prod (main branch)
3. ✅ Monitor deployments in GitHub Actions
4. ✅ Check AWS resources are created correctly

---

## References

- [GitHub OIDC Documentation](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-amazon-web-services)
- [AWS IAM OIDC Documentation](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_create_oidc.html)
- [OpenTofu AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
