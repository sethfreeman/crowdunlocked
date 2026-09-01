# Crowd Unlocked Architecture

## Overview

Crowd Unlocked is an artist management platform. The web experience and its API
are a single Next.js application deployed to Vercel; persistent data lives in
DynamoDB. A Flutter mobile app shares the same API.

> Historical note: an earlier design used Go microservices on EKS with Flux
> GitOps and CloudFront. That was removed to cut cost and complexity. The
> application logic now lives in the Next.js API routes. This document describes
> the current system.

## System Diagram

```
                 ┌─────────────────────────────┐
                 │        Vercel (CDN +         │
                 │     serverless functions)    │
                 │                              │
   Browser ─────▶│   Next.js app (apps/web)     │
   Flutter ─────▶│   UI (App Router)            │
                 │   API routes (pages/api/v1)  │
                 └───────────────┬──────────────┘
                                 │ AWS SDK (OIDC role)
                                 ▼
                 ┌─────────────────────────────┐
                 │          DynamoDB            │
                 │   venues-dev / bookings-dev  │
                 │        (us-west-2)           │
                 └─────────────────────────────┘
```

## Components

### Web + API (Next.js, `apps/web`)
- **UI**: App Router pages under `app/`.
- **API**: Next.js API routes under `pages/api/v1/**`:
  - `venues/index.ts` — list/create venues
  - `venues/[id].ts` — get/update/delete a venue
  - `venues/search.ts` — geohash / city / type search
  - `bookings/index.ts` — list/create bookings (validates the venue exists)
  - `health.ts` — health check
- Each route uses the AWS SDK to query DynamoDB directly. No intermediate service tier.
- Deployed to Vercel; API routes run as serverless functions.

### Data (DynamoDB)
- Region: `us-west-2`. Billing: pay-per-request (scales to ~zero when idle).
- `venues-dev`: primary key `id`, with GSIs for spatial (geohash), city/state,
  venue type, and external-source dedupe. The create/update routes maintain the
  geohash, `city_state`, and GSI key attributes.
- `bookings-dev`: primary key `id`, DynamoDB Streams enabled.
- `releases-dev`, `publicity-dev`, `social-dev`, `money-dev`: provisioned but
  currently unused (reserved for future features).

### AWS authentication (OIDC)
- Vercel serverless functions assume an AWS IAM role
  (`crowdunlocked-dev-vercel-role`) via OIDC web-identity federation. No
  long-lived access keys.
- The role's policy is scoped to the venues and bookings tables (and their indexes).
- Defined in `infra/terraform/modules/vercel-oidc`, wired up in
  `infra/terraform/dev`.

### Infrastructure (OpenTofu)
- `infra/terraform/dev`: DynamoDB tables + the Vercel OIDC module.
- `infra/terraform/mgmt`: AWS Organization, domains, ACM certs (out of the app's
  hot path; TLS for the app is handled by Vercel).
- State in S3 (`crowdunlocked-terraform-state`) with a DynamoDB lock table in us-east-1.

### Mobile (Flutter)
- iOS and Android client that consumes the same API.

## Request Flow

1. A client calls `https://<app>/api/v1/...`.
2. Vercel runs the corresponding Next.js API route as a serverless function.
3. The function obtains temporary AWS credentials via the OIDC role and queries
   DynamoDB in `us-west-2`.
4. The route validates input, performs the DynamoDB operation, and returns JSON.

## Deployment

- **App**: pushed to GitHub, deployed by Vercel's GitHub integration (production
  and preview). See `docs/VERCEL_DEPLOYMENT.md`.
- **Infrastructure**: OpenTofu applied from `infra/terraform/dev`.
- **CI** (`.github/workflows/ci.yaml`): lint, test, type-check, build the web app,
  plus a plan-only terraform check on PRs. CI does not deploy; Vercel does.

## Security

- OIDC for both Vercel→AWS and GitHub Actions→AWS (no static credentials).
- Least-privilege IAM (role limited to the required tables).
- Secrets kept out of git; runtime config via Vercel environment variables.

## Cost Posture

- Vercel Hobby tier: free at current usage.
- DynamoDB on-demand: ~$0 when idle.
- Net ongoing hosting cost is roughly $0/month for the current alpha stage, plus
  domain registration if retained.
