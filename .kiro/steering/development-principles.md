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

#### Terraform Organization Standards
**File Structure:**
- `main.tf` - Primary resource definitions with clear sections
- `variables.tf` - Input variables with descriptions
- `outputs.tf` - Output values with descriptions
- `*.tfvars` - Environment-specific values

**main.tf Section Order:**
1. Terraform Configuration (backend, providers)
2. Data Sources (external data lookups)
3. Local Variables (computed values)
4. VPC and Networking (foundation layer)
5. IAM Roles and Policies (security layer)
6. Compute Resources (EKS, EC2, etc.)
7. Data Storage (DynamoDB, RDS, S3)
8. Container Registry (ECR)
9. Monitoring (CloudWatch, alarms)

**Required Infrastructure Components:**
- **Networking**: VPC, public/private subnets, Internet Gateway, NAT Gateways, route tables
- **Security**: IAM roles with least privilege, security groups, NACLs
- **High Availability**: Multi-AZ deployment for critical resources
- **Monitoring**: CloudWatch logs, metrics, and alarms
- **Backup**: Point-in-time recovery for databases, lifecycle policies for images

**Best Practices:**
- Use section headers with `# ===` separators for clarity
- Group related resources together
- Add descriptions to all variables and outputs
- Use `for_each` instead of `count` for resource sets
- Declare dependencies explicitly with `depends_on`
- Tag all resources with Environment, Service, ManagedBy
- Use data sources for cross-stack references
- Keep resource names consistent: `{project}-{env}-{resource}`

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
- **Using OIDC for AWS authentication** - No long-lived access keys in CI/CD

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
- Keep terraform variables documented with descriptions
- Maintain architecture diagrams
- Document networking architecture (VPC, subnets, routing)
- Create runbooks for common infrastructure tasks
- Keep a changelog of infrastructure changes

## Current Status

### ✅ Implemented
- TDD for venue search feature
- Infrastructure as code (terraform)
- GitHub Actions CI/CD pipelines with OIDC authentication
- Separate dev/prod environments
- Automated testing in CI/CD
- ECR repositories for container images
- DynamoDB tables for all services

### 🚧 In Progress
- EKS cluster deployment
- Certificate management
- Google Places API integration

### 📋 TODO
- E2E testing framework
- Performance testing
- Security scanning in CI/CD
- Automated dependency updates
- Clean up old dev account (987470856210)

## Infrastructure Review Checklist

Before merging infrastructure changes, verify:

### Networking
- [ ] VPC has DNS support and hostnames enabled
- [ ] Public subnets exist for load balancers and NAT gateways
- [ ] Private subnets exist for application workloads
- [ ] Internet Gateway attached to VPC
- [ ] NAT Gateways in each AZ for high availability
- [ ] Route tables properly configured (public → IGW, private → NAT)
- [ ] Subnets tagged for Kubernetes ELB integration

### Security
- [ ] IAM roles follow least privilege principle
- [ ] Security groups restrict access appropriately
- [ ] No hardcoded credentials or secrets
- [ ] Encryption enabled for data at rest
- [ ] Encryption enabled for data in transit

### High Availability
- [ ] Resources distributed across multiple AZs
- [ ] Auto-scaling configured where appropriate
- [ ] Health checks configured
- [ ] Backup and recovery procedures in place

### Monitoring
- [ ] CloudWatch logs configured
- [ ] Alarms set for critical metrics
- [ ] Log retention policies defined
- [ ] Metrics exported for key resources

### Cost Optimization
- [ ] Right-sized instance types
- [ ] Lifecycle policies for ephemeral data
- [ ] Reserved capacity considered for stable workloads
- [ ] Unused resources cleaned up

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
