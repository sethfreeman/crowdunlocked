package repository

import (
	"context"

	"github.com/crowdunlocked/services/bookings/internal/domain"
)

// MockVenueRepository is an in-memory implementation for testing
type MockVenueRepository struct {
	venues map[string]*domain.Venue
}

// NewMockVenueRepository creates a new mock repository
func NewMockVenueRepository() *MockVenueRepository {
	return &MockVenueRepository{
		venues: make(map[string]*domain.Venue),
	}
}

func (r *MockVenueRepository) Create(ctx context.Context, venue *domain.Venue) error {
	r.venues[venue.ID] = venue
	return nil
}

func (r *MockVenueRepository) GetByID(ctx context.Context, id string) (*domain.Venue, error) {
	venue, ok := r.venues[id]
	if !ok {
		return nil, &VenueNotFoundError{}
	}
	return venue, nil
}

func (r *MockVenueRepository) Update(ctx context.Context, venue *domain.Venue) error {
	if _, ok := r.venues[venue.ID]; !ok {
		return &VenueNotFoundError{}
	}
	r.venues[venue.ID] = venue
	return nil
}

func (r *MockVenueRepository) Delete(ctx context.Context, id string) error {
	if _, ok := r.venues[id]; !ok {
		return &VenueNotFoundError{}
	}
	delete(r.venues, id)
	return nil
}

func (r *MockVenueRepository) SearchByGeohash(ctx context.Context, geohashPrefixes []string, limit int) ([]*domain.Venue, error) {
	results := make([]*domain.Venue, 0)
	for _, venue := range r.venues {
		for _, prefix := range geohashPrefixes {
			if len(venue.Location.Geohash) >= len(prefix) &&
				venue.Location.Geohash[:len(prefix)] == prefix {
				results = append(results, venue)
				break
			}
		}
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

func (r *MockVenueRepository) SearchByCity(ctx context.Context, city, state string, limit int) ([]*domain.Venue, error) {
	results := make([]*domain.Venue, 0)
	for _, venue := range r.venues {
		if venue.Address.City == city && venue.Address.State == state {
			results = append(results, venue)
			if len(results) >= limit {
				break
			}
		}
	}
	return results, nil
}

func (r *MockVenueRepository) SearchByType(ctx context.Context, venueType domain.VenueType, limit int) ([]*domain.Venue, error) {
	results := make([]*domain.Venue, 0)
	for _, venue := range r.venues {
		for _, vt := range venue.VenueTypes {
			if vt == venueType {
				results = append(results, venue)
				break
			}
		}
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

func (r *MockVenueRepository) GetByExternalID(ctx context.Context, source domain.DataSource, externalID string) (*domain.Venue, error) {
	for _, venue := range r.venues {
		if venue.Source == source {
			switch source {
			case domain.SourceSongkick:
				if venue.SongkickID == externalID {
					return venue, nil
				}
			case domain.SourceBandsintown:
				if venue.BandsintownID == externalID {
					return venue, nil
				}
			case domain.SourceGooglePlaces:
				if venue.GooglePlaceID == externalID {
					return venue, nil
				}
			}
		}
	}
	return nil, &VenueNotFoundError{}
}
