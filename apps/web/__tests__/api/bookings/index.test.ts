import { createMocks } from 'node-mocks-http';
import handler from '../../../pages/api/v1/bookings/index';
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
    GetCommand: jest.fn(),
    __getMockSend: () => send,
  };
});

jest.mock('@aws-sdk/client-dynamodb', () => ({
  DynamoDBClient: jest.fn(() => ({})),
}));

const { __getMockSend } = require('@aws-sdk/lib-dynamodb');
const mockSend: jest.Mock = __getMockSend();

describe('/api/v1/bookings', () => {
  beforeEach(() => {
    mockSend.mockReset();
  });

  describe('POST /api/v1/bookings', () => {
    it('should create a new booking with valid data', async () => {
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
        method: 'POST',
        body: {
          venue_id: 'venue-123',
          user_id: 'user-456',
          event_name: 'Birthday Party',
          start_time: '2024-12-31T20:00:00.000Z',
          end_time: '2024-12-31T23:59:59.000Z'
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(201);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('id');
      expect(data.venue_id).toBe('venue-123');
      expect(data.user_id).toBe('user-456');
      expect(data.event_name).toBe('Birthday Party');
      expect(data.status).toBe('pending');
      expect(mockSend).toHaveBeenCalledTimes(2);
    });

    it('should return 400 for missing required fields', async () => {
      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'POST',
        body: {
          venue_id: 'venue-123',
          // missing user_id, event_name, start_time, end_time
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(400);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('error');
    });

    it('should return 400 for invalid date format', async () => {
      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'POST',
        body: {
          venue_id: 'venue-123',
          user_id: 'user-456',
          event_name: 'Birthday Party',
          start_time: 'invalid-date',
          end_time: '2024-12-31T23:59:59.000Z'
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(400);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('error');
    });

    it('should return 400 when end_time is before start_time', async () => {
      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'POST',
        body: {
          venue_id: 'venue-123',
          user_id: 'user-456',
          event_name: 'Birthday Party',
          start_time: '2024-12-31T23:59:59.000Z',
          end_time: '2024-12-31T20:00:00.000Z'
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(400);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('error');
    });

    it('should return 404 when venue does not exist', async () => {
      mockSend.mockResolvedValue({});

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'POST',
        body: {
          venue_id: 'non-existent',
          user_id: 'user-456',
          event_name: 'Birthday Party',
          start_time: '2024-12-31T20:00:00.000Z',
          end_time: '2024-12-31T23:59:59.000Z'
        },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(404);
      const data = JSON.parse(res._getData());
      expect(data).toHaveProperty('error');
    });
  });

  describe('GET /api/v1/bookings', () => {
    it('should return bookings list', async () => {
      mockSend.mockResolvedValue({
        Items: [
          {
            id: 'booking-1',
            venue_id: 'venue-123',
            user_id: 'user-456',
            event_name: 'Birthday Party',
            start_time: '2024-12-31T20:00:00.000Z',
            end_time: '2024-12-31T23:59:59.000Z',
            status: 'confirmed',
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
      expect(data).toHaveProperty('bookings');
      expect(data.bookings).toHaveLength(1);
      expect(data.bookings[0].event_name).toBe('Birthday Party');
    });

    it('should filter bookings by venue_id', async () => {
      mockSend.mockResolvedValue({
        Items: [
          {
            id: 'booking-1',
            venue_id: 'venue-123',
            user_id: 'user-456',
            event_name: 'Birthday Party',
            status: 'confirmed'
          }
        ],
        Count: 1
      });

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'GET',
        query: { venue_id: 'venue-123' },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(200);
      const data = JSON.parse(res._getData());
      expect(data.bookings).toHaveLength(1);
      expect(data.bookings[0].venue_id).toBe('venue-123');
    });

    it('should filter bookings by user_id', async () => {
      mockSend.mockResolvedValue({
        Items: [
          {
            id: 'booking-1',
            venue_id: 'venue-123',
            user_id: 'user-456',
            event_name: 'Birthday Party',
            status: 'confirmed'
          }
        ],
        Count: 1
      });

      const { req, res } = createMocks<NextApiRequest, NextApiResponse>({
        method: 'GET',
        query: { user_id: 'user-456' },
      });

      await handler(req, res);

      expect(res._getStatusCode()).toBe(200);
      const data = JSON.parse(res._getData());
      expect(data.bookings).toHaveLength(1);
      expect(data.bookings[0].user_id).toBe('user-456');
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
      expect(data.bookings).toEqual([]);
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