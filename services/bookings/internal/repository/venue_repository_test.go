package repository

import (
	"context"
	"testing"

	"github.com/crowdunlocked/services/bookings/internal/domain"
)

func TestVenueRepository_Create(t *testing.T) {
	repo := NewMockVenueRepository()
	ctx := context.Background()

	venue := domain.NewVenue(
		"Test Brewery",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeBrewery},
		domain.SourceGooglePlaces,
	)

	err := repo.Create(ctx, venue)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify venue was created
	retrieved, err := repo.GetByID(ctx, venue.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if retrieved.ID != venue.ID {
		t.Errorf("GetByID() ID = %v, want %v", retrieved.ID, venue.ID)
	}
	if retrieved.Name != venue.Name {
		t.Errorf("GetByID() Name = %v, want %v", retrieved.Name, venue.Name)
	}
}

func TestVenueRepository_GetByID(t *testing.T) {
	repo := NewMockVenueRepository()
	ctx := context.Background()

	venue := domain.NewVenue(
		"Test Venue",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	_ = repo.Create(ctx, venue)

	retrieved, err := repo.GetByID(ctx, venue.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if retrieved.ID != venue.ID {
		t.Errorf("GetByID() ID = %v, want %v", retrieved.ID, venue.ID)
	}
}

func TestVenueRepository_GetByID_NotFound(t *testing.T) {
	repo := NewMockVenueRepository()
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "nonexistent-id")
	if err == nil {
		t.Error("GetByID() should return error for nonexistent ID")
	}
}

func TestVenueRepository_Update(t *testing.T) {
	repo := NewMockVenueRepository()
	ctx := context.Background()

	venue := domain.NewVenue(
		"Original Name",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeBrewery},
		domain.SourceGooglePlaces,
	)
	_ = repo.Create(ctx, venue)

	// Update venue
	venue.Name = "Updated Name"
	venue.Capacity = 200
	venue.Update()

	err := repo.Update(ctx, venue)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	// Verify update
	retrieved, _ := repo.GetByID(ctx, venue.ID)
	if retrieved.Name != "Updated Name" {
		t.Errorf("Update() Name = %v, want Updated Name", retrieved.Name)
	}
	if retrieved.Capacity != 200 {
		t.Errorf("Update() Capacity = %v, want 200", retrieved.Capacity)
	}
}

func TestVenueRepository_Delete(t *testing.T) {
	repo := NewMockVenueRepository()
	ctx := context.Background()

	venue := domain.NewVenue(
		"Test Venue",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	_ = repo.Create(ctx, venue)

	err := repo.Delete(ctx, venue.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(ctx, venue.ID)
	if err == nil {
		t.Error("GetByID() should return error after deletion")
	}
}

func TestVenueRepository_SearchByGeohash(t *testing.T) {
	repo := NewMockVenueRepository()
	ctx := context.Background()

	// Create venues in San Francisco area
	venue1 := domain.NewVenue(
		"SF Brewery 1",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeBrewery},
		domain.SourceGooglePlaces,
	)
	venue2 := domain.NewVenue(
		"SF Brewery 2",
		domain.GeoPoint{Latitude: 37.7849, Longitude: -122.4094, Geohash: "9q8yym"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeBrewery},
		domain.SourceGooglePlaces,
	)
	_ = repo.Create(ctx, venue1)
	_ = repo.Create(ctx, venue2)

	// Search by geohash prefix
	results, err := repo.SearchByGeohash(ctx, []string{"9q8yy"}, 10)
	if err != nil {
		t.Fatalf("SearchByGeohash() error = %v", err)
	}
	if len(results) != 2 {
		t.Errorf("SearchByGeohash() returned %v venues, want 2", len(results))
	}
}

func TestVenueRepository_SearchByCity(t *testing.T) {
	repo := NewMockVenueRepository()
	ctx := context.Background()

	venue1 := domain.NewVenue(
		"SF Venue",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	venue2 := domain.NewVenue(
		"LA Venue",
		domain.GeoPoint{Latitude: 34.0522, Longitude: -118.2437, Geohash: "9q5ct"},
		domain.Address{City: "Los Angeles", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	_ = repo.Create(ctx, venue1)
	_ = repo.Create(ctx, venue2)

	results, err := repo.SearchByCity(ctx, "San Francisco", "CA", 10)
	if err != nil {
		t.Fatalf("SearchByCity() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("SearchByCity() returned %v venues, want 1", len(results))
	}
	if results[0].Name != "SF Venue" {
		t.Errorf("SearchByCity() returned %v, want SF Venue", results[0].Name)
	}
}

func TestVenueRepository_SearchByType(t *testing.T) {
	repo := NewMockVenueRepository()
	ctx := context.Background()

	brewery := domain.NewVenue(
		"Test Brewery",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeBrewery},
		domain.SourceGooglePlaces,
	)
	winery := domain.NewVenue(
		"Test Winery",
		domain.GeoPoint{Latitude: 37.7849, Longitude: -122.4094, Geohash: "9q8yym"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeWinery},
		domain.SourceGooglePlaces,
	)
	_ = repo.Create(ctx, brewery)
	_ = repo.Create(ctx, winery)

	results, err := repo.SearchByType(ctx, domain.VenueTypeBrewery, 10)
	if err != nil {
		t.Fatalf("SearchByType() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("SearchByType() returned %v venues, want 1", len(results))
	}
	if results[0].Name != "Test Brewery" {
		t.Errorf("SearchByType() returned %v, want Test Brewery", results[0].Name)
	}
}

func TestVenueRepository_GetByExternalID(t *testing.T) {
	repo := NewMockVenueRepository()
	ctx := context.Background()

	venue := domain.NewVenue(
		"Test Venue",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceSongkick,
	)
	venue.SongkickID = "songkick-123"
	_ = repo.Create(ctx, venue)

	retrieved, err := repo.GetByExternalID(ctx, domain.SourceSongkick, "songkick-123")
	if err != nil {
		t.Fatalf("GetByExternalID() error = %v", err)
	}
	if retrieved.ID != venue.ID {
		t.Errorf("GetByExternalID() ID = %v, want %v", retrieved.ID, venue.ID)
	}
	if retrieved.SongkickID != "songkick-123" {
		t.Errorf("GetByExternalID() SongkickID = %v, want songkick-123", retrieved.SongkickID)
	}
}

// Mock repository for testing
type MockVenueRepository struct {
	venues map[string]*domain.Venue
}

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
		return nil, ErrVenueNotFound
	}
	return venue, nil
}

func (r *MockVenueRepository) Update(ctx context.Context, venue *domain.Venue) error {
	if _, ok := r.venues[venue.ID]; !ok {
		return ErrVenueNotFound
	}
	r.venues[venue.ID] = venue
	return nil
}

func (r *MockVenueRepository) Delete(ctx context.Context, id string) error {
	if _, ok := r.venues[id]; !ok {
		return ErrVenueNotFound
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
	return nil, ErrVenueNotFound
}

var ErrVenueNotFound = &VenueNotFoundError{}

type VenueNotFoundError struct{}

func (e *VenueNotFoundError) Error() string {
	return "venue not found"
}
