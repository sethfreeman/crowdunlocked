package service

import (
	"context"
	"testing"

	"github.com/crowdunlocked/services/bookings/internal/domain"
	"github.com/crowdunlocked/services/bookings/internal/repository"
)

func TestVenueService_Search_ByLocation(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	service := NewVenueService(repo)
	ctx := context.Background()

	// Create test venues
	venue1 := domain.NewVenue(
		"SF Brewery",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeBrewery},
		domain.SourceGooglePlaces,
	)
	venue2 := domain.NewVenue(
		"SF Winery",
		domain.GeoPoint{Latitude: 37.7849, Longitude: -122.4094, Geohash: "9q8yym"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeWinery},
		domain.SourceGooglePlaces,
	)
	_ = repo.Create(ctx, venue1)
	_ = repo.Create(ctx, venue2)

	// Search by location
	criteria := &domain.VenueSearchCriteria{
		Location: &domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194},
		RadiusKm: 5.0,
		Limit:    10,
	}

	result, err := service.Search(ctx, criteria)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(result.Venues) != 2 {
		t.Errorf("Search() returned %v venues, want 2", len(result.Venues))
	}
}

func TestVenueService_Search_WithDistanceCalculation(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	service := NewVenueService(repo)
	ctx := context.Background()

	venue := domain.NewVenue(
		"Test Venue",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	_ = repo.Create(ctx, venue)

	criteria := &domain.VenueSearchCriteria{
		Location: &domain.GeoPoint{Latitude: 37.7849, Longitude: -122.4094},
		RadiusKm: 10.0,
		Limit:    10,
	}

	result, err := service.Search(ctx, criteria)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(result.Venues) != 1 {
		t.Errorf("Search() returned %v venues, want 1", len(result.Venues))
	}
	if result.Venues[0].DistanceKm <= 0 {
		t.Error("Search() should calculate distance")
	}
}

func TestVenueService_Search_ByCity(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	service := NewVenueService(repo)
	ctx := context.Background()

	venue := domain.NewVenue(
		"SF Venue",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	_ = repo.Create(ctx, venue)

	criteria := &domain.VenueSearchCriteria{
		City:  "San Francisco",
		State: "CA",
		Limit: 10,
	}

	result, err := service.Search(ctx, criteria)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(result.Venues) != 1 {
		t.Errorf("Search() returned %v venues, want 1", len(result.Venues))
	}
}

func TestVenueService_Search_ByVenueType(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	service := NewVenueService(repo)
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

	criteria := &domain.VenueSearchCriteria{
		VenueTypes: []domain.VenueType{domain.VenueTypeBrewery},
		Limit:      10,
	}

	result, err := service.Search(ctx, criteria)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(result.Venues) != 1 {
		t.Errorf("Search() returned %v venues, want 1", len(result.Venues))
	}
	if result.Venues[0].Venue.Name != "Test Brewery" {
		t.Errorf("Search() returned %v, want Test Brewery", result.Venues[0].Venue.Name)
	}
}

func TestVenueService_Search_WithFilters(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	service := NewVenueService(repo)
	ctx := context.Background()

	// Create venues with different capacities
	smallVenue := domain.NewVenue(
		"Small Venue",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	smallVenue.Capacity = 50

	largeVenue := domain.NewVenue(
		"Large Venue",
		domain.GeoPoint{Latitude: 37.7849, Longitude: -122.4094, Geohash: "9q8yym"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	largeVenue.Capacity = 500

	_ = repo.Create(ctx, smallVenue)
	_ = repo.Create(ctx, largeVenue)

	// Search with capacity filter
	criteria := &domain.VenueSearchCriteria{
		City:        "San Francisco",
		State:       "CA",
		MinCapacity: 100,
		MaxCapacity: 1000,
		Limit:       10,
	}

	result, err := service.Search(ctx, criteria)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(result.Venues) != 1 {
		t.Errorf("Search() returned %v venues, want 1", len(result.Venues))
	}
	if result.Venues[0].Venue.Name != "Large Venue" {
		t.Errorf("Search() returned %v, want Large Venue", result.Venues[0].Venue.Name)
	}
}

func TestVenueService_Search_WithGenreFilter(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	service := NewVenueService(repo)
	ctx := context.Background()

	rockVenue := domain.NewVenue(
		"Rock Club",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	rockVenue.Genres = []string{"rock", "metal"}

	jazzVenue := domain.NewVenue(
		"Jazz Club",
		domain.GeoPoint{Latitude: 37.7849, Longitude: -122.4094, Geohash: "9q8yym"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	jazzVenue.Genres = []string{"jazz", "blues"}

	_ = repo.Create(ctx, rockVenue)
	_ = repo.Create(ctx, jazzVenue)

	criteria := &domain.VenueSearchCriteria{
		City:   "San Francisco",
		State:  "CA",
		Genres: []string{"rock"},
		Limit:  10,
	}

	result, err := service.Search(ctx, criteria)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(result.Venues) != 1 {
		t.Errorf("Search() returned %v venues, want 1", len(result.Venues))
	}
	if result.Venues[0].Venue.Name != "Rock Club" {
		t.Errorf("Search() returned %v, want Rock Club", result.Venues[0].Venue.Name)
	}
}

func TestVenueService_Search_WithRatingFilter(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	service := NewVenueService(repo)
	ctx := context.Background()

	lowRated := domain.NewVenue(
		"Low Rated",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	lowRated.Rating = 3.0

	highRated := domain.NewVenue(
		"High Rated",
		domain.GeoPoint{Latitude: 37.7849, Longitude: -122.4094, Geohash: "9q8yym"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	highRated.Rating = 4.5

	_ = repo.Create(ctx, lowRated)
	_ = repo.Create(ctx, highRated)

	criteria := &domain.VenueSearchCriteria{
		City:      "San Francisco",
		State:     "CA",
		MinRating: 4.0,
		Limit:     10,
	}

	result, err := service.Search(ctx, criteria)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(result.Venues) != 1 {
		t.Errorf("Search() returned %v venues, want 1", len(result.Venues))
	}
	if result.Venues[0].Venue.Name != "High Rated" {
		t.Errorf("Search() returned %v, want High Rated", result.Venues[0].Venue.Name)
	}
}

func TestVenueService_Search_VerifiedOnly(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	service := NewVenueService(repo)
	ctx := context.Background()

	unverified := domain.NewVenue(
		"Unverified",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)

	verified := domain.NewVenue(
		"Verified",
		domain.GeoPoint{Latitude: 37.7849, Longitude: -122.4094, Geohash: "9q8yym"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	verified.Verify()

	_ = repo.Create(ctx, unverified)
	_ = repo.Create(ctx, verified)

	criteria := &domain.VenueSearchCriteria{
		City:         "San Francisco",
		State:        "CA",
		VerifiedOnly: true,
		Limit:        10,
	}

	result, err := service.Search(ctx, criteria)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(result.Venues) != 1 {
		t.Errorf("Search() returned %v venues, want 1", len(result.Venues))
	}
	if result.Venues[0].Venue.Name != "Verified" {
		t.Errorf("Search() returned %v, want Verified", result.Venues[0].Venue.Name)
	}
}

func TestVenueService_Search_SortByDistance(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	service := NewVenueService(repo)
	ctx := context.Background()

	near := domain.NewVenue(
		"Near Venue",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)

	far := domain.NewVenue(
		"Far Venue",
		domain.GeoPoint{Latitude: 37.8049, Longitude: -122.3894, Geohash: "9q8yz"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)

	_ = repo.Create(ctx, near)
	_ = repo.Create(ctx, far)

	criteria := &domain.VenueSearchCriteria{
		Location: &domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194},
		RadiusKm: 10.0,
		SortBy:   domain.SortByDistance,
		SortOrder: domain.SortAsc,
		Limit:    10,
	}

	result, err := service.Search(ctx, criteria)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(result.Venues) != 2 {
		t.Errorf("Search() returned %v venues, want 2", len(result.Venues))
	}
	// First venue should be closer
	if result.Venues[0].Venue.Name != "Near Venue" {
		t.Errorf("Search() first venue = %v, want Near Venue", result.Venues[0].Venue.Name)
	}
	if result.Venues[0].DistanceKm >= result.Venues[1].DistanceKm {
		t.Error("Search() should sort by distance ascending")
	}
}

func TestVenueService_GetByID(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	service := NewVenueService(repo)
	ctx := context.Background()

	venue := domain.NewVenue(
		"Test Venue",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	_ = repo.Create(ctx, venue)

	retrieved, err := service.GetByID(ctx, venue.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if retrieved.ID != venue.ID {
		t.Errorf("GetByID() ID = %v, want %v", retrieved.ID, venue.ID)
	}
}

func TestVenueService_Create(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	service := NewVenueService(repo)
	ctx := context.Background()

	venue := domain.NewVenue(
		"New Venue",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeBrewery},
		domain.SourceUserSubmitted,
	)

	err := service.Create(ctx, venue)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify creation
	retrieved, _ := service.GetByID(ctx, venue.ID)
	if retrieved.Name != "New Venue" {
		t.Errorf("Create() name = %v, want New Venue", retrieved.Name)
	}
}

func TestVenueService_Update(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	service := NewVenueService(repo)
	ctx := context.Background()

	venue := domain.NewVenue(
		"Original",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	_ = service.Create(ctx, venue)

	venue.Name = "Updated"
	err := service.Update(ctx, venue)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	retrieved, _ := service.GetByID(ctx, venue.ID)
	if retrieved.Name != "Updated" {
		t.Errorf("Update() name = %v, want Updated", retrieved.Name)
	}
}

func TestVenueService_Delete(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	service := NewVenueService(repo)
	ctx := context.Background()

	venue := domain.NewVenue(
		"To Delete",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	_ = service.Create(ctx, venue)

	err := service.Delete(ctx, venue.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = service.GetByID(ctx, venue.ID)
	if err == nil {
		t.Error("GetByID() should return error after deletion")
	}
}
