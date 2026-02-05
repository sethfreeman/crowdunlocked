// Optional: configure or set up a testing framework before each test.
// If you delete this file, remove `setupFilesAfterEnv` from `jest.config.js`

// Mock AWS SDK
jest.mock('aws-amplify', () => ({
  Amplify: {
    configure: jest.fn(),
  },
}));

// Mock environment variables
process.env.AWS_REGION = 'us-east-1';
process.env.DYNAMODB_VENUES_TABLE = 'venues-test';
process.env.DYNAMODB_BOOKINGS_TABLE = 'bookings-test';