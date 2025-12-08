# Venue Search Architecture

## Overview
Geospatial venue search system for CrowdUnlocked bookings service, enabling artists to discover performance venues on a map with advanced filtering.

## Key Features
- **Map-based search** with radius filtering
- **Multi-criteria filtering**: venue type, capacity, genres, payment, amenities
- **Geospatial indexing** using geohash for efficient queries
- **Multi-source data**: Songkick, Bandsintown, Google Places, user submissions
- **Real-time availability** checking

## Domain Models

### Venue
Core entity representing a performance venue with:
- Geographic location (lat/lng + geohash)
- Venue characteristics (type, capacity, genres)
- Payment information (range, type)
- Amenities and features
- Contact information
- External API IDs for syncing
- Verification status

### Search Criteria
Flexible filtering system supporting:
- Geographic: location + radius, city, state
- Venue: types, capacity range, genres, amenities
- Payment: min/max pay, payment types
- Quality: rating, verified status
- Sorting: distance, rating, capacity, pay

## Geospatial Strategy

### Geohash Implementation
- **Precision levels**: 4-7 characters based on search radius
- **Neighbor calculation**: Query 9 geohash cells to cover radius
- **Distance calculation**: Haversine formula for accurate results

### DynamoDB Schema
```
Primary Table: venues-{env}
├── PK: VENUE#{id}
├── SK: METADATA
├── GSI1 (Geohash): GEO#{geohash} → VENUE#{id}
├── GSI2 (City): CITY#{city}#{state} → VENUE#{name}
├── GSI3 (Type): TYPE#{type} → RATING#{rating}#VENUE#{id}
└── GSI4 (External): EXTERNAL#{source}#{id} → VENUE#{id}
```

### Query Pattern
1. Calculate geohash prefixes for search radius
2. Query GSI1 for all matching geohash prefixes
3. Filter results by additional criteria
4. Calculate exact distances
5. Sort and paginate

## API Integration Flow

```
┌─────────────────┐
│  Google Places  │──┐
└─────────────────┘  │
                     │
┌─────────────────┐  │    ┌──────────────┐
│    Songkick     │──┼───▶│ Aggregation  │
└─────────────────┘  │    │   Pipeline   │
                     │    └──────┬───────┘
┌─────────────────┐  │           │
│  Bandsintown    │──┘           │
└─────────────────┘              │
                                 ▼
                         ┌───────────────┐
                         │   DynamoDB    │
                         │ venues table  │
                         └───────┬───────┘
                                 │
                                 ▼
                         ┌───────────────┐
                         │  Search API   │
                         └───────┬───────┘
                                 │
                                 ▼
                         ┌───────────────┐
                         │  Map Frontend │
                         │  (Mapbox GL)  │
                         └───────────────┘
```

## Implementation Status

### ✅ Completed
- Venue domain model with comprehensive fields
- Search criteria and result types
- Geohash encoding/decoding utilities
- Distance calculation (Haversine)
- Geohash prefix generation for radius search
- Unit tests for geospatial functions

### 🚧 Next Steps
1. **DynamoDB Setup**
   - Create table with GSIs
   - Implement venue repository
   - Add geohash indexing

2. **Google Places Integration**
   - API client setup
   - Venue import pipeline
   - Photo fetching
   - Geocoding integration

3. **Search Implementation**
   - Geospatial query logic
   - Filter combinations
   - Pagination and sorting
   - REST API endpoints

4. **External API Integration**
   - Songkick sync pipeline
   - Bandsintown sync pipeline
   - Deduplication logic
   - Scheduled sync jobs

5. **Frontend**
   - Mapbox GL JS integration
   - Venue markers with clustering
   - Filter UI
   - Venue detail cards

## Performance Considerations

### Query Optimization
- Use geohash prefixes to limit scan scope
- Cache popular searches (city-level)
- Implement pagination to limit result sets
- Use DynamoDB parallel scans for large areas

### Data Freshness
- Real-time: User submissions, booking status
- Daily: Songkick/Bandsintown sync
- Weekly: Google Places refresh
- On-demand: Venue claims and updates

### Scalability
- DynamoDB on-demand pricing scales automatically
- Geohash indexing enables efficient spatial queries
- API rate limiting prevents external API overuse
- CloudFront caching for static venue data

## Cost Optimization

### API Usage
- Cache Google Places results (24 hours)
- Batch Songkick/Bandsintown requests
- Use free tiers where available
- Implement request deduplication

### Storage
- Store only essential venue data
- Compress photos (use URLs, not binary)
- Archive inactive venues
- Use DynamoDB TTL for temporary data

## Future Enhancements

### Phase 2
- **Availability calendar**: Real-time booking slots
- **Venue reviews**: Artist ratings and feedback
- **Smart recommendations**: ML-based venue matching
- **Route planning**: Multi-city tour optimization

### Phase 3
- **Direct booking**: In-app booking confirmation
- **Payment integration**: Deposits and guarantees
- **Contract generation**: Automated booking agreements
- **Analytics**: Venue performance metrics

## Testing Strategy

### Unit Tests
- Geohash encoding/decoding accuracy
- Distance calculation precision
- Search criteria validation
- Domain model business logic

### Integration Tests
- DynamoDB query patterns
- API client error handling
- Data sync pipelines
- Deduplication logic

### E2E Tests
- Map search workflows
- Filter combinations
- Booking request flow
- Venue claim process

## Monitoring & Alerts

### Metrics
- Search query latency (p50, p95, p99)
- API error rates by source
- Venue data completeness
- User engagement (searches, bookings)

### Alerts
- API rate limit approaching
- Sync job failures
- High query latency
- Data quality issues (duplicates, missing fields)
