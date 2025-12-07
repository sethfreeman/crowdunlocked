package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReleaseType string

const (
	TypeSingle ReleaseType = "single"
	TypeEP     ReleaseType = "ep"
	TypeAlbum  ReleaseType = "album"
)

type Release struct {
	ID          string      `dynamodbav:"id" json:"id"`
	ArtistID    string      `dynamodbav:"artist_id" json:"artist_id"`
	Title       string      `dynamodbav:"title" json:"title"`
	Type        ReleaseType `dynamodbav:"type" json:"type"`
	ReleaseDate time.Time   `dynamodbav:"release_date" json:"release_date"`
	TrackCount  int         `dynamodbav:"track_count" json:"track_count"`
	CoverArtURL string      `dynamodbav:"cover_art_url" json:"cover_art_url"`
	CreatedAt   time.Time   `dynamodbav:"created_at" json:"created_at"`
	UpdatedAt   time.Time   `dynamodbav:"updated_at" json:"updated_at"`
}

func NewRelease(artistID, title string, releaseType ReleaseType, releaseDate time.Time, trackCount int) *Release {
	now := time.Now()
	return &Release{
		ID:          uuid.New().String(),
		ArtistID:    artistID,
		Title:       title,
		Type:        releaseType,
		ReleaseDate: releaseDate,
		TrackCount:  trackCount,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
