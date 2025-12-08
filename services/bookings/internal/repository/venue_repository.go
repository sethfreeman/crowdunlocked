package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/crowdunlocked/services/bookings/internal/domain"
)

// VenueRepository defines the interface for venue data access
type VenueRepository interface {
	Create(ctx context.Context, venue *domain.Venue) error
	GetByID(ctx context.Context, id string) (*domain.Venue, error)
	Update(ctx context.Context, venue *domain.Venue) error
	Delete(ctx context.Context, id string) error
	SearchByGeohash(ctx context.Context, geohashPrefixes []string, limit int) ([]*domain.Venue, error)
	SearchByCity(ctx context.Context, city, state string, limit int) ([]*domain.Venue, error)
	SearchByType(ctx context.Context, venueType domain.VenueType, limit int) ([]*domain.Venue, error)
	GetByExternalID(ctx context.Context, source domain.DataSource, externalID string) (*domain.Venue, error)
}

// DynamoDBVenueRepository implements VenueRepository using DynamoDB
type DynamoDBVenueRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDBVenueRepository creates a new DynamoDB venue repository
func NewDynamoDBVenueRepository(client *dynamodb.Client, tableName string) *DynamoDBVenueRepository {
	return &DynamoDBVenueRepository{
		client:    client,
		tableName: tableName,
	}
}

// venueItem represents the DynamoDB item structure with GSI attributes
type venueItem struct {
	*domain.Venue
	// GSI attributes
	Geohash           string `dynamodbav:"geohash"`
	GeohashSort       string `dynamodbav:"geohash_sort"`
	CityState         string `dynamodbav:"city_state"`
	VenueType         string `dynamodbav:"venue_type"`
	RatingID          string `dynamodbav:"rating_id"`
	ExternalSourceID  string `dynamodbav:"external_source_id,omitempty"`
}

// toVenueItem converts a domain.Venue to a venueItem with GSI attributes
func toVenueItem(v *domain.Venue) *venueItem {
	item := &venueItem{
		Venue:       v,
		Geohash:     v.Location.Geohash,
		GeohashSort: fmt.Sprintf("%s#%s", v.Location.Geohash, v.ID),
		CityState:   fmt.Sprintf("%s#%s", v.Address.City, v.Address.State),
		RatingID:    fmt.Sprintf("%010.2f#%s", v.Rating, v.ID),
	}

	// Set primary venue type for GSI3
	if len(v.VenueTypes) > 0 {
		item.VenueType = string(v.VenueTypes[0])
	}

	// Set external source ID for GSI4
	switch v.Source {
	case domain.SourceSongkick:
		if v.SongkickID != "" {
			item.ExternalSourceID = fmt.Sprintf("%s#%s", v.Source, v.SongkickID)
		}
	case domain.SourceBandsintown:
		if v.BandsintownID != "" {
			item.ExternalSourceID = fmt.Sprintf("%s#%s", v.Source, v.BandsintownID)
		}
	case domain.SourceGooglePlaces:
		if v.GooglePlaceID != "" {
			item.ExternalSourceID = fmt.Sprintf("%s#%s", v.Source, v.GooglePlaceID)
		}
	}

	return item
}

// Create creates a new venue in DynamoDB
func (r *DynamoDBVenueRepository) Create(ctx context.Context, venue *domain.Venue) error {
	item := toVenueItem(venue)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal venue: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to create venue: %w", err)
	}

	return nil
}

// GetByID retrieves a venue by ID
func (r *DynamoDBVenueRepository) GetByID(ctx context.Context, id string) (*domain.Venue, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get venue: %w", err)
	}

	if result.Item == nil {
		return nil, &VenueNotFoundError{}
	}

	var item venueItem
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal venue: %w", err)
	}

	return item.Venue, nil
}

// Update updates an existing venue
func (r *DynamoDBVenueRepository) Update(ctx context.Context, venue *domain.Venue) error {
	item := toVenueItem(venue)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal venue: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to update venue: %w", err)
	}

	return nil
}

// Delete deletes a venue by ID
func (r *DynamoDBVenueRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete venue: %w", err)
	}

	return nil
}

// SearchByGeohash searches venues by geohash prefixes
func (r *DynamoDBVenueRepository) SearchByGeohash(ctx context.Context, geohashPrefixes []string, limit int) ([]*domain.Venue, error) {
	venues := make([]*domain.Venue, 0)

	for _, prefix := range geohashPrefixes {
		result, err := r.client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(r.tableName),
			IndexName:              aws.String("GeohashIndex"),
			KeyConditionExpression: aws.String("geohash = :geohash"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":geohash": &types.AttributeValueMemberS{Value: prefix},
			},
			Limit: aws.Int32(int32(limit)),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query by geohash: %w", err)
		}

		for _, item := range result.Items {
			var venueItem venueItem
			err = attributevalue.UnmarshalMap(item, &venueItem)
			if err != nil {
				continue
			}
			venues = append(venues, venueItem.Venue)
		}

		if len(venues) >= limit {
			break
		}
	}

	return venues, nil
}

// SearchByCity searches venues by city and state
func (r *DynamoDBVenueRepository) SearchByCity(ctx context.Context, city, state string, limit int) ([]*domain.Venue, error) {
	cityState := fmt.Sprintf("%s#%s", city, state)

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("CityIndex"),
		KeyConditionExpression: aws.String("city_state = :city_state"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":city_state": &types.AttributeValueMemberS{Value: cityState},
		},
		Limit: aws.Int32(int32(limit)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query by city: %w", err)
	}

	venues := make([]*domain.Venue, 0, len(result.Items))
	for _, item := range result.Items {
		var venueItem venueItem
		err = attributevalue.UnmarshalMap(item, &venueItem)
		if err != nil {
			continue
		}
		venues = append(venues, venueItem.Venue)
	}

	return venues, nil
}

// SearchByType searches venues by venue type
func (r *DynamoDBVenueRepository) SearchByType(ctx context.Context, venueType domain.VenueType, limit int) ([]*domain.Venue, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("VenueTypeIndex"),
		KeyConditionExpression: aws.String("venue_type = :venue_type"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":venue_type": &types.AttributeValueMemberS{Value: string(venueType)},
		},
		Limit:            aws.Int32(int32(limit)),
		ScanIndexForward: aws.Bool(false), // Sort by rating descending
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query by type: %w", err)
	}

	venues := make([]*domain.Venue, 0, len(result.Items))
	for _, item := range result.Items {
		var venueItem venueItem
		err = attributevalue.UnmarshalMap(item, &venueItem)
		if err != nil {
			continue
		}
		venues = append(venues, venueItem.Venue)
	}

	return venues, nil
}

// VenueNotFoundError is returned when a venue is not found
type VenueNotFoundError struct{}

func (e *VenueNotFoundError) Error() string {
	return "venue not found"
}

// GetByExternalID retrieves a venue by external source ID
func (r *DynamoDBVenueRepository) GetByExternalID(ctx context.Context, source domain.DataSource, externalID string) (*domain.Venue, error) {
	externalSourceID := fmt.Sprintf("%s#%s", source, externalID)

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("ExternalIdIndex"),
		KeyConditionExpression: aws.String("external_source_id = :external_source_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":external_source_id": &types.AttributeValueMemberS{Value: externalSourceID},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query by external ID: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, &VenueNotFoundError{}
	}

	var venueItem venueItem
	err = attributevalue.UnmarshalMap(result.Items[0], &venueItem)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal venue: %w", err)
	}

	return venueItem.Venue, nil
}
