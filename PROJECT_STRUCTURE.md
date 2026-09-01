# Crowd Unlocked - Project Structure

## Overview

Monorepo for the Crowd Unlocked artist management platform. The web app and its
API are a single Next.js project deployed to Vercel; data lives in DynamoDB. A
Flutter mobile app shares the same backend.

> Historical note: this repo previously contained Go microservices, Kubernetes
> manifests (`k8s/`), Flux GitOps config (`flux/`), and an EKS deployment. Those
> were removed to reduce cost and operational overhead. The application logic now
> lives entirely in the Next.js API routes. Prior structure is in git history.

## Directory Structure

```
crowdunlocked/
├── apps/
│   ├── web/                    # Next.js web app + API
│   │   ├── app/                # App Router pages (UI)
│   │   ├── pages/api/          # API routes
│   │   │   ├── health.ts
│   │   │   └── v1/
│   │   │       ├── venues/     # index, [id], search
│   │   │       └── bookings/   # index
│   │   ├── lib/                # types, utils (geohash)
│   │   ├── __tests__/          # Jest API route tests
│   │   └── package.json
│   └── mobile/                 # Flutter app (iOS & Android)
│
├── infra/
│   └── terraform/
│       ├── mgmt/               # AWS Organization, domains, ACM certs
│       ├── dev/                # DynamoDB tables + Vercel OIDC role
│       └── modules/
│           ├── vercel-oidc/    # OIDC provider + IAM role for Vercel
│           └── github-oidc/    # OIDC role for GitHub Actions
│
├── docs/                       # Architecture and setup guides
├── scripts/                    # Helper scripts
└── .github/workflows/ci.yaml   # Lint, test, type-check, build, tf plan
```

## Technology Stack

### Web + API
- **Framework**: Next.js 14 (React, TypeScript)
- **API**: Next.js API routes (`pages/api/v1/**`)
- **AWS SDK**: `@aws-sdk/client-dynamodb`, `@aws-sdk/lib-dynamodb`
- **Testing**: Jest + node-mocks-http
- **Hosting**: Vercel (serverless functions for API routes)

### Mobile
- **Flutter** 3.2+ (iOS & Android)

### Data
- **DynamoDB** (pay-per-request), region `us-west-2`
  - `venues-dev` with GSIs: GeohashIndex, CityIndex, VenueTypeIndex, ExternalIdIndex
  - `bookings-dev` (streams enabled)
  - `releases-dev`, `publicity-dev`, `social-dev`, `money-dev` (currently unused)

### Infrastructure
- **IaC**: OpenTofu (`infra/terraform`)
- **State**: S3 bucket `crowdunlocked-terraform-state` + DynamoDB lock (us-east-1)
- **Auth**: Vercel → AWS via OIDC; GitHub Actions → AWS via OIDC (no long-lived keys)

## Data Model Highlights

The `venues` table drives search via secondary indexes:
- **GeohashIndex** — spatial/location search
- **CityIndex** — city/state lookups
- **VenueTypeIndex** — filter by venue type
- **ExternalIdIndex** — dedupe against external sources

The venues API route writes the geohash, `city_state`, and GSI key attributes on
create/update so these queries work.

## Deployment

- **Web app**: Vercel, via its GitHub integration. See `docs/VERCEL_DEPLOYMENT.md`.
- **Infrastructure**: OpenTofu applied from `infra/terraform/dev`.

## Development Workflow

1. Feature branch from `develop`.
2. Make changes and update tests (`apps/web/__tests__`).
3. Run `npm test`, `npm run lint`, `npm run type-check`, `npm run build`.
4. PR to `develop`; CI validates. Vercel deploys on merge.

See `docs/VERCEL_DEPLOYMENT.md` for deployment and `docs/ARCHITECTURE.md` for design.
