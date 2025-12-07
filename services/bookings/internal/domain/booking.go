package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID          string        `dynamodbav:"id" json:"id"`
	ArtistID    string        `dynamodbav:"artist_id" json:"artist_id"`
	VenueID     string        `dynamodbav:"venue_id" json:"venue_id"`
	EventDate   time.Time     `dynamodbav:"event_date" json:"event_date"`
	Status      BookingStatus `dynamodbav:"status" json:"status"`
	Fee         float64       `dynamodbav:"fee" json:"fee"`
	CreatedAt   time.Time     `dynamodbav:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `dynamodbav:"updated_at" json:"updated_at"`
}

func NewBooking(artistID, venueID string, eventDate time.Time, fee float64) *Booking {
	now := time.Now()
	return &Booking{
		ID:        uuid.New().String(),
		ArtistID:  artistID,
		VenueID:   venueID,
		EventDate: eventDate,
		Status:    StatusPending,
		Fee:       fee,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (b *Booking) Confirm() {
	b.Status = StatusConfirmed
	b.UpdatedAt = time.Now()
}

func (b *Booking) Cancel() {
	b.Status = StatusCancelled
	b.UpdatedAt = time.Now()
}
