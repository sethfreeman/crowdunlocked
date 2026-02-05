import { createMocks } from 'node-mocks-http';
import handler from '../../../pages/api/v1/venues/[id]';
import type { NextApiRequest, NextApiResponse } from 'next';

// Mock AWS SDK
jest.mock('@aws-sdk/lib-dynamodb', () => ({
  DynamoDBDocumentClient: {
    from: jest.fn(() => ({
      send: jest.fn(),
    })),
  },
  GetCommand: jest.fn(),
  UpdateCommand: jest.fn(),
  DeleteCommand: jest.fn(),
}));

jest.mock('@aws-sdk/client-dynamodb', () => ({
  DynamoDBClient: jest.fn(() => ({})),
}));

describe('/api/v1/venues/[id]', () => {
  let mockSend: jest.Mock;

  beforeEach(() => {
    jest.clearAllMocks();
    const { DynamoDBDocumentClient } = require('@aws-sdk/lib-dynamodb');
    mockSend = jest.fn();
    DynamoDBDocumentClient.from.mockReturnValue({ send: mockSend });
  });

  describe('GET /api/v1/venues/[id]', () => {
    it('should return venue by id', async () => {
      mockSend.mockResolvedValue({
        Item: {
          id: 'venue-123',
          name: 'Test Venue',
          location: { latitude: 40.7128, longitude: -74.0060 },
          address: {
            street: '123 Main St',
            city: 'New York',
            state: 'NY',
            postal_code: '10001',
            country: 'US'
          },
          venue_types: ['club'],
          created_at: '2024-01-01T00:00:00.000Z',
          updated_at: '2024-01-01T00:00:00.000Z'
        }
      });

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'GET',
        query: { id: 'venue-123' },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(200);
      const data = JSON.parse(res._getData());
      expect(data.id).toBe('venue-123');
      expect(data.name).toBe('Test Venue');
      expect(mockSend).toHaveBeenCalledTimes(1);
    });

    it('should return 404 for non-existent venue', async () => {
      mockSend.mockResolvedValue({});

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'GET',
        query: { id: 'non-existent' },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(404);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('error');
    });
  });

  describe('PUT /api/v1/venues/[id]', () => {
    it('should update venue with valid data', async () => {
      mockSend
        .mockResolvedValueOnce({
          Item: {
            id: 'venue-123',
            name: 'Old Name',
            location: { latitude: 40.7128, longitude: -74.0060 },
            address: {
              street: '123 Main St',
              city: 'New York',
              state: 'NY',
              postal_code: '10001',
              country: 'US'
            },
            venue_types: ['club'],
            created_at: '2024-01-01T00:00:00.000Z',
            updated_at: '2024-01-01T00:00:00.000Z'
          }
        })
        .mockResolvedValueOnce({
          Attributes: {
            id: 'venue-123',
            name: 'Updated Name',
            location: { latitude: 40.7128, longitude: -74.0060 },
            address: {
              street: '123 Main St',
              city: 'New York',
              state: 'NY',
              postal_code: '10001',
              country: 'US'
            },
            venue_types: ['club'],
            created_at: '2024-01-01T00:00:00.000Z',
            updated_at: '2024-01-01T01:00:00.000Z'
          }
        });

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'PUT',
        query: { id: 'venue-123' },
        body: {
          name: 'Updated Name'
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(200);
      const data = JSON.parse(res._getData());
      expect(data.name).toBe('Updated Name');
      expect(mockSend).toHaveBeenCalledTimes(2);
    });

    it('should return 404 when updating non-existent venue', async () => {
      mockSend.mockResolvedValue({});

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'PUT',
        query: { id: 'non-existent' },
        body: { name: 'Updated Name' },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(404);
    });

    it('should return 400 for invalid update data', async () => {
      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'PUT',
        query: { id: 'venue-123' },
        body: {
          location: { latitude: 91, longitude: -74.0060 } // invalid latitude
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(400);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('error');
    });
  });

  describe('DELETE /api/v1/venues/[id]', () => {
    it('should delete venue', async () => {
      mockSend
        .mockResolvedValueOnce({
          Item: {
            id: 'venue-123',
            name: 'Test Venue',
            venue_types: ['club']
          }
        })
        .mockResolvedValueOnce({});

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'DELETE',
        query: { id: 'venue-123' },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(204);
      expect(mockSend).toHaveBeenCalledTimes(2);
    });

    it('should return 404 when deleting non-existent venue', async () => {
      mockSend.mockResolvedValue({});

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'DELETE',
        query: { id: 'non-existent' },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(404);
    });
  });

  it('should return 405 for unsupported methods', async () => {
    const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
      method: 'PATCH',
      query: { id: 'venue-123' },
    });

    await handler(req, res);

    expect(res._getStatusCode()).toBe(405);
    expect(res._getHeaders().allow).toBe('GET, PUT, DELETE');
  });
});