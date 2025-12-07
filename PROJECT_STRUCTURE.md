# Crowd Unlocked - Project Structure

## Overview

Enterprise-grade monorepo for Crowd Unlocked artist management platform with Go microservices, Next.js web app, Flutter mobile app, and complete AWS infrastructure.

## Directory Structure

```
crowdunlocked/
├── services/                    # Go microservices
│   ├── bookings/               # Bookings service (TDD-ready)
│   │   ├── cmd/server/         # Main entry point
│   │   ├── internal/
│   │   │   ├── domain/         # Business logic + tests
│   │   │   ├── repository/     # DynamoDB repository
│   │   │   └── handler/        # HTTP handlers
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── releases/               # Music releases service
│   ├── publicity/              # PR and publicity service
│   ├── social/                 # Social media service
│   └── money/                  # Revenue management service
│
├── apps/
│   ├── web/                    # Next.js web application
│   │   ├── app/                # App router pages
│   │   ├── Dockerfile
│   │   ├── package.json
│   │   └── next.config.js
│   └── mobile/                 # Flutter mobile app
│       ├── lib/main.dart
│       ├── android/
│       ├── ios/
│       └── pubspec.yaml
│
├── infra/
│   └── terraform/
│       ├── management/         # AWS Organizations, SSO, Route 53, ACM
│       ├── dev/                # Dev environment (EKS, DynamoDB)
│       └── prod/               # Prod environment (EKS, CloudFront, DynamoDB)
│
├── k8s/                        # Kubernetes manifests
│   ├── base/                   # Base configurations
│   │   ├── bookings/
│   │   ├── releases/
│   │   └── xray/              # AWS X-Ray daemon
│   └── overlays/
│       ├── dev/               # Dev-specific configs
│       └── prod/              # Prod-specific configs
│
├── flux/                       # Flux CD GitOps
│   └── clusters/
│       ├── dev/               # Watches 'develop' branch
│       └── prod/              # Watches 'main' branch
│
├── docs/
│   ├── ARCHITECTURE.md        # System architecture
│   └── SETUP.md              # Setup instructions
│
├── scripts/
│   ├── bootstrap.sh          # Initial setup script
│   └── create-dynamodb-tables.sh
│
├── .github/workflows/
│   └── ci.yaml               # CI/CD pipeline
│
├── go.work                    # Go workspace
├── Makefile                   # Build automation
├── docker-compose.yml         # Local development
└── README.md
```

## Technology Stack

### Backend Services
- **Language**: Go 1.22
- **Framework**: Chi router
- **Database**: DynamoDB (all services)
- **Tracing**: AWS X-Ray
- **Testing**: TDD with testify

### Frontend
- **Web**: Next.js 14 (React, TypeScript, Tailwind CSS)
- **Mobile**: Flutter 3.2+ (iOS & Android)

### Infrastructure
- **Cloud**: AWS (multi-account setup)
- **IaC**: Terraform
- **Container Orchestration**: EKS Auto Mode (Fargate)
- **GitOps**: Flux CD
- **CDN**: CloudFront
- **DNS**: Route 53
- **Certificates**: ACM
- **Monitoring**: CloudWatch + X-Ray

### AWS Account Structure
1. **Management**: Organizations, SSO, shared services
2. **Dev**: Development environment
3. **Prod**: Production environment

## Key Features

### Microservices Architecture
- 5 independent Go services
- Each with dedicated DynamoDB table
- Containerized with Docker
- Deployed to EKS Fargate

### GitOps Workflow
- `develop` branch → Auto-deploy to dev cluster
- `main` branch → Auto-deploy to prod cluster
- Flux reconciles every 10 minutes

### Observability
- Distributed tracing with X-Ray
- Centralized logging with CloudWatch
- Automated alarms for error rates

### Security
- IAM roles for service accounts (IRSA)
- Private EKS subnets
- TLS everywhere
- Multi-account isolation
- SSO for human access

## Quick Start

```bash
# 1. Bootstrap project
./scripts/bootstrap.sh

# 2. Run tests
make test

# 3. Local development
docker-compose up

# 4. Deploy infrastructure
cd infra/terraform/mgmt
terraform init && terraform apply

# 5. Setup GitOps
flux bootstrap github --owner=crowdunlocked --repository=crowdunlocked
```

## Development Workflow

1. **Feature Development**: Work on `develop` branch
2. **Testing**: TDD - write tests first, then implementation
3. **Local Testing**: Use docker-compose for integration tests
4. **Push**: Commit triggers CI pipeline
5. **Auto-Deploy**: Flux deploys to dev cluster
6. **Production**: Merge to `main` for prod deployment

## CI/CD Pipeline

- **Test**: Run Go tests for all services
- **Build**: Build Docker images
- **Push**: Push to Amazon ECR
- **Deploy**: Flux auto-deploys from Git

## Monitoring & Alerts

- CloudWatch Logs: `/aws/eks/crowdunlocked-{env}/{service}`
- X-Ray Traces: Distributed tracing across services
- CloudWatch Alarms: High error rate detection
- Metrics: Request counts, latencies, errors

## Next Steps

1. Configure AWS credentials
2. Set domain name in Terraform variables
3. Run infrastructure deployment
4. Bootstrap Flux on clusters
5. Push code to trigger deployments

See `docs/SETUP.md` for detailed instructions.
