import { createMocks } from 'node-mocks-http';
import handler from '../../../pages/api/v1/venues/index';
import type { NextApiRequest, NextApiResponse } from 'next';

// Mock AWS SDK.
// The handler module captures docClient (and its `send`) at import time, which
// happens before test setup runs. So the mock owns a single persistent send fn
// created inside the factory, and the test retrieves it via __getMockSend().
jest.mock('@aws-sdk/lib-dynamodb', () => {
  const send = jest.fn();
  return {
    DynamoDBDocumentClient: {
      from: jest.fn(() => ({ send })),
    },
    PutCommand: jest.fn(),
    ScanCommand: jest.fn(),
    __getMockSend: () => send,
  };
});

jest.mock('@aws-sdk/client-dynamodb', () => ({
  DynamoDBClient: jest.fn(() => ({})),
}));

const { __getMockSend } = require('@aws-sdk/lib-dynamodb');
const mockSend: jest.Mock = __getMockSend();

describe('/api/v1/venues', () => {
  beforeEach(() => {
    mockSend.mockReset();
  });

  describe('POST /api/v1/venues', () => {
    it('should create a new venue with valid data', async () => {
      mockSend.mockResolvedValue({});

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'POST',
        body: {
          name: 'Test Venue',
          location: { latitude: 40.7128, longitude: -74.0060 },
          address: {
            street: '123 Main St',
            city: 'New York',
            state: 'NY',
            postal_code: '10001',
            country: 'US'
          },
          venue_types: ['club']
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(201);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('id');
      expect(data.name).toBe('Test Venue');
      expect(data.location).toEqual({ latitude: 40.7128, longitude: -74.0060 });
      expect(data.venue_types).toEqual(['club']);
      expect(mockSend).toHaveBeenCalledTimes(1);
    });

    it('should return 400 for missing required fields', async () => {
      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'POST',
        body: {
          name: 'Test Venue',
          // missing location, address, venue_types
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(400);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('error');
    });

    it('should return 400 for invalid location coordinates', async () => {
      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'POST',
        body: {
          name: 'Test Venue',
          location: { latitude: 91, longitude: -74.0060 }, // invalid latitude
          address: {
            street: '123 Main St',
            city: 'New York',
            state: 'NY',
            postal_code: '10001',
            country: 'US'
          },
          venue_types: ['club']
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(400);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('error');
    });

    it('should return 400 for invalid venue types', async () => {
      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'POST',
        body: {
          name: 'Test Venue',
          location: { latitude: 40.7128, longitude: -74.0060 },
          address: {
            street: '123 Main St',
            city: 'New York',
            state: 'NY',
            postal_code: '10001',
            country: 'US'
          },
          venue_types: ['invalid_type']
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(400);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('error');
    });
  });

  describe('GET /api/v1/venues', () => {
    it('should return venues list', async () => {
      mockSend.mockResolvedValue({
        Items: [
          {
            id: 'venue-1',
            name: 'Test Venue 1',
            location: { latitude: 40.7128, longitude: -74.0060 },
            venue_types: ['club'],
            created_at: '2024-01-01T00:00:00.000Z',
            updated_at: '2024-01-01T00:00:00.000Z'
          }
        ],
        Count: 1
      });

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'GET',
        query: { limit: '10' },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(200);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('venues');
      expect(data.venues).toHaveLength(1);
      expect(data.venues[0].name).toBe('Test Venue 1');
    });

    it('should handle empty results', async () => {
      mockSend.mockResolvedValue({
        Items: [],
        Count: 0
      });

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'GET',
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(200);
      const data = JSON.parse(res._getData());
      expect(data.venues).toEqual([]);
      expect(data.total).toBe(0);
    });
  });

  it('should return 405 for unsupported methods', async () => {
    const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
      method: 'DELETE',
    });

    await handler(req, res);

    expect(res._getStatusCode()).toBe(405);
    expect(res._getHeaders().allow).toBe('GET, POST');
  });
});