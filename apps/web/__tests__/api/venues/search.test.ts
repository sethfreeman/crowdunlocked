import { createMocks } from 'node-mocks-http';
import handler from '../../../pages/api/v1/venues/search';
import type { NextApiRequest, NextApiResponse } from 'next';

// Mock AWS SDK
jest.mock('@aws-sdk/lib-dynamodb', () => ({
  DynamoDBDocumentClient: {
    from: jest.fn(() => ({
      send: jest.fn(),
    })),
  },
  QueryCommand: jest.fn(),
}));

jest.mock('@aws-sdk/client-dynamodb', () => ({
  DynamoDBClient: jest.fn(() => ({})),
}));

describe('/api/v1/venues/search', () => {
  let mockSend: jest.Mock;

  beforeEach(() => {
    jest.clearAllMocks();
    const { DynamoDBDocumentClient } = require('@aws-sdk/lib-dynamodb');
    mockSend = jest.fn();
    DynamoDBDocumentClient.from.mockReturnValue({ send: mockSend });
  });

  describe('GET /api/v1/venues/search', () => {
    it('should search venues by location', async () => {
      mockSend.mockResolvedValue({
        Items: [
          {
            id: 'venue-1',
            name: 'Test Club',
            location: { latitude: 40.7128, longitude: -74.0060 },
            venue_types: ['club'],
            geohash: 'dr5ru',
            created_at: '2024-01-01T00:00:00.000Z',
            updated_at: '2024-01-01T00:00:00.000Z'
          }
        ],
        Count: 1
      });

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'GET',
        query: {
          latitude: '40.7128',
          longitude: '-74.0060',
          radius_km: '5',
          venue_types: 'club'
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(200);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('venues');
      expect(data.venues).toHaveLength(1);
      expect(data.venues[0].name).toBe('Test Club');
      expect(mockSend).toHaveBeenCalledTimes(1);
    });

    it('should search venues by city and state', async () => {
      mockSend.mockResolvedValue({
        Items: [
          {
            id: 'venue-2',
            name: 'NYC Bar',
            location: { latitude: 40.7589, longitude: -73.9851 },
            venue_types: ['bar'],
            city_state: 'new york#ny',
            created_at: '2024-01-01T00:00:00.000Z',
            updated_at: '2024-01-01T00:00:00.000Z'
          }
        ],
        Count: 1
      });

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'GET',
        query: {
          city: 'New York',
          state: 'NY',
          venue_types: 'bar'
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(200);
      const data = JSON.parse(res._getData());
      expect(data.venues).toHaveLength(1);
      expect(data.venues[0].name).toBe('NYC Bar');
    });

    it('should search venues by venue type only', async () => {
      mockSend.mockResolvedValue({
        Items: [
          {
            id: 'venue-3',
            name: 'Theater Venue',
            location: { latitude: 40.7505, longitude: -73.9934 },
            venue_types: ['theater'],
            venue_type: 'theater',
            created_at: '2024-01-01T00:00:00.000Z',
            updated_at: '2024-01-01T00:00:00.000Z'
          }
        ],
        Count: 1
      });

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'GET',
        query: {
          venue_types: 'theater'
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(200);
      const data = JSON.parse(res._getData());
      expect(data.venues).toHaveLength(1);
      expect(data.venues[0].name).toBe('Theater Venue');
    });

    it('should handle multiple venue types', async () => {
      mockSend
        .mockResolvedValueOnce({
          Items: [{ id: 'venue-1', name: 'Club 1', venue_types: ['club'] }],
          Count: 1
        })
        .mockResolvedValueOnce({
          Items: [{ id: 'venue-2', name: 'Bar 1', venue_types: ['bar'] }],
          Count: 1
        });

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'GET',
        query: {
          venue_types: 'club,bar'
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(200);
      const data = JSON.parse(res._getData());
      expect(data.venues).toHaveLength(2);
      expect(mockSend).toHaveBeenCalledTimes(2);
    });

    it('should return 400 for invalid coordinates', async () => {
      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'GET',
        query: {
          latitude: '91', // invalid
          longitude: '-74.0060',
          radius_km: '5'
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(400);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('error');
    });

    it('should return 400 for invalid venue types', async () => {
      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'GET',
        query: {
          venue_types: 'invalid_type'
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(400);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('error');
    });

    it('should handle empty results', async () => {
      mockSend.mockResolvedValue({
        Items: [],
        Count: 0
      });

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'GET',
        query: {
          venue_types: 'club'
        },
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
      method: 'POST',
    });

    await handler(req, res);

    expect(res._getStatusCode()).toBe(405);
    expect(res._getHeaders().allow).toBe('GET');
  });
});