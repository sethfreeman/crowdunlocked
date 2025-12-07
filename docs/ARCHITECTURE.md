# Crowd Unlocked Architecture

## Overview

Crowd Unlocked is a cloud-native artist management platform built on AWS with a microservices architecture.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        CloudFront CDN                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Application Load Balancer                 │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│   Bookings   │    │   Releases   │    │   Publicity  │
│   Service    │    │   Service    │    │   Service    │
└──────────────┘    └──────────────┘    └──────────────┘
        │                     │                     │
        ▼                     ▼                     ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  DynamoDB    │    │  DynamoDB    │    │  DynamoDB    │
└──────────────┘    └──────────────┘    └──────────────┘
```

## Microservices

### 1. Bookings Service
- Manages artist bookings and event scheduling
- DynamoDB table: `bookings-{env}`
- Port: 8080

### 2. Releases Service
- Tracks music releases and distribution
- DynamoDB table: `releases-{env}`
- Port: 8080

### 3. Publicity Service
- Handles PR campaigns and media outreach
- DynamoDB table: `publicity-{env}`
- Port: 8080

### 4. Social Service
- Monitors social media presence
- DynamoDB table: `social-{env}`
- Port: 8080

### 5. Money Service
- Manages revenue streams and payments
- DynamoDB table: `money-{env}`
- Port: 8080

## Infrastructure

### AWS Accounts
- **Management**: Organizations, SSO, Route 53, ACM
- **Dev**: Development EKS cluster, DynamoDB tables
- **Prod**: Production EKS cluster, DynamoDB tables, CloudFront

### EKS Auto Mode
- Fargate-based compute (serverless)
- Automatic scaling and patching
- No node management required

### GitOps with Flux
- `develop` branch → Dev cluster
- `main` branch → Prod cluster
- Automatic reconciliation every 10 minutes

### Observability
- **Tracing**: AWS X-Ray
- **Logs**: CloudWatch Logs
- **Metrics**: CloudWatch Metrics
- **Alarms**: Error rate monitoring

## Data Storage

All services use DynamoDB with:
- Pay-per-request billing
- Point-in-time recovery enabled
- Streams enabled for event-driven patterns

## Security

- IAM roles for service accounts (IRSA)
- Private subnets for EKS
- TLS everywhere (ACM certificates)
- AWS Organizations for account isolation
- SSO for human access

## Frontend Applications

### Web (Next.js)
- Server-side rendering
- Deployed to EKS
- CloudFront CDN distribution

### Mobile (Flutter)
- iOS and Android support
- Native performance
- Shared codebase
