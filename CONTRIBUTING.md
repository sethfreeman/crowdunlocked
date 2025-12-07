# Contributing to Crowd Unlocked

Thank you for your interest in contributing to Crowd Unlocked! 🎉

We're building an open-source artist management platform, and we welcome contributions from developers of all skill levels.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Workflow](#development-workflow)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Getting Help](#getting-help)

## Code of Conduct

This project follows a simple code of conduct:
- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Help others learn and grow

## How Can I Contribute?

### Reporting Bugs

Found a bug? Please open an issue with:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, Go version, etc.)

### Suggesting Features

Have an idea? Open an issue with:
- Clear description of the feature
- Use case and benefits
- Possible implementation approach

### Contributing Code

1. **Find an issue** or create one
2. **Comment** that you'd like to work on it
3. **Fork** the repository
4. **Create a branch** and implement your changes
5. **Submit a PR** for review

## Development Workflow

### 1. Set Up Your Environment

```bash
# Fork and clone the repo
git clone https://github.com/YOUR_USERNAME/crowdunlocked.git
cd crowdunlocked

# Install dependencies
make test  # This will download Go modules

# Start local environment
docker-compose up -d
AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test \
  bash scripts/create-dynamodb-tables.sh
```

### 2. Create a Feature Branch

```bash
# Always branch from develop
git checkout develop
git pull origin develop
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test additions/changes

### 3. Follow TDD (Test-Driven Development)

We practice TDD - write tests first!

```bash
# 1. Write a failing test
cd services/bookings
# Edit internal/domain/booking_test.go

# 2. Run tests (should fail)
go test ./...

# 3. Implement the feature
# Edit internal/domain/booking.go

# 4. Run tests (should pass)
go test ./...

# 5. Refactor if needed
```

See [docs/TDD_GUIDE.md](docs/TDD_GUIDE.md) for detailed examples.

### 4. Make Your Changes

- Write clean, readable code
- Follow existing code style
- Add comments for complex logic
- Update documentation if needed

### 5. Test Your Changes

```bash
# Run all tests
make test

# Run specific service tests
cd services/bookings
go test -v -cover ./...

# Test locally with Docker
docker-compose up --build
```

### 6. Commit Your Changes

Use conventional commit messages:

```bash
git add .
git commit -m "feat(bookings): add booking cancellation feature"
```

Commit message format:
- `feat(scope): description` - New feature
- `fix(scope): description` - Bug fix
- `docs: description` - Documentation
- `test(scope): description` - Tests
- `refactor(scope): description` - Code refactoring
- `chore: description` - Maintenance

### 7. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## Pull Request Process

### Before Submitting

- ✅ All tests pass locally
- ✅ Code is formatted (`make fmt`)
- ✅ No linting errors (`make lint`)
- ✅ Documentation is updated
- ✅ Commit messages follow conventions

### PR Description

Include:
- **What**: What does this PR do?
- **Why**: Why is this change needed?
- **How**: How does it work?
- **Testing**: How did you test it?
- **Screenshots**: If UI changes

### Review Process

1. Automated checks will run (tests, linting, etc.)
2. A maintainer will review your code
3. Address any feedback
4. Once approved, your PR will be merged!

### After Merge

- Your changes will auto-deploy to dev environment
- After testing, they'll be promoted to production
- You'll be added to our contributors list! 🎉

## Coding Standards

### Go Services

```go
// Good: Clear function names, error handling
func (s *BookingService) CreateBooking(ctx context.Context, req *CreateBookingRequest) (*Booking, error) {
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    
    booking := &Booking{
        ID:        uuid.New().String(),
        ArtistID:  req.ArtistID,
        VenueID:   req.VenueID,
        Status:    StatusPending,
        CreatedAt: time.Now(),
    }
    
    if err := s.repo.Create(ctx, booking); err != nil {
        return nil, fmt.Errorf("failed to create booking: %w", err)
    }
    
    return booking, nil
}
```

### TypeScript/Next.js

```typescript
// Good: Type safety, clear naming
interface BookingFormProps {
  onSubmit: (booking: CreateBookingRequest) => Promise<void>;
  initialData?: Booking;
}

export function BookingForm({ onSubmit, initialData }: BookingFormProps) {
  // Component implementation
}
```

### General Guidelines

- **DRY**: Don't Repeat Yourself
- **KISS**: Keep It Simple, Stupid
- **YAGNI**: You Aren't Gonna Need It
- **Single Responsibility**: One function, one purpose
- **Error Handling**: Always handle errors gracefully

## Testing Guidelines

### Test Coverage Goals

- **Unit tests**: 80%+ coverage
- **Integration tests**: Critical paths
- **E2E tests**: Happy paths + key errors

### Writing Good Tests

```go
func TestBooking_Confirm(t *testing.T) {
    // Arrange
    booking := &Booking{
        ID:     "test-123",
        Status: StatusPending,
    }
    
    // Act
    err := booking.Confirm()
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, StatusConfirmed, booking.Status)
    assert.NotNil(t, booking.ConfirmedAt)
}
```

### Test Naming

- `TestFunctionName` for unit tests
- `TestFunctionName_Scenario` for specific cases
- `TestFunctionName_Error` for error cases

## Project Structure

```
crowdunlocked/
├── services/           # Go microservices
│   ├── bookings/
│   │   ├── cmd/       # Main applications
│   │   ├── internal/  # Private code
│   │   └── tests/     # Integration tests
│   └── ...
├── apps/              # Frontend applications
│   ├── web/          # Next.js
│   └── mobile/       # Flutter
├── infra/            # Infrastructure as code
├── k8s/              # Kubernetes manifests
└── docs/             # Documentation
```

## Getting Help

- 📖 Read the [documentation](./docs)
- 💬 Ask in [Discussions](https://github.com/YOUR_USERNAME/crowdunlocked/discussions)
- 🐛 Check existing [Issues](https://github.com/YOUR_USERNAME/crowdunlocked/issues)
- 📧 Reach out to maintainers

## Recognition

Contributors will be:
- Added to our contributors list
- Mentioned in release notes
- Credited in documentation

## Questions?

Don't hesitate to ask! We're here to help. Open an issue or discussion, and we'll get back to you.

---

**Thank you for contributing to Crowd Unlocked!** 🎵✨
