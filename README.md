# Crowd Unlocked

Enterprise monorepo for Crowd Unlocked platform.

## Architecture

- **Microservices**: 5 Go services (bookings, releases, publicity, social, money)
- **Database**: DynamoDB across all services
- **Container Orchestration**: EKS Auto Mode on Fargate
- **GitOps**: Flux CD (develop → dev, main → prod)
- **Infrastructure**: Terraform-managed AWS resources
- **Frontend**: Next.js web app
- **Mobile**: Flutter app

## Structure

```
crowdunlocked/
├── services/          # Go microservices
├── apps/             # Frontend applications
├── infra/   # Terraform IaC
├── k8s/             # Kubernetes manifests
└── flux/            # Flux CD configuration
```

## AWS Accounts

- **Management**: Organizations, SSO, shared services
- **Dev**: Development environment
- **Prod**: Production environment

## Quick Start

```bash
# Run all tests
make test

# Run tests in watch mode (TDD)
make test-watch

# Format code
make fmt

# Build services
make build
```

## Documentation

- **[Initial Setup Summary](INITIAL_SETUP_SUMMARY.md)** - ⭐ Start here for overview
- **[Development Workflow](DEVELOPMENT_WORKFLOW.md)** - Daily development guide
- **[TDD Guide](docs/TDD_GUIDE.md)** - Test-driven development practices
- **[Architecture](docs/ARCHITECTURE.md)** - System design and components
- **[Setup](docs/SETUP.md)** - Initial setup instructions
- **[CI/CD Setup](infra/terraform/CICD_SETUP.md)** - GitHub Actions and OpenTofu deployment

## Development Workflow

1. Create feature branch from `develop`
2. Write tests first (TDD)
3. Implement feature
4. Create PR to `develop`
5. Merge → Auto-deploys to dev
6. Test on dev environment
7. Merge `develop` to `main` → Auto-deploys to prod

**No manual approvals needed** - Just PR review

## Environments

- **Dev**: https://crowdunlockedbeta.com (auto-deploy from `develop`)
- **Prod**: https://crowdunlocked.com (auto-deploy from `main`)
