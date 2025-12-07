package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/crowdunlocked/services/bookings/internal/domain"
)

type BookingRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewBookingRepository(client *dynamodb.Client, tableName string) *BookingRepository {
	return &BookingRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *BookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	item, err := attributevalue.MarshalMap(booking)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

func (r *BookingRepository) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var booking domain.Booking
	err = attributevalue.UnmarshalMap(result.Item, &booking)
	return &booking, err
}

func (r *BookingRepository) Update(ctx context.Context, booking *domain.Booking) error {
	item, err := attributevalue.MarshalMap(booking)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}
