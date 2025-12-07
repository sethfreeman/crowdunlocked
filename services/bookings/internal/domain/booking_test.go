package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBooking(t *testing.T) {
	artistID := "artist-123"
	venueID := "venue-456"
	eventDate := time.Now().Add(24 * time.Hour)
	fee := 5000.0

	booking := NewBooking(artistID, venueID, eventDate, fee)

	assert.NotEmpty(t, booking.ID)
	assert.Equal(t, artistID, booking.ArtistID)
	assert.Equal(t, venueID, booking.VenueID)
	assert.Equal(t, eventDate, booking.EventDate)
	assert.Equal(t, StatusPending, booking.Status)
	assert.Equal(t, fee, booking.Fee)
	assert.False(t, booking.CreatedAt.IsZero())
	assert.False(t, booking.UpdatedAt.IsZero())
}

func TestBooking_Confirm(t *testing.T) {
	booking := NewBooking("artist-1", "venue-1", time.Now(), 1000.0)
	originalUpdatedAt := booking.UpdatedAt

	time.Sleep(time.Millisecond)
	booking.Confirm()

	assert.Equal(t, StatusConfirmed, booking.Status)
	assert.True(t, booking.UpdatedAt.After(originalUpdatedAt))
}

func TestBooking_Cancel(t *testing.T) {
	booking := NewBooking("artist-1", "venue-1", time.Now(), 1000.0)
	originalUpdatedAt := booking.UpdatedAt

	time.Sleep(time.Millisecond)
	booking.Cancel()

	assert.Equal(t, StatusCancelled, booking.Status)
	assert.True(t, booking.UpdatedAt.After(originalUpdatedAt))
}
