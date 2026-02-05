import type { NextApiRequest, NextApiResponse } from 'next';
import { DynamoDBClient } from '@aws-sdk/client-dynamodb';
import { DynamoDBDocumentClient, QueryCommand } from '@aws-sdk/lib-dynamodb';
import { encodeGeohash } from '../../../../lib/utils/geohash';
import type { Venue, VenueSearchParams } from '../../../../lib/types/venue';

const client = new DynamoDBClient({ region: process.env.AWS_REGION || 'us-east-1' });
const docClient = DynamoDBDocumentClient.from(client);

const VENUES_TABLE = process.env.DYNAMODB_VENUES_TABLE || 'venues-dev';
const VALID_VENUE_TYPES = ['club', 'bar', 'restaurant', 'theater', 'arena', 'stadium', 'conference_center', 'other'];

interface VenueSearchResponse {
  venues: Venue[];
  total: number;
  limit: number;
  offset: number;
}

interface ErrorResponse {
  error: string;
  details?: string;
}

function validateSearchParams(query: any): { isValid: boolean; errors: string[] } {
  const errors: string[] = [];

  // Validate coordinates if provided
  if (query.latitude !== undefined) {
    const lat = parseFloat(query.latitude);
    if (isNaN(lat) || lat < -90 || lat > 90) {
      errors.push('Latitude must be a number between -90 and 90');
    }
  }

  if (query.longitude !== undefined) {
    const lng = parseFloat(query.longitude);
    if (isNaN(lng) || lng < -180 || lng > 180) {
      errors.push('Longitude must be a number between -180 and 180');
    }
  }

  // Validate radius if provided
  if (query.radius_km !== undefined) {
    const radius = parseFloat(query.radius_km);
    if (isNaN(radius) || radius <= 0) {
      errors.push('Radius must be a positive number');
    }
  }

  // Validate venue types if provided
  if (query.venue_types) {
    const types = query.venue_types.split(',').map((t: string) => t.trim());
    const invalidTypes = types.filter((type: string) => !VALID_VENUE_TYPES.includes(type));
    if (invalidTypes.length > 0) {
      errors.push(`Invalid venue types: ${invalidTypes.join(', ')}. Valid types are: ${VALID_VENUE_TYPES.join(', ')}`);
    }
  }

  // Validate limit if provided
  if (query.limit !== undefined) {
    const limit = parseInt(query.limit);
    if (isNaN(limit) || limit <= 0 || limit > 100) {
      errors.push('Limit must be a number between 1 and 100');
    }
  }

  return { isValid: errors.length === 0, errors };
}

async function searchVenuesByLocation(
  latitude: number,
  longitude: number,
  radiusKm: number,
  venueTypes?: string[],
  limit: number = 20
): Promise<Venue[]> {
  // Generate geohash for the search location
  const geohash = encodeGeohash(latitude, longitude, 5);
  
  const command = new QueryCommand({
    TableName: VENUES_TABLE,
    IndexName: 'GeohashIndex',
    KeyConditionExpression: 'geohash = :geohash',
    ExpressionAttributeValues: {
      ':geohash': geohash,
    },
    Limit: limit,
  });

  const result = await docClient.send(command);
  let venues = result.Items || [];

  // Filter by distance (simple approximation)
  venues = venues.filter(item => {
    const distance = calculateDistance(
      latitude, longitude,
      item.location.latitude, item.location.longitude
    );
    return distance <= radiusKm;
  });

  // Filter by venue types if specified
  if (venueTypes && venueTypes.length > 0) {
    venues = venues.filter(item => 
      item.venue_types.some((type: string) => venueTypes.includes(type))
    );
  }

  return venues.map(mapDynamoItemToVenue);
}

async function searchVenuesByCity(
  city: string,
  state: string,
  venueTypes?: string[],
  limit: number = 20
): Promise<Venue[]> {
  const cityState = `${city.toLowerCase()}#${state.toLowerCase()}`;
  
  const command = new QueryCommand({
    TableName: VENUES_TABLE,
    IndexName: 'CityIndex',
    KeyConditionExpression: 'city_state = :city_state',
    ExpressionAttributeValues: {
      ':city_state': cityState,
    },
    Limit: limit,
  });

  const result = await docClient.send(command);
  let venues = result.Items || [];

  // Filter by venue types if specified
  if (venueTypes && venueTypes.length > 0) {
    venues = venues.filter(item => 
      item.venue_types.some((type: string) => venueTypes.includes(type))
    );
  }

  return venues.map(mapDynamoItemToVenue);
}

async function searchVenuesByType(
  venueTypes: string[],
  limit: number = 20
): Promise<Venue[]> {
  const allVenues: any[] = [];

  // Query each venue type separately
  for (const venueType of venueTypes) {
    const command = new QueryCommand({
      TableName: VENUES_TABLE,
      IndexName: 'VenueTypeIndex',
      KeyConditionExpression: 'venue_type = :venue_type',
      ExpressionAttributeValues: {
        ':venue_type': venueType,
      },
      Limit: Math.ceil(limit / venueTypes.length),
    });

    const result = await docClient.send(command);
    allVenues.push(...(result.Items || []));
  }

  // Remove duplicates and limit results
  const uniqueVenues = allVenues.filter((venue, index, self) => 
    index === self.findIndex(v => v.id === venue.id)
  ).slice(0, limit);

  return uniqueVenues.map(mapDynamoItemToVenue);
}

function mapDynamoItemToVenue(item: any): Venue {
  return {
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
  };
}

function calculateDistance(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const R = 6371; // Earth's radius in kilometers
  const dLat = (lat2 - lat1) * Math.PI / 180;
  const dLon = (lon2 - lon1) * Math.PI / 180;
  const a = 
    Math.sin(dLat/2) * Math.sin(dLat/2) +
    Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) * 
    Math.sin(dLon/2) * Math.sin(dLon/2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));
  return R * c;
}

async function searchVenues(req: NextApiRequest, res: NextApiResponse<VenueSearchResponse | ErrorResponse>) {
  try {
    const { isValid, errors } = validateSearchParams(req.query);
    
    if (!isValid) {
      return res.status(400).json({ error: 'Validation failed', details: errors.join('; ') });
    }

    const params: VenueSearchParams = {
      venue_types: req.query.venue_types ? req.query.venue_types.toString().split(',').map(t => t.trim()) : undefined,
      latitude: req.query.latitude ? parseFloat(req.query.latitude.toString()) : undefined,
      longitude: req.query.longitude ? parseFloat(req.query.longitude.toString()) : undefined,
      radius_km: req.query.radius_km ? parseFloat(req.query.radius_km.toString()) : 10,
      city: req.query.city?.toString(),
      state: req.query.state?.toString(),
      limit: req.query.limit ? parseInt(req.query.limit.toString()) : 20,
      offset: req.query.offset ? parseInt(req.query.offset.toString()) : 0,
    };

    let venues: Venue[] = [];

    // Determine search strategy based on provided parameters
    if (params.latitude !== undefined && params.longitude !== undefined) {
      // Location-based search
      venues = await searchVenuesByLocation(
        params.latitude,
        params.longitude,
        params.radius_km || 10,
        params.venue_types,
        params.limit
      );
    } else if (params.city && params.state) {
      // City-based search
      venues = await searchVenuesByCity(
        params.city,
        params.state,
        params.venue_types,
        params.limit
      );
    } else if (params.venue_types && params.venue_types.length > 0) {
      // Venue type search
      venues = await searchVenuesByType(params.venue_types, params.limit);
    } else {
      return res.status(400).json({ 
        error: 'Invalid search parameters', 
        details: 'Must provide either coordinates, city/state, or venue types' 
      });
    }

    const response: VenueSearchResponse = {
      venues,
      total: venues.length,
      limit: params.limit || 20,
      offset: params.offset || 0,
    };

    res.status(200).json(response);
  } catch (error) {
    console.error('Error searching venues:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<VenueSearchResponse | ErrorResponse>
) {
  switch (req.method) {
    case 'GET':
      return searchVenues(req, res);
    default:
      res.setHeader('Allow', 'GET');
      return res.status(405).end(`Method ${req.method} Not Allowed`);
  }
}