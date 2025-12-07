# Crowd Unlocked

[![Tests](https://github.com/sethfreeman/crowdunlocked/actions/workflows/services.yaml/badge.svg)](https://github.com/sethfreeman/crowdunlocked/actions/workflows/services.yaml)
[![Infrastructure](https://github.com/sethfreeman/crowdunlocked/actions/workflows/terraform.yaml/badge.svg)](https://github.com/sethfreeman/crowdunlocked/actions/workflows/terraform.yaml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22-blue.svg)](https://golang.org)
[![OpenTofu](https://img.shields.io/badge/OpenTofu-1.10-purple.svg)](https://opentofu.org)

> Open-source artist management platform built with microservices, Next.js, Flutter, and AWS

**Enterprise monorepo for Crowd Unlocked platform.**

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

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](./docs)
- 🚀 [Quick Start](./QUICKSTART.md)
- 💬 [Discussions](https://github.com/sethfreeman/crowdunlocked/discussions)
- 🐛 [Issues](https://github.com/sethfreeman/crowdunlocked/issues)

---

**Built with ❤️ for artists and their teams**
