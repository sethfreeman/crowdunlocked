# Dev infrastructure (OpenTofu)

Manages the dev DynamoDB tables and the Vercel OIDC role. The web app itself is
deployed to Vercel, not to AWS compute (see `docs/VERCEL_DEPLOYMENT.md`).

## Gotchas (read this before running tofu)

### 1. Use the native arm64 tofu binary

On Apple Silicon, the Homebrew Intel `tofu` (`/usr/local/bin/tofu`) runs under
Rosetta and the large AWS provider times out during plugin startup. Use the
arm64 build:

```bash
/opt/homebrew/bin/tofu
```

If `which tofu` points at `/usr/local/bin/tofu`, call the full path above (or put
`/opt/homebrew/bin` earlier on your PATH).

### 2. Two AWS profiles: backend vs provider

- **Backend (state):** S3 bucket `crowdunlocked-terraform-state` lives in the
  account behind the **`default`** profile (021645491430).
- **Provider (resources):** the DynamoDB tables and IAM role live in the dev
  account (179151668767), reached via the **`crowdunlocked-dev`** profile.

The backend profile is pinned via `backend.hcl`. The provider profile is passed
at runtime with `AWS_PROFILE`.

## Commands

Init (once per machine / after backend changes):

```bash
AWS_PROFILE=default /opt/homebrew/bin/tofu init -reconfigure -backend-config=backend.hcl
```

Plan / apply (provider = dev account):

```bash
AWS_PROFILE=crowdunlocked-dev /opt/homebrew/bin/tofu plan  -var-file=dev.tfvars
AWS_PROFILE=crowdunlocked-dev /opt/homebrew/bin/tofu apply -var-file=dev.tfvars
```

Read outputs (state read = backend/default profile):

```bash
AWS_PROFILE=default /opt/homebrew/bin/tofu output
```

## Key outputs

- `vercel_role_arn` = `arn:aws:iam::179151668767:role/crowdunlocked-dev-vercel-role`
- `vercel_oidc_provider_arn` = `arn:aws:iam::179151668767:oidc-provider/oidc.vercel.com`
- `dynamodb_tables` = map of the six `*-dev` table names (region **us-west-2**)

## Vercel project IDs

After creating the Vercel project, add its IDs to `vercel_project_ids` in
`main.tf` (module "vercel_oidc") and re-apply so the OIDC role trusts it:

```hcl
vercel_project_ids = [
  "team_YOUR_TEAM_ID:project_YOUR_PROJECT_ID:environment_production",
  "team_YOUR_TEAM_ID:project_YOUR_PROJECT_ID:environment_preview",
]
```
