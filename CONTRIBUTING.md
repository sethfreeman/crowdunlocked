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

# Install web app dependencies
cd apps/web
npm install

# Run the dev server
npm run dev   # http://localhost:3000
```

The API routes read AWS config from environment variables (`apps/web/.env.local`).
To develop against the real dev tables, use an AWS profile with access to the
`venues-dev` / `bookings-dev` tables in `us-west-2`.

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
cd apps/web

# 1. Write a failing test in __tests__/ (e.g. a new API route test)
# 2. Run tests (should fail)
npm test

# 3. Implement the feature in pages/api/... or app/...
# 4. Run tests (should pass)
npm test

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
cd apps/web

# Run all checks the CI runs
npm test
npm run lint
npm run type-check
npm run build
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

- ✅ Tests pass locally (`npm test` in `apps/web`)
- ✅ Lint clean (`npm run lint`)
- ✅ Types clean (`npm run type-check`)
- ✅ Builds (`npm run build`)
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

- Vercel deploys your changes automatically
- You'll be added to our contributors list! 🎉

## Coding Standards

### TypeScript / Next.js API routes

```typescript
// Good: validate input, handle errors, return typed JSON
export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method !== 'POST') {
    res.setHeader('Allow', 'POST');
    return res.status(405).end(`Method ${req.method} Not Allowed`);
  }

  const { isValid, errors } = validateCreateVenueRequest(req.body);
  if (!isValid) {
    return res.status(400).json({ error: 'Validation failed', details: errors.join('; ') });
  }

  try {
    // ... DynamoDB operation via the AWS SDK ...
    return res.status(201).json(venue);
  } catch (err) {
    console.error('Error creating venue:', err);
    return res.status(500).json({ error: 'Internal server error' });
  }
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

API route tests use `node-mocks-http` and mock the DynamoDB client. See the
existing suites in `apps/web/__tests__/api/`.

```typescript
it('should return 400 for invalid venue types', async () => {
  const { req, res } = createMocks({
    method: 'POST',
    body: { name: 'Test', /* ... */, venue_types: ['invalid_type'] },
  });

  await handler(req, res);

  expect(res._getStatusCode()).toBe(400);
  expect(JSON.parse(res._getData())).toHaveProperty('error');
});
```

### Test Naming

- Describe the route/behavior under test, e.g. `should create a venue with valid data`
- Cover the happy path plus validation and error cases

## Project Structure

```
crowdunlocked/
├── apps/
│   ├── web/           # Next.js web app + API routes (deployed to Vercel)
│   │   ├── app/       # UI (App Router)
│   │   ├── pages/api/ # API routes
│   │   ├── lib/       # types, utils
│   │   └── __tests__/ # Jest tests
│   └── mobile/        # Flutter app
├── infra/terraform/   # OpenTofu (DynamoDB + OIDC)
└── docs/              # Documentation
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
