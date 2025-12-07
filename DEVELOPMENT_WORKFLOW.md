# Development Workflow

Quick reference for daily development on Crowd Unlocked.

## Daily Workflow

### 1. Start New Feature

```bash
# Update develop branch
git checkout develop
git pull origin develop

# Create feature branch
git checkout -b feature/booking-cancellation

# Or for bug fixes
git checkout -b fix/booking-date-validation
```

### 2. Write Tests First (TDD)

```bash
# Run tests in watch mode for rapid feedback
make test-watch

# Or run tests manually
make test-unit
```

**Write failing test → Implement → Make it pass → Refactor**

See [TDD_GUIDE.md](docs/TDD_GUIDE.md) for detailed examples.

### 3. Implement Feature

```bash
# Format code as you go
make fmt

# Run linters
make lint

# Run all tests
make test
```

### 4. Commit Changes

```bash
git add .
git commit -m "feat(bookings): add booking cancellation"

# Commit message format:
# feat(scope): description     - New feature
# fix(scope): description      - Bug fix
# test(scope): description     - Test changes
# docs: description            - Documentation
# refactor(scope): description - Code refactoring
# chore: description           - Maintenance
```

### 5. Push and Create PR

```bash
git push origin feature/booking-cancellation
```

Then create PR on GitHub:
- **Base**: `develop`
- **Title**: Same as commit message
- **Description**: What and why

### 6. Review and Merge

GitHub Actions will automatically:
- ✅ Run all tests
- ✅ Check code formatting
- ✅ Run terraform plan
- ✅ Build Docker images

Once checks pass:
- Review your own PR (solo dev)
- Merge to `develop`

### 7. Auto-Deploy to Dev

After merge to `develop`:
- Terraform applies to dev environment
- Docker images built and pushed
- Flux deploys to dev EKS cluster
- E2E tests run on dev

**No manual steps needed!**

### 8. Test on Dev Environment

```bash
# Dev environment
https://crowdunlockedbeta.com

# Check deployment status
kubectl --context dev get pods
flux get kustomizations
```

### 9. Promote to Production

When ready to deploy to prod:

```bash
git checkout main
git pull origin main
git merge develop
git push origin main
```

Or create PR from `develop` to `main` on GitHub.

After merge to `main`:
- Terraform applies to prod environment
- Docker images built and pushed
- Flux deploys to prod EKS cluster
- Smoke tests run on prod

**No manual approval needed!**

## Quick Commands

### Testing

```bash
make test              # Run all tests
make test-unit         # Unit tests only
make test-integration  # Integration tests (needs Docker)
make test-web          # Web app tests
make test-mobile       # Mobile app tests
make test-watch        # Watch mode for rapid feedback
```

### Code Quality

```bash
make fmt               # Format all code
make lint              # Run all linters
```

### Infrastructure

```bash
make terraform-validate  # Validate terraform
make terraform-plan      # Plan changes (needs AWS creds)
```

### Build

```bash
make build             # Build all Go services
make docker-build      # Build Docker images
```

## Branching Strategy

```
feature/new-feature
        ↓
    develop ──────────→ Dev (crowdunlockedbeta.com)
        ↓
     main ────────────→ Prod (crowdunlocked.com)
```

## Environment URLs

- **Dev**: https://crowdunlockedbeta.com
- **Prod**: https://crowdunlocked.com

## Deployment Status

### Check Flux Status

```bash
# Dev cluster
flux --context dev get kustomizations

# Prod cluster
flux --context prod get kustomizations
```

### Check Pod Status

```bash
# Dev
kubectl --context dev get pods -n crowdunlocked

# Prod
kubectl --context prod get pods -n crowdunlocked
```

### View Logs

```bash
# Dev
kubectl --context dev logs -f deployment/bookings-service -n crowdunlocked

# Prod
kubectl --context prod logs -f deployment/bookings-service -n crowdunlocked
```

## Troubleshooting

### Tests Failing Locally

```bash
# Run specific test
cd services/bookings
go test -v -run TestCreateBooking ./internal/handler

# Check test output carefully
go test -v ./...

# Clean and retry
go clean -testcache
go test ./...
```

### CI/CD Failing

1. Check GitHub Actions tab for error details
2. Ensure all tests pass locally first
3. Check terraform plan output in PR
4. Verify AWS credentials are configured

### Deployment Issues

```bash
# Check Flux reconciliation
flux --context dev reconcile kustomization flux-system

# Check pod status
kubectl --context dev get pods -n crowdunlocked

# View pod logs
kubectl --context dev logs -f <pod-name> -n crowdunlocked

# Describe pod for events
kubectl --context dev describe pod <pod-name> -n crowdunlocked
```

### Rollback

```bash
# Revert the merge commit
git revert <commit-sha>
git push origin develop  # or main

# Flux will auto-deploy the reverted state
```

## Best Practices

### Do ✅
- Write tests first (TDD)
- Run tests locally before pushing
- Keep commits small and focused
- Use descriptive commit messages
- Format code before committing
- Review your own PRs carefully
- Test on dev before promoting to prod

### Don't ❌
- Skip writing tests
- Commit broken tests
- Push without running tests
- Merge failing PRs
- Deploy directly to prod
- Commit sensitive data
- Disable CI checks

## Getting Help

- **TDD Guide**: [docs/TDD_GUIDE.md](docs/TDD_GUIDE.md)
- **Architecture**: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- **Setup**: [docs/SETUP.md](docs/SETUP.md)
- **CI/CD**: [infra/terraform/CICD_SETUP.md](infra/terraform/CICD_SETUP.md)

## Solo Developer Notes

Since you're working solo:
- Self-review is fine
- No need to wait for approvals
- But still follow the process for discipline
- Tests are your safety net
- CI/CD gives you confidence

The workflow is designed to catch issues early and deploy safely, even without a team.
