# Test-Driven Development Guide

This project follows TDD. New features and changes should be developed test-first.

## Philosophy

- **Write tests first** — define behavior before implementation
- **Red, Green, Refactor** — fail, pass, improve
- **Fast feedback** — tests run in well under a second locally
- **Cover behavior** — happy path plus validation and error cases

## Where tests live

The app is a Next.js project in `apps/web`. Tests use **Jest** with
**node-mocks-http** for API route testing.

```
apps/web/
├── pages/api/
│   ├── health.ts
│   └── v1/
│       ├── venues/        # index.ts, [id].ts, search.ts
│       └── bookings/      # index.ts
├── lib/                   # types, utils (e.g. geohash)
└── __tests__/
    └── api/
        ├── venues/        # index.test.ts, [id].test.ts, search.test.ts
        └── bookings/      # index.test.ts
```

## Running tests

```bash
cd apps/web

npm test                 # run all tests
npm test -- --watch      # watch mode
npm test -- --coverage   # coverage report
npm run type-check       # tsc --noEmit
npm run lint             # next lint
```

## Mocking DynamoDB in API route tests

The API route handlers create the DynamoDB document client at module load time:

```ts
const client = new DynamoDBClient({ region: process.env.AWS_REGION });
const docClient = DynamoDBDocumentClient.from(client);
```

Because `docClient` (and its `send`) is captured at import time, the mock must
own a single persistent `send` function that `from()` always returns. The test
retrieves that function via a helper exported from the mock factory:

```ts
jest.mock('@aws-sdk/lib-dynamodb', () => {
  const send = jest.fn();
  return {
    DynamoDBDocumentClient: { from: jest.fn(() => ({ send })) },
    QueryCommand: jest.fn(),
    // ...other commands the route uses...
    __getMockSend: () => send,
  };
});

jest.mock('@aws-sdk/client-dynamodb', () => ({
  DynamoDBClient: jest.fn(() => ({})),
}));

const { __getMockSend } = require('@aws-sdk/lib-dynamodb');
const mockSend: jest.Mock = __getMockSend();

beforeEach(() => mockSend.mockReset());
```

Then drive behavior per-test with `mockSend.mockResolvedValue({ Items: [...] })`.

## TDD workflow

### 1. Feature branch

```bash
git checkout develop
git pull origin develop
git checkout -b feature/booking-cancellation
```

### 2. Write a failing test

```ts
// apps/web/__tests__/api/bookings/cancel.test.ts
it('cancels a booking', async () => {
  mockSend.mockResolvedValueOnce({ Item: { id: 'b-1', status: 'pending' } }); // get
  mockSend.mockResolvedValueOnce({});                                          // update

  const { req, res } = createMocks({ method: 'DELETE', query: { id: 'b-1' } });
  await handler(req, res);

  expect(res._getStatusCode()).toBe(200);
});
```

`npm test` → **should FAIL** ❌

### 3. Implement the minimum to pass

Add the handler logic in the relevant `pages/api/...` file.

`npm test` → **should PASS** ✅

### 4. Refactor

Improve naming, extract helpers, tighten validation, keeping tests green.

### 5. Commit and open a PR

```bash
git add .
git commit -m "feat(bookings): add booking cancellation"
git push origin feature/booking-cancellation
```

Open a PR to `develop`.

## CI

`.github/workflows/ci.yaml` runs on every PR and push:

```
web job:
  npm ci → lint → test → type-check → build

terraform-plan job (PRs):
  tofu init → validate → plan (dev, plan-only)
```

All checks must pass before merge. CI does not deploy — **Vercel** deploys the
web app automatically via its GitHub integration.

## Branching

```
feature/*  ──PR──▶  develop  ──PR──▶  main
```

- Feature branches from `develop`
- PRs require passing CI (self-review OK for solo dev)
- Vercel deploys on merge

## Best practices

### Do ✅
- Write tests before code
- Test behavior, not implementation details
- Use descriptive test names
- Keep tests fast and isolated
- Reset mocks between tests
- Run tests locally before pushing

### Don't ❌
- Skip tests for "simple" code
- Share state between tests
- Commit broken tests
- Disable tests to make CI pass

## Tools

- **Test runner**: Jest
- **HTTP mocking**: node-mocks-http
- **AWS SDK**: mocked via `jest.mock` (see above)
- **Types**: TypeScript (`tsc --noEmit`)
- **Lint**: `next lint` (ESLint)

## Reference

See the existing suites in `apps/web/__tests__/api/` for working examples of the
mock setup and assertion patterns.
