import type { NextApiRequest, NextApiResponse } from 'next';
import { DynamoDBClient } from '@aws-sdk/client-dynamodb';
import { DynamoDBDocumentClient, GetCommand, UpdateCommand, DeleteCommand } from '@aws-sdk/lib-dynamodb';
import { encodeGeohash } from '../../../../lib/utils/geohash';
import type { Venue, UpdateVenueRequest } from '../../../../lib/types/venue';

const client = new DynamoDBClient({ region: process.env.AWS_REGION || 'us-east-1' });
const docClient = DynamoDBDocumentClient.from(client);

const VENUES_TABLE = process.env.DYNAMODB_VENUES_TABLE || 'venues-dev';
const VALID_VENUE_TYPES = ['club', 'bar', 'restaurant', 'theater', 'arena', 'stadium', 'conference_center', 'other'];

interface ErrorResponse {
  error: string;
  details?: string;
}

function validateUpdateVenueRequest(body: any): { isValid: boolean; errors: string[] } {
  const errors: string[] = [];

  if (body.name !== undefined && (typeof body.name !== 'string' || body.name.trim().length === 0)) {
    errors.push('Name must be a non-empty string');
  }

  if (body.location !== undefined) {
    if (typeof body.location !== 'object') {
      errors.push('Location must be an object');
    } else {
      if (body.location.latitude !== undefined && (typeof body.location.latitude !== 'number' || body.location.latitude < -90 || body.location.latitude > 90)) {
        errors.push('Location latitude must be a number between -90 and 90');
      }
      if (body.location.longitude !== undefined && (typeof body.location.longitude !== 'number' || body.location.longitude < -180 || body.location.longitude > 180)) {
        errors.push('Location longitude must be a number between -180 and 180');
      }
    }
  }

  if (body.address !== undefined) {
    if (typeof body.address !== 'object') {
      errors.push('Address must be an object');
    } else {
      const addressFields = ['street', 'city', 'state', 'postal_code', 'country'];
      for (const field of addressFields) {
        if (body.address[field] !== undefined && (typeof body.address[field] !== 'string' || body.address[field].trim().length === 0)) {
          errors.push(`Address ${field} must be a non-empty string`);
        }
      }
    }
  }

  if (body.venue_types !== undefined) {
    if (!Array.isArray(body.venue_types) || body.venue_types.length === 0) {
      errors.push('Venue types must be a non-empty array');
    } else {
      const invalidTypes = body.venue_types.filter((type: any) => !VALID_VENUE_TYPES.includes(type));
      if (invalidTypes.length > 0) {
        errors.push(`Invalid venue types: ${invalidTypes.join(', ')}. Valid types are: ${VALID_VENUE_TYPES.join(', ')}`);
      }
    }
  }

  if (body.capacity !== undefined && (typeof body.capacity !== 'number' || body.capacity < 0)) {
    errors.push('Capacity must be a non-negative number');
  }

  return { isValid: errors.length === 0, errors };
}

async function getVenue(req: NextApiRequest, res: NextApiResponse<Venue | ErrorResponse>) {
  try {
    const { id } = req.query;

    if (!id || typeof id !== 'string') {
      return res.status(400).json({ error: 'Invalid venue ID' });
    }

    const command = new GetCommand({
      TableName: VENUES_TABLE,
      Key: { id },
    });

    const result = await docClient.send(command);

    if (!result.Item) {
      return res.status(404).json({ error: 'Venue not found' });
    }

    const venue: Venue = {
      id: result.Item.id,
      name: result.Item.name,
      location: result.Item.location,
      address: result.Item.address,
      venue_types: result.Item.venue_types,
      description: result.Item.description,
      capacity: result.Item.capacity,
      rating: result.Item.rating,
      external_source: result.Item.external_source,
      external_source_id: result.Item.external_source_id,
      created_at: result.Item.created_at,
      updated_at: result.Item.updated_at,
    };

    res.status(200).json(venue);
  } catch (error) {
    console.error('Error fetching venue:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
}

async function updateVenue(req: NextApiRequest, res: NextApiResponse<Venue | ErrorResponse>) {
  try {
    const { id } = req.query;

    if (!id || typeof id !== 'string') {
      return res.status(400).json({ error: 'Invalid venue ID' });
    }

    const { isValid, errors } = validateUpdateVenueRequest(req.body);
    
    if (!isValid) {
      return res.status(400).json({ error: 'Validation failed', details: errors.join('; ') });
    }

    // First, check if venue exists
    const getCommand = new GetCommand({
      TableName: VENUES_TABLE,
      Key: { id },
    });

    const existingVenue = await docClient.send(getCommand);

    if (!existingVenue.Item) {
      return res.status(404).json({ error: 'Venue not found' });
    }

    const updateRequest: UpdateVenueRequest = req.body;
    const now = new Date().toISOString();

    // Build update expression
    const updateExpressions: string[] = [];
    const expressionAttributeNames: Record<string, string> = {};
    const expressionAttributeValues: Record<string, any> = {};

    if (updateRequest.name !== undefined) {
      updateExpressions.push('#name = :name');
      expressionAttributeNames['#name'] = 'name';
      expressionAttributeValues[':name'] = updateRequest.name;
    }

    if (updateRequest.location !== undefined) {
      updateExpressions.push('#location = :location');
      expressionAttributeNames['#location'] = 'location';
      expressionAttributeValues[':location'] = updateRequest.location;

      // Update geohash if location changed
      const geohash = encodeGeohash(updateRequest.location.latitude, updateRequest.location.longitude, 5);
      const geohashSort = `${geohash}#${id}`;
      
      updateExpressions.push('geohash = :geohash', 'geohash_sort = :geohash_sort');
      expressionAttributeValues[':geohash'] = geohash;
      expressionAttributeValues[':geohash_sort'] = geohashSort;
    }

    if (updateRequest.address !== undefined) {
      updateExpressions.push('#address = :address');
      expressionAttributeNames['#address'] = 'address';
      expressionAttributeValues[':address'] = updateRequest.address;

      // Update city_state if address changed
      const cityState = `${updateRequest.address.city.toLowerCase()}#${updateRequest.address.state.toLowerCase()}`;
      updateExpressions.push('city_state = :city_state');
      expressionAttributeValues[':city_state'] = cityState;
    }

    if (updateRequest.venue_types !== undefined) {
      updateExpressions.push('venue_types = :venue_types');
      expressionAttributeValues[':venue_types'] = updateRequest.venue_types;

      // Update primary venue_type for GSI
      updateExpressions.push('venue_type = :venue_type');
      expressionAttributeValues[':venue_type'] = updateRequest.venue_types[0];
    }

    if (updateRequest.description !== undefined) {
      updateExpressions.push('description = :description');
      expressionAttributeValues[':description'] = updateRequest.description;
    }

    if (updateRequest.capacity !== undefined) {
      updateExpressions.push('capacity = :capacity');
      expressionAttributeValues[':capacity'] = updateRequest.capacity;
    }

    // Always update the updated_at timestamp
    updateExpressions.push('updated_at = :updated_at');
    expressionAttributeValues[':updated_at'] = now;

    const updateCommand = new UpdateCommand({
      TableName: VENUES_TABLE,
      Key: { id },
      UpdateExpression: `SET ${updateExpressions.join(', ')}`,
      ExpressionAttributeNames: Object.keys(expressionAttributeNames).length > 0 ? expressionAttributeNames : undefined,
      ExpressionAttributeValues: expressionAttributeValues,
      ReturnValues: 'ALL_NEW',
    });

    const result = await docClient.send(updateCommand);

    if (!result.Attributes) {
      return res.status(500).json({ error: 'Failed to update venue' });
    }

    const venue: Venue = {
      id: result.Attributes.id,
      name: result.Attributes.name,
      location: result.Attributes.location,
      address: result.Attributes.address,
      venue_types: result.Attributes.venue_types,
      description: result.Attributes.description,
      capacity: result.Attributes.capacity,
      rating: result.Attributes.rating,
      external_source: result.Attributes.external_source,
      external_source_id: result.Attributes.external_source_id,
      created_at: result.Attributes.created_at,
      updated_at: result.Attributes.updated_at,
    };

    res.status(200).json(venue);
  } catch (error) {
    console.error('Error updating venue:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
}

async function deleteVenue(req: NextApiRequest, res: NextApiResponse<void | ErrorResponse>) {
  try {
    const { id } = req.query;

    if (!id || typeof id !== 'string') {
      return res.status(400).json({ error: 'Invalid venue ID' });
    }

    // First, check if venue exists
    const getCommand = new GetCommand({
      TableName: VENUES_TABLE,
      Key: { id },
    });

    const existingVenue = await docClient.send(getCommand);

    if (!existingVenue.Item) {
      return res.status(404).json({ error: 'Venue not found' });
    }

    const deleteCommand = new DeleteCommand({
      TableName: VENUES_TABLE,
      Key: { id },
    });

    await docClient.send(deleteCommand);

    res.status(204).end();
  } catch (error) {
    console.error('Error deleting venue:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Venue | void | ErrorResponse>
) {
  switch (req.method) {
    case 'GET':
      return getVenue(req, res);
    case 'PUT':
      return updateVenue(req, res);
    case 'DELETE':
      return deleteVenue(req, res);
    default:
      res.setHeader('Allow', 'GET, PUT, DELETE');
      return res.status(405).end(`Method ${req.method} Not Allowed`);
  }
}