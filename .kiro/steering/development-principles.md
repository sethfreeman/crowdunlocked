---
inclusion: always
---

# Development Principles for Crowd Unlocked

## Core Principles

### 1. Test-Driven Development (TDD)
- **Write tests first, then implementation**
- All new features must have tests before implementation
- Test coverage should be comprehensive (unit, integration, e2e where appropriate)
- Tests must pass before merging to develop or main
- Example: The venue search feature was built entirely with TDD

### 2. Infrastructure as Code (IaC)
- All infrastructure must be defined in code (Terraform/OpenTofu)
- No manual infrastructure changes in production
- Infrastructure changes go through PR review process
- State is managed in S3 with DynamoDB locking

### 3. CI/CD First
- **Do not push features if CI/CD is broken**
- All code must pass through automated pipelines
- Deployments are automated via GitHub Actions
- No manual deployments to production
- Feature branches must pass all checks before merging

### 4. Twelve-Factor App Methodology
Following https://12factor.net principles:

1. **Codebase**: One codebase tracked in git, many deploys
2. **Dependencies**: Explicitly declare and isolate dependencies (go.mod, package.json, pubspec.yaml)
3. **Config**: Store config in environment variables (never commit secrets)
4. **Backing Services**: Treat backing services as attached resources (DynamoDB, S3, etc.)
5. **Build, Release, Run**: Strictly separate build and run stages
6. **Processes**: Execute the app as stateless processes
7. **Port Binding**: Export services via port binding
8. **Concurrency**: Scale out via the process model
9. **Disposability**: Maximize robustness with fast startup and graceful shutdown
10. **Dev/Prod Parity**: Keep development, staging, and production as similar as possible
11. **Logs**: Treat logs as event streams (CloudWatch)
12. **Admin Processes**: Run admin/management tasks as one-off processes

## Workflow Rules

### Feature Development
1. Create feature branch from `develop`
2. Write tests first (TDD)
3. Implement feature to pass tests
4. Ensure all tests pass locally
5. Push and create PR
6. Wait for CI/CD checks to pass
7. Get code review
8. Merge to `develop` (auto-deploys to dev environment)

### Infrastructure Changes
1. Create feature/fix branch
2. Make infrastructure changes in terraform
3. Test locally with `tofu plan`
4. Push and create PR
5. Review terraform plan in CI/CD
6. Merge to `develop` (auto-applies to dev)
7. Verify in dev environment
8. Merge to `main` for production deployment

### Deployment Strategy
- **develop branch** → auto-deploys to dev environment
- **main branch** → auto-deploys to prod environment
- All deployments go through CI/CD
- No direct pushes to develop or main (use PRs)

## Code Quality Standards

### Go Services
- Follow standard Go project layout
- Use `gofmt` for formatting
- Run `go vet` and `golangci-lint`
- Minimum 80% test coverage
- Use table-driven tests where appropriate

### TypeScript/React (Web)
- Use TypeScript strict mode
- Follow ESLint rules
- Use Prettier for formatting
- Component tests with React Testing Library
- E2E tests with Playwright

### Flutter (Mobile)
- Follow Flutter style guide
- Use `flutter analyze`
- Widget tests for all components
- Integration tests for critical flows

## Security Practices

### Secrets Management
- **Never commit secrets to git**
- Use AWS Secrets Manager for application secrets
- Use GitHub Secrets for CI/CD credentials
- Rotate secrets regularly
- **Temporary:** Using AWS access keys for CI/CD (github-actions-ci user in each account)
- **Goal:** Switch to OIDC for AWS authentication (no long-lived keys) once GitHub Actions issue is resolved

### API Keys
- Store in AWS Secrets Manager
- Separate keys for dev and prod
- Dev keys in dev account, prod keys in prod account
- Never share keys across environments

## Documentation Requirements

### Code Documentation
- Public functions/methods must have doc comments
- Complex logic should have inline comments
- README in each service directory

### API Documentation
- OpenAPI/Swagger specs for all HTTP APIs
- Keep docs in sync with implementation
- Include examples and error responses

### Infrastructure Documentation
- Document manual setup steps (like AWS Organization)
- Keep terraform variables documented
- Maintain architecture diagrams

## Current Status

### ✅ Implemented
- TDD for venue search feature
- Infrastructure as code (terraform)
- GitHub Actions CI/CD pipelines
- Separate dev/prod environments
- Automated testing in CI/CD

### 🚧 In Progress
- Certificate management
- Google Places API integration

### 📋 TODO
- **Fix OIDC authentication** - Currently using access keys temporarily due to GitHub Actions secret interpolation issues
- E2E testing framework
- Performance testing
- Security scanning in CI/CD
- Automated dependency updates

## When CI/CD is Broken

**STOP and fix it before continuing feature work.**

If CI/CD is broken:
1. Create a fix branch
2. Fix the issue
3. Test locally
4. Push and verify in CI/CD
5. Merge fix
6. Resume feature work

Do not work around broken CI/CD. Do not merge PRs with failing checks.

## Questions?

If you're unsure about any of these principles or how to apply them, ask before proceeding. It's better to clarify than to build the wrong thing.
