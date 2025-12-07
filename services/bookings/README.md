# Bookings Service

Manages artist bookings and event scheduling for Crowd Unlocked.

## API Endpoints

### Create Booking
```
POST /api/v1/bookings
Content-Type: application/json

{
  "artist_id": "artist-123",
  "venue_id": "venue-456",
  "event_date": "2024-12-15T20:00:00Z",
  "fee": 5000.0
}
```

### Get Booking
```
GET /api/v1/bookings/{id}
```

### Confirm Booking
```
POST /api/v1/bookings/{id}/confirm
```

## Environment Variables

- `PORT`: Server port (default: 8080)
- `DYNAMODB_TABLE`: DynamoDB table name
- `AWS_REGION`: AWS region
- `AWS_XRAY_DAEMON_ADDRESS`: X-Ray daemon address

## Development

```bash
# Run tests
go test -v ./...

# Run locally
go run cmd/server/main.go

# Build
go build -o bookings ./cmd/server
```

## Testing

The service follows TDD principles. Tests are located in `*_test.go` files.

```bash
# Run tests with coverage
go test -v -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```
