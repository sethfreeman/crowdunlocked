package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crowdunlocked/services/bookings/internal/domain"
	"github.com/crowdunlocked/services/bookings/internal/repository"
	"github.com/crowdunlocked/services/bookings/internal/service"
	"github.com/go-chi/chi/v5"
)

func TestVenueHandler_Search(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	svc := service.NewVenueService(repo)
	handler := NewVenueHandler(svc)

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
	_ = repo.Create(context.Background(), venue1)
	_ = repo.Create(context.Background(), venue2)

	// Test search by city
	req := httptest.NewRequest(http.MethodGet, "/api/v1/venues/search?city=San+Francisco&state=CA&limit=10", nil)
	w := httptest.NewRecorder()

	handler.Search(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Search() status = %v, want %v", w.Code, http.StatusOK)
	}

	var result domain.VenueSearchResult
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result.Venues) != 2 {
		t.Errorf("Search() returned %v venues, want 2", len(result.Venues))
	}
}

func TestVenueHandler_Search_ByLocation(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	svc := service.NewVenueService(repo)
	handler := NewVenueHandler(svc)

	venue := domain.NewVenue(
		"Test Venue",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	_ = repo.Create(context.Background(), venue)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/venues/search?lat=37.7749&lng=-122.4194&radius=5&limit=10", nil)
	w := httptest.NewRecorder()

	handler.Search(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Search() status = %v, want %v", w.Code, http.StatusOK)
	}

	var result domain.VenueSearchResult
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result.Venues) != 1 {
		t.Errorf("Search() returned %v venues, want 1", len(result.Venues))
	}
}

func TestVenueHandler_Search_WithFilters(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	svc := service.NewVenueService(repo)
	handler := NewVenueHandler(svc)

	brewery := domain.NewVenue(
		"Test Brewery",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeBrewery},
		domain.SourceGooglePlaces,
	)
	brewery.Capacity = 100

	winery := domain.NewVenue(
		"Test Winery",
		domain.GeoPoint{Latitude: 37.7849, Longitude: -122.4094, Geohash: "9q8yym"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeWinery},
		domain.SourceGooglePlaces,
	)
	winery.Capacity = 200

	_ = repo.Create(context.Background(), brewery)
	_ = repo.Create(context.Background(), winery)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/venues/search?city=San+Francisco&state=CA&venue_types=brewery&min_capacity=50&max_capacity=150&limit=10", nil)
	w := httptest.NewRecorder()

	handler.Search(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Search() status = %v, want %v", w.Code, http.StatusOK)
	}

	var result domain.VenueSearchResult
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result.Venues) != 1 {
		t.Errorf("Search() returned %v venues, want 1", len(result.Venues))
	}
	if result.Venues[0].Venue.Name != "Test Brewery" {
		t.Errorf("Search() returned %v, want Test Brewery", result.Venues[0].Venue.Name)
	}
}

func TestVenueHandler_GetByID(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	svc := service.NewVenueService(repo)
	handler := NewVenueHandler(svc)

	venue := domain.NewVenue(
		"Test Venue",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	_ = repo.Create(context.Background(), venue)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/venues/"+venue.ID, nil)
	w := httptest.NewRecorder()

	// Set up chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", venue.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetByID() status = %v, want %v", w.Code, http.StatusOK)
	}

	var result domain.Venue
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.ID != venue.ID {
		t.Errorf("GetByID() ID = %v, want %v", result.ID, venue.ID)
	}
}

func TestVenueHandler_GetByID_NotFound(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	svc := service.NewVenueService(repo)
	handler := NewVenueHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/venues/nonexistent", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "nonexistent")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("GetByID() status = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestVenueHandler_Create(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	svc := service.NewVenueService(repo)
	handler := NewVenueHandler(svc)

	reqBody := CreateVenueRequest{
		Name: "New Brewery",
		Location: LocationRequest{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
		Address: AddressRequest{
			Street:     "123 Main St",
			City:       "San Francisco",
			State:      "CA",
			PostalCode: "94102",
			Country:    "US",
		},
		VenueTypes: []string{"brewery"},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/venues", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Create() status = %v, want %v. Body: %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var result domain.Venue
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.Name != "New Brewery" {
		t.Errorf("Create() name = %v, want New Brewery", result.Name)
	}
	if result.Location.Geohash == "" {
		t.Error("Create() should generate geohash")
	}
}

func TestVenueHandler_Update(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	svc := service.NewVenueService(repo)
	handler := NewVenueHandler(svc)

	venue := domain.NewVenue(
		"Original Name",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	_ = repo.Create(context.Background(), venue)

	reqBody := UpdateVenueRequest{
		Name:     stringPtr("Updated Name"),
		Capacity: intPtr(200),
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/venues/"+venue.ID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", venue.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.Update(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Update() status = %v, want %v. Body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var result domain.Venue
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.Name != "Updated Name" {
		t.Errorf("Update() name = %v, want Updated Name", result.Name)
	}
	if result.Capacity != 200 {
		t.Errorf("Update() capacity = %v, want 200", result.Capacity)
	}
}

func TestVenueHandler_Delete(t *testing.T) {
	repo := repository.NewMockVenueRepository()
	svc := service.NewVenueService(repo)
	handler := NewVenueHandler(svc)

	venue := domain.NewVenue(
		"To Delete",
		domain.GeoPoint{Latitude: 37.7749, Longitude: -122.4194, Geohash: "9q8yyk"},
		domain.Address{City: "San Francisco", State: "CA", Country: "US"},
		[]domain.VenueType{domain.VenueTypeClub},
		domain.SourceUserSubmitted,
	)
	_ = repo.Create(context.Background(), venue)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/venues/"+venue.ID, nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", venue.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Delete() status = %v, want %v", w.Code, http.StatusNoContent)
	}

	// Verify deletion
	_, err := repo.GetByID(context.Background(), venue.ID)
	if err == nil {
		t.Error("Delete() should delete the venue")
	}
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
