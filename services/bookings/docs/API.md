# Bookings Service API Documentation

## Base URL
- **Development**: `http://localhost:8080/api/v1`
- **Production**: `https://api.crowdunlocked.com/api/v1`

## Authentication
Currently no authentication required. Will be added in future versions.

## Common Headers
```
Content-Type: application/json
Accept: application/json
```

---

## Venues API

### Search Venues
Search for venues with flexible filtering and sorting.

**Endpoint**: `GET /venues/search`

**Query Parameters**:

| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `lat` | float | No* | Latitude for location search | `37.7749` |
| `lng` | float | No* | Longitude for location search | `-122.4194` |
| `radius` | float | No* | Search radius in kilometers | `5.0` |
| `city` | string | No* | City name | `San Francisco` |
| `state` | string | No | State/province code | `CA` |
| `country` | string | No | Country code | `US` |
| `venue_types` | string | No | Comma-separated venue types | `brewery,winery` |
| `min_capacity` | int | No | Minimum venue capacity | `50` |
| `max_capacity` | int | No | Maximum venue capacity | `500` |
| `genres` | string | No | Comma-separated genres | `rock,indie` |
| `min_pay` | int | No | Minimum payment amount | `100` |
| `max_pay` | int | No | Maximum payment amount | `1000` |
| `min_rating` | float | No | Minimum rating (0-5) | `4.0` |
| `verified_only` | boolean | No | Only verified venues | `true` |
| `active_only` | boolean | No | Only active venues | `true` |
| `limit` | int | No | Results per page (default: 10) | `20` |
| `offset` | int | No | Pagination offset (default: 0) | `0` |
| `sort_by` | string | No | Sort field | `distance`, `rating`, `capacity`, `pay`, `name`, `created_at` |
| `sort_order` | string | No | Sort direction | `asc`, `desc` |

*At least one of: `lat/lng/radius`, `city`, or `venue_types` is required.

**Venue Types**:
- `club` - Music club/venue
- `theater` - Theater
- `brewery` - Brewery
- `winery` - Winery
- `coffeehouse` - Coffee house
- `festival` - Festival
- `bar` - Bar
- `restaurant` - Restaurant
- `arena` - Arena
- `other` - Other

**Example Requests**:

```bash
# Search by location
GET /venues/search?lat=37.7749&lng=-122.4194&radius=5&limit=10

# Search by city with filters
GET /venues/search?city=San+Francisco&state=CA&venue_types=brewery,winery&min_capacity=50&max_capacity=200&min_rating=4.0&verified_only=true

# Search with sorting
GET /venues/search?city=Portland&state=OR&sort_by=rating&sort_order=desc&limit=20
```

**Response**: `200 OK`
```json
{
  "venues": [
    {
      "venue": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "The Brewery",
        "location": {
          "latitude": 37.7749,
          "longitude": -122.4194,
          "geohash": "9q8yyk"
        },
        "address": {
          "street": "123 Main St",
          "city": "San Francisco",
          "state": "CA",
          "postal_code": "94102",
          "country": "US"
        },
        "venue_types": ["brewery"],
        "capacity": 150,
        "genres": ["rock", "indie", "folk"],
        "pay_range": {
          "min": 200,
          "max": 500,
          "currency": "USD",
          "type": "guarantee",
          "notes": "Plus bar tab"
        },
        "amenities": ["sound_system", "parking", "green_room"],
        "photos": [
          "https://example.com/photo1.jpg"
        ],
        "contact_info": {
          "email": "booking@thebrewery.com",
          "phone": "+1-555-0123",
          "website": "https://thebrewery.com",
          "booking_url": "https://thebrewery.com/booking"
        },
        "rating": 4.5,
        "review_count": 42,
        "description": "Intimate brewery with great acoustics",
        "verified": true,
        "active": true,
        "source": "google_places",
        "created_at": "2025-01-15T10:30:00Z",
        "updated_at": "2025-01-15T10:30:00Z"
      },
      "distance_km": 2.3
    }
  ],
  "total": 25,
  "limit": 10,
  "offset": 0,
  "has_more": true
}
```

---

### Get Venue by ID
Retrieve detailed information about a specific venue.

**Endpoint**: `GET /venues/{id}`

**Path Parameters**:
- `id` (string, required): Venue ID

**Example Request**:
```bash
GET /venues/550e8400-e29b-41d4-a716-446655440000
```

**Response**: `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "The Brewery",
  "location": {
    "latitude": 37.7749,
    "longitude": -122.4194,
    "geohash": "9q8yyk"
  },
  "address": {
    "street": "123 Main St",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94102",
    "country": "US"
  },
  "venue_types": ["brewery"],
  "capacity": 150,
  "genres": ["rock", "indie"],
  "rating": 4.5,
  "verified": true,
  "active": true
}
```

**Error Response**: `404 Not Found`
```json
{
  "error": "venue not found"
}
```

---

### Create Venue
Create a new venue (user-submitted).

**Endpoint**: `POST /venues`

**Request Body**:
```json
{
  "name": "New Brewery",
  "location": {
    "latitude": 37.7749,
    "longitude": -122.4194
  },
  "address": {
    "street": "123 Main St",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94102",
    "country": "US"
  },
  "venue_types": ["brewery"],
  "capacity": 150,
  "genres": ["rock", "indie"],
  "description": "Great local brewery with live music"
}
```

**Required Fields**:
- `name`
- `location` (latitude and longitude)
- `venue_types` (at least one)

**Response**: `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "New Brewery",
  "location": {
    "latitude": 37.7749,
    "longitude": -122.4194,
    "geohash": "9q8yyk"
  },
  "verified": false,
  "active": true,
  "source": "user_submitted",
  "created_at": "2025-01-15T10:30:00Z"
}
```

**Error Response**: `400 Bad Request`
```json
{
  "error": "name is required"
}
```

---

### Update Venue
Update an existing venue (partial update supported).

**Endpoint**: `PUT /venues/{id}`

**Path Parameters**:
- `id` (string, required): Venue ID

**Request Body** (all fields optional):
```json
{
  "name": "Updated Name",
  "capacity": 200,
  "genres": ["rock", "metal"],
  "description": "Updated description"
}
```

**Response**: `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Updated Name",
  "capacity": 200,
  "updated_at": "2025-01-15T11:00:00Z"
}
```

**Error Response**: `404 Not Found`
```json
{
  "error": "venue not found"
}
```

---

### Delete Venue
Delete a venue.

**Endpoint**: `DELETE /venues/{id}`

**Path Parameters**:
- `id` (string, required): Venue ID

**Response**: `204 No Content`

**Error Response**: `500 Internal Server Error`
```json
{
  "error": "failed to delete venue"
}
```

---

## Bookings API

### Create Booking
Create a new booking request.

**Endpoint**: `POST /bookings`

**Request Body**:
```json
{
  "artist_id": "artist-123",
  "venue_id": "550e8400-e29b-41d4-a716-446655440000",
  "event_date": "2025-03-15T20:00:00Z",
  "fee": 500.00
}
```

**Response**: `201 Created`
```json
{
  "id": "booking-456",
  "artist_id": "artist-123",
  "venue_id": "550e8400-e29b-41d4-a716-446655440000",
  "event_date": "2025-03-15T20:00:00Z",
  "status": "pending",
  "fee": 500.00,
  "created_at": "2025-01-15T10:30:00Z"
}
```

---

### Get Booking by ID
Retrieve booking details.

**Endpoint**: `GET /bookings/{id}`

**Response**: `200 OK`
```json
{
  "id": "booking-456",
  "artist_id": "artist-123",
  "venue_id": "550e8400-e29b-41d4-a716-446655440000",
  "event_date": "2025-03-15T20:00:00Z",
  "status": "confirmed",
  "fee": 500.00
}
```

---

### Confirm Booking
Confirm a pending booking.

**Endpoint**: `POST /bookings/{id}/confirm`

**Response**: `200 OK`
```json
{
  "id": "booking-456",
  "status": "confirmed",
  "updated_at": "2025-01-15T11:00:00Z"
}
```

---

## Health Check

### Health Check
Check service health status.

**Endpoint**: `GET /health`

**Response**: `200 OK`
```
OK
```

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "error message description"
}
```

**HTTP Status Codes**:
- `200 OK` - Request succeeded
- `201 Created` - Resource created successfully
- `204 No Content` - Request succeeded with no response body
- `400 Bad Request` - Invalid request parameters
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## Rate Limiting
Currently no rate limiting. Will be added in future versions.

## Pagination
Use `limit` and `offset` parameters for pagination:
- `limit`: Number of results per page (default: 10, max: 100)
- `offset`: Number of results to skip (default: 0)

Example:
```bash
# Page 1
GET /venues/search?city=Portland&state=OR&limit=20&offset=0

# Page 2
GET /venues/search?city=Portland&state=OR&limit=20&offset=20

# Page 3
GET /venues/search?city=Portland&state=OR&limit=20&offset=40
```

---

## Examples

### Find breweries near me
```bash
curl "http://localhost:8080/api/v1/venues/search?lat=37.7749&lng=-122.4194&radius=10&venue_types=brewery&limit=10"
```

### Find verified venues in a city
```bash
curl "http://localhost:8080/api/v1/venues/search?city=Portland&state=OR&verified_only=true&min_rating=4.0&sort_by=rating&sort_order=desc"
```

### Create a new venue
```bash
curl -X POST http://localhost:8080/api/v1/venues \
  -H "Content-Type: application/json" \
  -d '{
    "name": "The Local Brewery",
    "location": {"latitude": 45.5231, "longitude": -122.6765},
    "address": {"city": "Portland", "state": "OR", "country": "US"},
    "venue_types": ["brewery"],
    "capacity": 100
  }'
```

### Update venue capacity
```bash
curl -X PUT http://localhost:8080/api/v1/venues/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"capacity": 200}'
```
