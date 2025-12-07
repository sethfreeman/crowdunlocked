# Test-Driven Development Guide

This project follows TDD principles. All new features and changes should be developed test-first.

## Testing Philosophy

- **Write tests first** - Define behavior before implementation
- **Red, Green, Refactor** - Fail, pass, improve
- **Fast feedback** - Tests should run quickly locally
- **Comprehensive coverage** - Unit, integration, and e2e tests

## Test Structure

### Go Services (Microservices)

Each service follows this structure:
```
services/[service-name]/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── booking.go
│   │   └── booking_test.go          # Unit tests
│   ├── handler/
│   │   ├── handler.go
│   │   └── handler_test.go          # Handler tests
│   └── repository/
│       ├── dynamodb.go
│       └── dynamodb_test.go         # Repository tests
└── tests/
    ├── integration/                  # Integration tests
    │   └── booking_flow_test.go
    └── e2e/                          # End-to-end tests
        └── api_test.go
```

### Web App (Next.js)

```
apps/web/
├── app/
│   └── [feature]/
│       ├── page.tsx
│       └── page.test.tsx            # Component tests
├── lib/
│   ├── utils.ts
│   └── utils.test.ts                # Unit tests
└── e2e/
    └── booking-flow.spec.ts         # Playwright e2e tests
```

### Mobile App (Flutter)

```
apps/mobile/
├── lib/
│   ├── features/
│   │   └── bookings/
│   │       ├── booking_screen.dart
│   │       └── booking_screen_test.dart
│   └── services/
│       ├── api_service.dart
│       └── api_service_test.dart
└── integration_test/
    └── app_test.dart
```

## Running Tests Locally

### Go Services

```bash
# Run all tests in a service
cd services/bookings
go test ./...

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestCreateBooking ./internal/handler

# Run integration tests (requires local DynamoDB)
docker-compose up -d dynamodb-local
go test -tags=integration ./tests/integration/...

# Watch mode (using entr or similar)
find . -name "*.go" | entr -c go test ./...
```

### Web App (Next.js)

```bash
cd apps/web

# Run all tests
npm test

# Watch mode
npm test -- --watch

# Coverage
npm test -- --coverage

# E2E tests
npm run test:e2e

# Type checking
npm run type-check
```

### Mobile App (Flutter)

```bash
cd apps/mobile

# Run all tests
flutter test

# Watch mode
flutter test --watch

# Coverage
flutter test --coverage

# Integration tests
flutter test integration_test/
```

### Infrastructure (OpenTofu)

```bash
cd infra/terraform/mgmt

# Validate
tofu validate

# Format check
tofu fmt -check

# Plan (dry run)
tofu plan -var-file=management.tfvars

# Use local tfvars for testing
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with test values
tofu plan -var-file=management.tfvars
```

## TDD Workflow

### 1. Feature Branch

```bash
# Create feature branch from develop
git checkout develop
git pull origin develop
git checkout -b feature/booking-cancellation
```

### 2. Write Failing Test

```go
// services/bookings/internal/handler/handler_test.go
func TestCancelBooking(t *testing.T) {
    // Arrange
    handler := NewHandler(mockRepo)
    req := httptest.NewRequest("DELETE", "/bookings/123", nil)
    
    // Act
    resp := handler.CancelBooking(req)
    
    // Assert
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

Run test: `go test ./internal/handler` → **Should FAIL** ❌

### 3. Implement Minimum Code

```go
// services/bookings/internal/handler/handler.go
func (h *Handler) CancelBooking(r *http.Request) *Response {
    // Minimal implementation to pass test
    return &Response{StatusCode: http.StatusOK}
}
```

Run test: `go test ./internal/handler` → **Should PASS** ✅

### 4. Refactor

Improve code quality while keeping tests green:
- Extract functions
- Improve naming
- Add error handling
- Optimize performance

Run tests after each change to ensure nothing breaks.

### 5. Commit and Push

```bash
git add .
git commit -m "feat(bookings): add booking cancellation"
git push origin feature/booking-cancellation
```

### 6. Create Pull Request

- PR title: `feat(bookings): add booking cancellation`
- Description: What and why
- Link to issue/ticket
- Tests pass locally ✅

## CI/CD Pipeline

### Pull Request (Any Branch → develop or main)

```
┌─────────────────────────────────────────────────────┐
│ 1. Lint & Format Check                              │
│    - Go: gofmt, golangci-lint                       │
│    - TypeScript: eslint, prettier                   │
│    - Terraform: terraform fmt                       │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│ 2. Unit Tests                                       │
│    - Go: go test ./...                              │
│    - Web: npm test                                  │
│    - Mobile: flutter test                           │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│ 3. Integration Tests                                │
│    - Services with test DynamoDB                    │
│    - API integration tests                          │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│ 4. Terraform Plan                                   │
│    - Plan for all environments                      │
│    - Post plan output to PR                         │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│ 5. Build & Security Scan                            │
│    - Docker image build                             │
│    - Trivy security scan                            │
└─────────────────────────────────────────────────────┘
```

**All checks must pass before merge** ✅

### Merge to develop → Deploy to Dev

```
PR Merged to develop
        ↓
┌─────────────────────────────────────────────────────┐
│ 1. Run All Tests Again                              │
└─────────────────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────────────────┐
│ 2. Terraform Apply (Management + Dev)               │
│    - Auto-approve (no manual step)                  │
└─────────────────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────────────────┐
│ 3. Build & Push Docker Images                       │
│    - Tag: dev-{git-sha}                             │
└─────────────────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────────────────┐
│ 4. Flux Auto-Deploy to Dev EKS                      │
│    - Watches develop branch                         │
│    - Applies k8s/overlays/dev                       │
│    - No manual approval needed                      │
└─────────────────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────────────────┐
│ 5. E2E Tests on Dev Environment                     │
│    - Smoke tests                                    │
│    - Critical user flows                            │
└─────────────────────────────────────────────────────┘
```

### Merge to main → Deploy to Prod

```
PR Merged to main (from develop)
        ↓
┌─────────────────────────────────────────────────────┐
│ 1. Run All Tests Again                              │
└─────────────────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────────────────┐
│ 2. Terraform Apply (Management + Prod)              │
│    - Auto-approve (no manual step)                  │
└─────────────────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────────────────┐
│ 3. Build & Push Docker Images                       │
│    - Tag: prod-{git-sha}                            │
└─────────────────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────────────────┐
│ 4. Flux Auto-Deploy to Prod EKS                     │
│    - Watches main branch                            │
│    - Applies k8s/overlays/prod                      │
│    - No manual approval needed                      │
└─────────────────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────────────────┐
│ 5. E2E Tests on Prod Environment                    │
│    - Smoke tests only                               │
│    - Alert on failure                               │
└─────────────────────────────────────────────────────┘
```

## Branching Strategy

```
feature/booking-cancellation
        ↓ (PR + Review)
    develop ──────────────→ Dev Environment (crowdunlockedbeta.com)
        ↓ (PR + Review)
     main ────────────────→ Prod Environment (crowdunlocked.com)
```

**Rules:**
- Feature branches created from `develop`
- PRs require passing tests (enforced by GitHub)
- Self-review allowed (solo developer)
- No manual approval in deployment pipeline
- Flux auto-deploys on branch changes

## Test Coverage Goals

- **Unit tests**: 80%+ coverage
- **Integration tests**: Critical paths covered
- **E2E tests**: Happy paths + critical errors

## Best Practices

### Do ✅
- Write tests before code
- Test behavior, not implementation
- Use descriptive test names
- Keep tests fast and isolated
- Mock external dependencies
- Run tests locally before pushing
- Fix broken tests immediately

### Don't ❌
- Skip tests for "simple" code
- Test implementation details
- Write slow tests
- Share state between tests
- Commit broken tests
- Disable tests to make CI pass
- Write tests after the fact

## Tools & Libraries

### Go
- Testing: `testing` (stdlib)
- Assertions: `github.com/stretchr/testify`
- Mocking: `github.com/stretchr/testify/mock`
- HTTP testing: `httptest` (stdlib)

### TypeScript/Next.js
- Testing: Jest
- React testing: React Testing Library
- E2E: Playwright
- Mocking: MSW (Mock Service Worker)

### Flutter
- Testing: `flutter_test`
- Mocking: `mockito`
- Integration: `integration_test`

## Example Test Files

See these examples for reference:
- Go service: `services/bookings/internal/handler/handler_test.go` (to be created)
- Next.js: `apps/web/app/bookings/page.test.tsx` (to be created)
- Flutter: `apps/mobile/lib/features/bookings/booking_screen_test.dart` (to be created)

## Getting Help

- Check existing tests for patterns
- Review test output carefully
- Use `-v` flag for verbose output
- Run single test to isolate issues
- Ask in PR comments for review
