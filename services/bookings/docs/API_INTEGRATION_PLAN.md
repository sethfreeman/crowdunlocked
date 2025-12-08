# Venue Data API Integration Plan

## Overview
Integrate multiple third-party APIs to build a comprehensive venue database for artist bookings, with geospatial search capabilities.

## API Integrations

### 1. Songkick API
**Purpose**: Venue data, event history, capacity information

**Endpoints**:
- `GET /api/3.0/venues.json` - Search venues
- `GET /api/3.0/venues/{id}.json` - Get venue details
- `GET /api/3.0/venues/{id}/calendar.json` - Venue events

**Data Mapping**:
```
Songkick → Our Model
- id → SongkickID
- displayName → Name
- lat/lng → Location
- capacity → Capacity
- metroArea → Address.City
```

**Rate Limits**: 5 requests/second
**Authentication**: API key
**Status**: Need to apply for API access

### 2. Bandsintown API
**Purpose**: Venue info, artist touring data, venue relationships

**Endpoints**:
- `GET /venues/{id}` - Venue details
- `GET /venues/search` - Search venues
- `GET /artists/{artist}/events` - Artist events at venues

**Data Mapping**:
```
Bandsintown → Our Model
- id → BandsintownID
- name → Name
- latitude/longitude → Location
- city/region/country → Address
```

**Rate Limits**: Unknown (need to verify)
**Authentication**: App ID
**Status**: Need to apply for API access

### 3. Google Places API
**Purpose**: Breweries, wineries, venue details, photos, reviews

**Endpoints**:
- `POST /places:searchNearby` - Find venues by type
- `GET /places/{id}` - Place details
- `GET /places/{id}/photos` - Venue photos

**Place Types**:
- `bar`
- `night_club`
- `restaurant`
- `cafe`
- Custom: brewery, winery (via text search)

**Data Mapping**:
```
Google Places → Our Model
- place_id → GooglePlaceID
- displayName → Name
- location → Location
- formattedAddress → Address
- rating → Rating
- userRatingCount → ReviewCount
- photos → Photos
```

**Rate Limits**: Pay per request
**Authentication**: API key
**Status**: Can start immediately

### 4. Mapbox Geocoding API
**Purpose**: Address geocoding, reverse geocoding, geohash generation

**Endpoints**:
- `GET /geocoding/v5/mapbox.places/{query}.json` - Forward geocoding
- `GET /geocoding/v5/mapbox.places/{lng},{lat}.json` - Reverse geocoding

**Usage**:
- Convert addresses to coordinates
- Generate geohashes for spatial indexing
- Validate venue locations

**Rate Limits**: 100,000 requests/month free tier
**Authentication**: Access token
**Status**: Can start immediately

## Data Aggregation Strategy

### Phase 1: Initial Data Collection
1. **Google Places** - Seed database with breweries, wineries, bars, clubs
2. **Mapbox** - Geocode and generate geohashes
3. Store in DynamoDB with `source=google_places`

### Phase 2: Venue Enrichment
1. **Songkick** - Match venues by name/location, add capacity and event history
2. **Bandsintown** - Cross-reference and fill gaps
3. Update records with external IDs

### Phase 3: User Contributions
1. Allow venues to claim profiles
2. Artists can submit new venues
3. Community verification system

### Phase 4: Continuous Sync
1. Nightly sync from Songkick/Bandsintown
2. Weekly refresh from Google Places
3. Real-time updates from user submissions

## DynamoDB Schema Design

### Primary Table: `venues-{env}`

**Primary Key**:
- `PK`: `VENUE#{id}`
- `SK`: `METADATA`

**GSI 1 - Geohash Index**:
- `GSI1PK`: `GEO#{geohash_prefix}`
- `GSI1SK`: `VENUE#{id}`
- Purpose: Spatial queries by geohash

**GSI 2 - City Index**:
- `GSI2PK`: `CITY#{city}#{state}`
- `GSI2SK`: `VENUE#{name}`
- Purpose: Search by city

**GSI 3 - Type Index**:
- `GSI3PK`: `TYPE#{venue_type}`
- `GSI3SK`: `RATING#{rating}#VENUE#{id}`
- Purpose: Filter by venue type, sort by rating

**GSI 4 - External ID Index**:
- `GSI4PK`: `EXTERNAL#{source}#{external_id}`
- `GSI4SK`: `VENUE#{id}`
- Purpose: Prevent duplicates, sync updates

### Geohash Strategy
- Use 6-character geohash for ~1.2km precision
- Store multiple precision levels (4, 5, 6 chars)
- Query by prefix for radius search

## Implementation Phases

### Phase 1: Foundation (Week 1-2)
- [ ] Create venue domain models ✅
- [ ] Set up DynamoDB table with GSIs
- [ ] Implement geohash utilities
- [ ] Create venue repository with CRUD operations

### Phase 2: Google Places Integration (Week 2-3)
- [ ] Set up Google Places API client
- [ ] Build venue import pipeline
- [ ] Implement photo fetching
- [ ] Add geocoding with Mapbox

### Phase 3: Search & Filtering (Week 3-4)
- [ ] Implement geospatial search
- [ ] Add filter combinations
- [ ] Build search API endpoints
- [ ] Add pagination and sorting

### Phase 4: Songkick/Bandsintown (Week 4-6)
- [ ] Apply for API access
- [ ] Build sync pipeline
- [ ] Implement deduplication logic
- [ ] Schedule periodic syncs

### Phase 5: Frontend Integration (Week 6-8)
- [ ] Map visualization with Mapbox GL JS
- [ ] Venue cards with filters
- [ ] Booking request flow
- [ ] Venue detail pages

## API Endpoints to Build

### Venue Search
```
GET /api/v1/venues/search
Query params:
  - lat, lng, radius_km
  - city, state, country
  - venue_types (comma-separated)
  - genres (comma-separated)
  - min_capacity, max_capacity
  - min_pay, max_pay
  - amenities (comma-separated)
  - min_rating
  - verified_only
  - limit, offset
  - sort_by, sort_order

Response: VenueSearchResult
```

### Venue Details
```
GET /api/v1/venues/{id}
Response: Venue
```

### Create Venue (User Submitted)
```
POST /api/v1/venues
Body: Venue (without ID)
Response: Venue
```

### Update Venue
```
PUT /api/v1/venues/{id}
Body: Partial<Venue>
Response: Venue
```

### Claim Venue
```
POST /api/v1/venues/{id}/claim
Body: { owner_id, verification_docs }
Response: { status: "pending" }
```

## Cost Estimates

### Google Places API
- $17 per 1,000 Place Details requests
- $7 per 1,000 Nearby Search requests
- $7 per 1,000 Photo requests
- **Estimated**: $200-500/month for initial seeding, $50-100/month ongoing

### Mapbox
- 100,000 requests/month free
- $0.50 per 1,000 requests after
- **Estimated**: Free tier sufficient initially

### Songkick/Bandsintown
- Typically free for non-commercial or partnership basis
- **Estimated**: $0 (pending approval)

### DynamoDB
- On-demand pricing
- **Estimated**: $10-50/month depending on scale

**Total Estimated**: $260-650 initial month, $60-150/month ongoing

## Security & Compliance

### API Key Management
- Store in AWS Secrets Manager
- Rotate keys quarterly
- Use separate keys per environment

### Rate Limiting
- Implement client-side rate limiting
- Queue requests to avoid hitting limits
- Cache responses where appropriate

### Data Privacy
- Don't store personal contact info without consent
- Allow venues to opt-out
- GDPR compliance for EU venues

## Success Metrics

### Data Quality
- Venue coverage: 10,000+ venues in first 6 months
- Data completeness: 80%+ venues with photos, contact info
- Accuracy: <5% duplicate venues

### Performance
- Search response time: <500ms p95
- Map load time: <2s
- Sync job completion: <1 hour

### User Engagement
- Venue searches per user: Track
- Booking requests sent: Track
- Venue claims: Track
