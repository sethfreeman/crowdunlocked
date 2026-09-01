# Crowd Unlocked

[![CI](https://github.com/sethfreeman/crowdunlocked/actions/workflows/ci.yaml/badge.svg)](https://github.com/sethfreeman/crowdunlocked/actions/workflows/ci.yaml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> Artist management platform. Next.js web app + API on Vercel, backed by DynamoDB.

## Architecture

- **Web + API**: Next.js 14 (`apps/web`). API routes under `pages/api/v1/**` talk
  to DynamoDB directly. Deployed to Vercel.
- **Database**: DynamoDB (pay-per-request), tables in `us-west-2`.
- **Auth to AWS**: Vercel serverless functions assume an AWS IAM role via OIDC
  (no long-lived keys).
- **Infrastructure**: OpenTofu (`infra/terraform/dev`) manages the DynamoDB tables
  and the Vercel OIDC role.
- **Mobile**: Flutter app (`apps/mobile`).

> Note: this project previously ran on EKS + Go microservices + Flux GitOps. That
> architecture was removed to cut cost and complexity; the app now lives entirely
> in the Next.js project. See git history if you need the old setup.

## Structure

```
crowdunlocked/
├── apps/
│   ├── web/            # Next.js web app + API routes (deployed to Vercel)
│   └── mobile/         # Flutter app
├── infra/terraform/
│   ├── dev/            # DynamoDB tables + Vercel OIDC role
│   ├── mgmt/           # AWS Organization, domains, certificates
│   └── modules/        # Reusable modules (vercel-oidc, github-oidc)
└── docs/               # Architecture and setup guides
```

## Quick Start (local web dev)

```bash
cd apps/web
npm install
npm run dev          # http://localhost:3000
```

The API routes read AWS config from environment variables (see `apps/web/.env.local`).
For local development against the real dev tables, use an AWS profile with access to
the `venues-dev` / `bookings-dev` tables in `us-west-2`.

## Testing

```bash
cd apps/web
npm test             # Jest (API route tests)
npm run lint
npm run type-check
npm run build
```

## Deployment

The web app deploys to **Vercel** via its GitHub integration. See
**[docs/VERCEL_DEPLOYMENT.md](docs/VERCEL_DEPLOYMENT.md)** for the full setup
(create project, wire OIDC, set env vars, custom domain).

Infrastructure changes go through OpenTofu:

```bash
cd infra/terraform/dev
AWS_PROFILE=crowdunlocked-dev tofu plan -var-file=dev.tfvars
```

## Documentation

- **[Vercel Deployment](docs/VERCEL_DEPLOYMENT.md)** — deploy the app (start here)
- **[Architecture](docs/ARCHITECTURE.md)** — system design
- **[TDD Guide](docs/TDD_GUIDE.md)** — testing practices
- **[AWS Organization Setup](docs/AWS_ORGANIZATION_SETUP.md)** — one-time account setup
- **[GitHub OIDC Setup](docs/GITHUB_OIDC_SETUP.md)** — CI auth to AWS

## Development Workflow

1. Create a feature branch from `develop`.
2. Make changes; add/adjust tests.
3. `npm test && npm run lint && npm run type-check && npm run build` in `apps/web`.
4. Push and open a PR to `develop`. CI runs lint/test/type-check/build and a
   terraform plan check.
5. Merge. Vercel deploys automatically.

## License

MIT — see [LICENSE](LICENSE).
