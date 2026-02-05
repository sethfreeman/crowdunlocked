import type { NextApiRequest, NextApiResponse } from 'next';
import { DynamoDBClient } from '@aws-sdk/client-dynamodb';
import { DynamoDBDocumentClient, PutCommand, ScanCommand, GetCommand } from '@aws-sdk/lib-dynamodb';
import { ulid } from 'ulid';
import type { Booking, CreateBookingRequest } from '../../../../lib/types/venue';

const client = new DynamoDBClient({ region: process.env.AWS_REGION || 'us-east-1' });
const docClient = DynamoDBDocumentClient.from(client);

const BOOKINGS_TABLE = process.env.DYNAMODB_BOOKINGS_TABLE || 'bookings-dev';
const VENUES_TABLE = process.env.DYNAMODB_VENUES_TABLE || 'venues-dev';

interface BookingsListResponse {
  bookings: Booking[];
  total: number;
  limit: number;
  offset: number;
}

interface ErrorResponse {
  error: string;
  details?: string;
}

function validateCreateBookingRequest(body: any): { isValid: boolean; errors: string[] } {
  const errors: string[] = [];

  if (!body.venue_id || typeof body.venue_id !== 'string' || body.venue_id.trim().length === 0) {
    errors.push('Venue ID is required and must be a non-empty string');
  }

  if (!body.user_id || typeof body.user_id !== 'string' || body.user_id.trim().length === 0) {
    errors.push('User ID is required and must be a non-empty string');
  }

  if (!body.event_name || typeof body.event_name !== 'string' || body.event_name.trim().length === 0) {
    errors.push('Event name is required and must be a non-empty string');
  }

  if (!body.start_time || typeof body.start_time !== 'string') {
    errors.push('Start time is required and must be a string');
  } else {
    const startTime = new Date(body.start_time);
    if (isNaN(startTime.getTime())) {
      errors.push('Start time must be a valid ISO 8601 date string');
    }
  }

  if (!body.end_time || typeof body.end_time !== 'string') {
    errors.push('End time is required and must be a string');
  } else {
    const endTime = new Date(body.end_time);
    if (isNaN(endTime.getTime())) {
      errors.push('End time must be a valid ISO 8601 date string');
    }
  }

  // Validate that end_time is after start_time
  if (body.start_time && body.end_time) {
    const startTime = new Date(body.start_time);
    const endTime = new Date(body.end_time);
    if (!isNaN(startTime.getTime()) && !isNaN(endTime.getTime()) && endTime <= startTime) {
      errors.push('End time must be after start time');
    }
  }

  return { isValid: errors.length === 0, errors };
}

async function createBooking(req: NextApiRequest, res: NextApiResponse<Booking | ErrorResponse>) {
  try {
    const { isValid, errors } = validateCreateBookingRequest(req.body);
    
    if (!isValid) {
      return res.status(400).json({ error: 'Validation failed', details: errors.join('; ') });
    }

    const createRequest: CreateBookingRequest = req.body;

    // Verify that the venue exists
    const venueCommand = new GetCommand({
      TableName: VENUES_TABLE,
      Key: { id: createRequest.venue_id },
    });

    const venueResult = await docClient.send(venueCommand);

    if (!venueResult.Item) {
      return res.status(404).json({ error: 'Venue not found' });
    }

    const now = new Date().toISOString();
    const id = ulid();
    
    const booking: Booking = {
      id,
      venue_id: createRequest.venue_id,
      user_id: createRequest.user_id,
      event_name: createRequest.event_name,
      start_time: createRequest.start_time,
      end_time: createRequest.end_time,
      status: 'pending',
      created_at: now,
      updated_at: now,
    };

    const command = new PutCommand({
      TableName: BOOKINGS_TABLE,
      Item: booking,
    });

    await docClient.send(command);

    res.status(201).json(booking);
  } catch (error) {
    console.error('Error creating booking:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
}

async function getBookings(req: NextApiRequest, res: NextApiResponse<BookingsListResponse | ErrorResponse>) {
  try {
    const limit = Math.min(parseInt(req.query.limit as string) || 20, 100);
    const offset = parseInt(req.query.offset as string) || 0;
    const venueId = req.query.venue_id as string;
    const userId = req.query.user_id as string;

    let command;

    if (venueId || userId) {
      // Use scan with filter for now - in production, you'd want GSIs for these queries
      const filterExpressions: string[] = [];
      const expressionAttributeValues: Record<string, any> = {};

      if (venueId) {
        filterExpressions.push('venue_id = :venue_id');
        expressionAttributeValues[':venue_id'] = venueId;
      }

      if (userId) {
        filterExpressions.push('user_id = :user_id');
        expressionAttributeValues[':user_id'] = userId;
      }

      command = new ScanCommand({
        TableName: BOOKINGS_TABLE,
        FilterExpression: filterExpressions.join(' AND '),
        ExpressionAttributeValues: expressionAttributeValues,
        Limit: limit,
      });
    } else {
      command = new ScanCommand({
        TableName: BOOKINGS_TABLE,
        Limit: limit,
      });
    }

    const result = await docClient.send(command);
    
    const bookings: Booking[] = (result.Items || []).map(item => ({
      id: item.id,
      venue_id: item.venue_id,
      user_id: item.user_id,
      event_name: item.event_name,
      start_time: item.start_time,
      end_time: item.end_time,
      status: item.status,
      created_at: item.created_at,
      updated_at: item.updated_at,
    }));

    const response: BookingsListResponse = {
      bookings,
      total: result.Count || 0,
      limit,
      offset,
    };

    res.status(200).json(response);
  } catch (error) {
    console.error('Error fetching bookings:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Booking | BookingsListResponse | ErrorResponse>
) {
  switch (req.method) {
    case 'POST':
      return createBooking(req, res);
    case 'GET':
      return getBookings(req, res);
    default:
      res.setHeader('Allow', 'GET, POST');
      return res.status(405).end(`Method ${req.method} Not Allowed`);
  }
}