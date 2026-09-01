# Vercel Deployment Guide (Next.js + DynamoDB via AWS OIDC)

This guide walks through deploying the Crowd Unlocked web app (Next.js, including
its API routes) to Vercel, with the API routes reading/writing DynamoDB using
keyless AWS OIDC authentication.

This is the current architecture. The app is a single Next.js project in
`apps/web`; its API routes under `pages/api/v1/**` talk to DynamoDB directly.
There are no separate backend services to deploy.

## Prerequisites

1. **AWS infrastructure already applied.** The dev terraform (`infra/terraform/dev`)
   provisions the DynamoDB tables and the Vercel OIDC module (provider + IAM role
   scoped to the venues and bookings tables). This has already been applied; you
   only re-apply after adding your Vercel project IDs (Step 3).
2. **Vercel account.** Sign up at [vercel.com](https://vercel.com). The Hobby
   (free) tier is sufficient for this project.
3. **GitHub repository.** Your code is already at `sethfreeman/crowdunlocked`.

## Region note (important)

The dev DynamoDB tables live in **us-west-2**. The only thing in us-east-1 is the
terraform state-lock table. Set `AWS_REGION=us-west-2` in Vercel and locally, or
DynamoDB calls will hit the wrong region and fail.

## Step 1: Confirm the AWS infrastructure

The Vercel OIDC role and DynamoDB tables already exist. To confirm and grab the
role ARN you'll need later:

```bash
cd infra/terraform/dev
AWS_PROFILE=crowdunlocked-dev tofu output vercel_role_arn
# arn:aws:iam::179151668767:role/crowdunlocked-dev-vercel-role
```

The OIDC module created:
- An OIDC provider for `oidc.vercel.com`
- An IAM role (`crowdunlocked-dev-vercel-role`) that Vercel assumes via web identity
- A DynamoDB policy scoped to the `venues-dev` and `bookings-dev` tables (plus their indexes)

## Step 2: Create the Vercel project

1. Go to [vercel.com/dashboard](https://vercel.com/dashboard) → **Add New → Project**.
2. Import the `crowdunlocked` GitHub repository.
3. **Set the Root Directory to `apps/web`.** This is the monorepo web app.
4. Vercel auto-detects Next.js. Defaults are correct:
   - Framework Preset: **Next.js**
   - Build Command: `npm run build`
   - Output Directory: `.next`
5. Deploy. The first build will succeed but the API routes won't reach DynamoDB
   yet, that's expected until Steps 3-4 are done.

## Step 3: Wire the Vercel project IDs into terraform

The OIDC role only trusts specific Vercel projects. After the project exists,
grab its IDs and add them to the trust condition.

1. In Vercel: **Project Settings → General**, note the **Project ID**, and
   **Account/Team Settings** for the **Team ID**.
2. Edit `infra/terraform/dev/main.tf`, in the `module "vercel_oidc"` block, set:

   ```hcl
   vercel_project_ids = [
     "team_YOUR_TEAM_ID:project_YOUR_PROJECT_ID:environment_production",
     "team_YOUR_TEAM_ID:project_YOUR_PROJECT_ID:environment_preview",
   ]
   ```

   Include the `preview` entry if you want PR preview deployments to reach AWS too.
3. Apply:

   ```bash
   cd infra/terraform/dev
   AWS_PROFILE=crowdunlocked-dev tofu apply -var-file=dev.tfvars
   ```

## Step 4: Configure Vercel environment variables

In **Project Settings → Environment Variables**, add:

```bash
AWS_REGION=us-west-2
DYNAMODB_VENUES_TABLE=venues-dev
DYNAMODB_BOOKINGS_TABLE=bookings-dev

# The OIDC role Vercel assumes (from Step 1)
AWS_ROLE_ARN=arn:aws:iam::179151668767:role/crowdunlocked-dev-vercel-role
```

When `AWS_ROLE_ARN` is set and Vercel's AWS OIDC integration is enabled, the AWS
SDK inside your API routes exchanges the Vercel OIDC token for temporary AWS
credentials automatically. No access keys are stored anywhere.

Redeploy so the new environment variables take effect (Deployments → redeploy, or
push a commit).

## Step 5: Verify

```bash
# Health check
curl https://YOUR-APP.vercel.app/api/health

# Create a venue
curl -X POST https://YOUR-APP.vercel.app/api/v1/venues \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Venue",
    "location": {"latitude": 40.7128, "longitude": -74.0060},
    "address": {"street": "123 Main St", "city": "New York", "state": "NY", "postal_code": "10001", "country": "US"},
    "venue_types": ["club"]
  }'

# Search venues
curl "https://YOUR-APP.vercel.app/api/v1/venues/search?venue_types=club&limit=5"
```

## Step 6: Custom domain (optional)

The `crowdunlockedbeta.com` DNS hosted zone lives in AWS account `021645491430`
(the account that also holds the terraform state and the `default` CLI profile).
Vercel issues and manages TLS for custom domains itself, so you do **not** need an
ACM certificate.

1. In Vercel: **Project Settings → Domains → Add**, enter `crowdunlockedbeta.com`.
2. Vercel shows the DNS records to add (A/CNAME).
3. Add those records in the Route 53 hosted zone for `crowdunlockedbeta.com`.
4. Vercel validates and provisions TLS automatically.

## Deployments

Vercel deploys automatically from GitHub:
- Push/merge to a branch with the project connected → production or preview deploy.
- No GitHub Actions deploy job is needed; CI (`.github/workflows/ci.yaml`) only
  lints, tests, type-checks, and builds the web app, plus runs a plan-only
  terraform check on PRs.

## Troubleshooting

- **403 / AccessDenied assuming the role:** the Vercel project IDs in
  `vercel_project_ids` don't match the deployment's team/project/environment.
  Re-check Step 3 and re-apply.
- **DynamoDB ResourceNotFound / timeouts:** wrong region. Confirm
  `AWS_REGION=us-west-2` (see the region note above).
- **Credentials errors:** confirm `AWS_ROLE_ARN` is set and Vercel's AWS OIDC
  integration is enabled for the project.
- **Check logs:** Vercel dashboard → Functions/Logs for the API routes.

## Cost

- **Vercel Hobby tier:** free (fine for alpha / no-traffic).
- **DynamoDB (on-demand):** pay-per-request, effectively ~$0 when idle.
- **AWS OIDC:** free.
- **Domain:** ~$12/year if you keep `crowdunlockedbeta.com` registered.

Total ongoing hosting: roughly $0/month at current (alpha, no traffic) usage.

## Security notes

- No long-lived AWS access keys; Vercel uses OIDC web-identity federation.
- The IAM role is least-privilege: it can only touch the venues and bookings tables.
- To scope things down further, you can split production vs preview into separate
  role conditions.
