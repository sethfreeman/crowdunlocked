import type { NextApiRequest, NextApiResponse } from 'next';
import { PutCommand, ScanCommand } from '@aws-sdk/lib-dynamodb';
import { ulid } from 'ulid';
import { docClient } from '../../../../lib/dynamodb';
import { encodeGeohash } from '../../../../lib/utils/geohash';
import type { Venue, CreateVenueRequest } from '../../../../lib/types/venue';

const VENUES_TABLE = process.env.DYNAMODB_VENUES_TABLE || 'venues-dev';
const VALID_VENUE_TYPES = ['club', 'bar', 'restaurant', 'theater', 'arena', 'stadium', 'conference_center', 'other'];

interface VenuesListResponse {
  venues: Venue[];
  total: number;
  limit: number;
  offset: number;
}

interface ErrorResponse {
  error: string;
  details?: string;
}

function validateCreateVenueRequest(body: any): { isValid: boolean; errors: string[] } {
  const errors: string[] = [];

  if (!body.name || typeof body.name !== 'string' || body.name.trim().length === 0) {
    errors.push('Name is required and must be a non-empty string');
  }

  if (!body.location || typeof body.location !== 'object') {
    errors.push('Location is required');
  } else {
    if (typeof body.location.latitude !== 'number' || body.location.latitude < -90 || body.location.latitude > 90) {
      errors.push('Location latitude must be a number between -90 and 90');
    }
    if (typeof body.location.longitude !== 'number' || body.location.longitude < -180 || body.location.longitude > 180) {
      errors.push('Location longitude must be a number between -180 and 180');
    }
  }

  if (!body.address || typeof body.address !== 'object') {
    errors.push('Address is required');
  } else {
    const requiredAddressFields = ['street', 'city', 'state', 'postal_code', 'country'];
    for (const field of requiredAddressFields) {
      if (!body.address[field] || typeof body.address[field] !== 'string' || body.address[field].trim().length === 0) {
        errors.push(`Address ${field} is required and must be a non-empty string`);
      }
    }
  }

  if (!Array.isArray(body.venue_types) || body.venue_types.length === 0) {
    errors.push('Venue types is required and must be a non-empty array');
  } else {
    const invalidTypes = body.venue_types.filter((type: any) => !VALID_VENUE_TYPES.includes(type));
    if (invalidTypes.length > 0) {
      errors.push(`Invalid venue types: ${invalidTypes.join(', ')}. Valid types are: ${VALID_VENUE_TYPES.join(', ')}`);
    }
  }

  if (body.capacity !== undefined && (typeof body.capacity !== 'number' || body.capacity < 0)) {
    errors.push('Capacity must be a non-negative number');
  }

  return { isValid: errors.length === 0, errors };
}

async function createVenue(req: NextApiRequest, res: NextApiResponse<Venue | ErrorResponse>) {
  try {
    const { isValid, errors } = validateCreateVenueRequest(req.body);
    
    if (!isValid) {
      return res.status(400).json({ error: 'Validation failed', details: errors.join('; ') });
    }

    const createRequest: CreateVenueRequest = req.body;
    const now = new Date().toISOString();
    const id = ulid();
    
    // Generate geohash for location-based queries
    const geohash = encodeGeohash(createRequest.location.latitude, createRequest.location.longitude, 5);
    const geohashSort = `${geohash}#${id}`;
    
    // Create city_state for city-based queries
    const cityState = `${createRequest.address.city.toLowerCase()}#${createRequest.address.state.toLowerCase()}`;
    
    const venue: Venue = {
      id,
      name: createRequest.name,
      location: createRequest.location,
      address: createRequest.address,
      venue_types: createRequest.venue_types,
      description: createRequest.description,
      capacity: createRequest.capacity,
      external_source: createRequest.external_source,
      external_source_id: createRequest.external_source_id,
      created_at: now,
      updated_at: now,
    };

    // Add DynamoDB-specific fields for GSI queries
    const dynamoItem = {
      ...venue,
      geohash,
      geohash_sort: geohashSort,
      city_state: cityState,
      // Add venue_type entries for each type (for VenueTypeIndex)
      venue_type: createRequest.venue_types[0], // Primary venue type
      rating_id: `0#${id}`, // Default rating of 0
    };

    const command = new PutCommand({
      TableName: VENUES_TABLE,
      Item: dynamoItem,
    });

    await docClient.send(command);

    res.status(201).json(venue);
  } catch (error) {
    console.error('Error creating venue:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
}

async function getVenues(req: NextApiRequest, res: NextApiResponse<VenuesListResponse | ErrorResponse>) {
  try {
    const limit = Math.min(parseInt(req.query.limit as string) || 20, 100);
    const offset = parseInt(req.query.offset as string) || 0;

    const command = new ScanCommand({
      TableName: VENUES_TABLE,
      Limit: limit,
      // Note: DynamoDB Scan doesn't support offset directly
      // In a production app, you'd use pagination tokens
    });

    const result = await docClient.send(command);
    
    const venues: Venue[] = (result.Items || []).map(item => ({
      id: item.id,
      name: item.name,
      location: item.location,
      address: item.address,
      venue_types: item.venue_types,
      description: item.description,
      capacity: item.capacity,
      rating: item.rating,
      external_source: item.external_source,
      external_source_id: item.external_source_id,
      created_at: item.created_at,
      updated_at: item.updated_at,
    }));

    const response: VenuesListResponse = {
      venues,
      total: result.Count || 0,
      limit,
      offset,
    };

    res.status(200).json(response);
  } catch (error) {
    console.error('Error fetching venues:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Venue | VenuesListResponse | ErrorResponse>
) {
  switch (req.method) {
    case 'POST':
      return createVenue(req, res);
    case 'GET':
      return getVenues(req, res);
    default:
      res.setHeader('Allow', 'GET, POST');
      return res.status(405).end(`Method ${req.method} Not Allowed`);
  }
}