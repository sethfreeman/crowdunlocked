package domain

import (
	"testing"
)

func TestVenueSearchCriteria_Defaults(t *testing.T) {
	criteria := &VenueSearchCriteria{
		Limit:  20,
		Offset: 0,
	}

	if criteria.Limit != 20 {
		t.Errorf("Default limit = %v, want 20", criteria.Limit)
	}
	if criteria.Offset != 0 {
		t.Errorf("Default offset = %v, want 0", criteria.Offset)
	}
}

func TestVenueSearchCriteria_GeographicFilters(t *testing.T) {
	location := &GeoPoint{
		Latitude:  37.7749,
		Longitude: -122.4194,
	}

	criteria := &VenueSearchCriteria{
		Location: location,
		RadiusKm: 10.0,
		City:     "San Francisco",
		State:    "CA",
		Country:  "US",
	}

	if criteria.Location.Latitude != 37.7749 {
		t.Errorf("Location latitude = %v, want 37.7749", criteria.Location.Latitude)
	}
	if criteria.RadiusKm != 10.0 {
		t.Errorf("RadiusKm = %v, want 10.0", criteria.RadiusKm)
	}
	if criteria.City != "San Francisco" {
		t.Errorf("City = %v, want San Francisco", criteria.City)
	}
}

func TestVenueSearchCriteria_VenueFilters(t *testing.T) {
	criteria := &VenueSearchCriteria{
		VenueTypes:  []VenueType{VenueTypeBrewery, VenueTypeWinery},
		MinCapacity: 50,
		MaxCapacity: 500,
		Genres:      []string{"rock", "indie"},
		Amenities:   []Amenity{AmenitySoundSystem, AmenityParking},
	}

	if len(criteria.VenueTypes) != 2 {
		t.Errorf("VenueTypes count = %v, want 2", len(criteria.VenueTypes))
	}
	if criteria.MinCapacity != 50 {
		t.Errorf("MinCapacity = %v, want 50", criteria.MinCapacity)
	}
	if criteria.MaxCapacity != 500 {
		t.Errorf("MaxCapacity = %v, want 500", criteria.MaxCapacity)
	}
	if len(criteria.Genres) != 2 {
		t.Errorf("Genres count = %v, want 2", len(criteria.Genres))
	}
	if len(criteria.Amenities) != 2 {
		t.Errorf("Amenities count = %v, want 2", len(criteria.Amenities))
	}
}

func TestVenueSearchCriteria_PaymentFilters(t *testing.T) {
	criteria := &VenueSearchCriteria{
		MinPay:       100,
		MaxPay:       1000,
		PaymentTypes: []PaymentType{PaymentGuarantee, PaymentDoorSplit},
	}

	if criteria.MinPay != 100 {
		t.Errorf("MinPay = %v, want 100", criteria.MinPay)
	}
	if criteria.MaxPay != 1000 {
		t.Errorf("MaxPay = %v, want 1000", criteria.MaxPay)
	}
	if len(criteria.PaymentTypes) != 2 {
		t.Errorf("PaymentTypes count = %v, want 2", len(criteria.PaymentTypes))
	}
}

func TestVenueSearchCriteria_QualityFilters(t *testing.T) {
	criteria := &VenueSearchCriteria{
		MinRating:    4.0,
		VerifiedOnly: true,
		ActiveOnly:   true,
	}

	if criteria.MinRating != 4.0 {
		t.Errorf("MinRating = %v, want 4.0", criteria.MinRating)
	}
	if !criteria.VerifiedOnly {
		t.Error("VerifiedOnly should be true")
	}
	if !criteria.ActiveOnly {
		t.Error("ActiveOnly should be true")
	}
}

func TestVenueSearchCriteria_Sorting(t *testing.T) {
	criteria := &VenueSearchCriteria{
		SortBy:    SortByDistance,
		SortOrder: SortAsc,
	}

	if criteria.SortBy != SortByDistance {
		t.Errorf("SortBy = %v, want %v", criteria.SortBy, SortByDistance)
	}
	if criteria.SortOrder != SortAsc {
		t.Errorf("SortOrder = %v, want %v", criteria.SortOrder, SortAsc)
	}
}

func TestVenueSortFields(t *testing.T) {
	// Test that all sort field constants are defined
	fields := []VenueSortField{
		SortByDistance,
		SortByRating,
		SortByCapacity,
		SortByPay,
		SortByName,
		SortByCreatedAt,
	}

	for _, f := range fields {
		if string(f) == "" {
			t.Errorf("VenueSortField should not be empty: %v", f)
		}
	}
}

func TestSortOrders(t *testing.T) {
	// Test that sort order constants are defined
	if SortAsc != "asc" {
		t.Errorf("SortAsc = %v, want asc", SortAsc)
	}
	if SortDesc != "desc" {
		t.Errorf("SortDesc = %v, want desc", SortDesc)
	}
}

func TestVenueSearchResult(t *testing.T) {
	venues := []*VenueWithDistance{
		{
			Venue: &Venue{
				ID:   "venue1",
				Name: "Test Venue 1",
			},
			DistanceKm: 2.5,
		},
		{
			Venue: &Venue{
				ID:   "venue2",
				Name: "Test Venue 2",
			},
			DistanceKm: 5.0,
		},
	}

	result := &VenueSearchResult{
		Venues:  venues,
		Total:   100,
		Limit:   20,
		Offset:  0,
		HasMore: true,
	}

	if len(result.Venues) != 2 {
		t.Errorf("Venues count = %v, want 2", len(result.Venues))
	}
	if result.Total != 100 {
		t.Errorf("Total = %v, want 100", result.Total)
	}
	if result.Limit != 20 {
		t.Errorf("Limit = %v, want 20", result.Limit)
	}
	if result.Offset != 0 {
		t.Errorf("Offset = %v, want 0", result.Offset)
	}
	if !result.HasMore {
		t.Error("HasMore should be true")
	}
}

func TestVenueWithDistance(t *testing.T) {
	venue := &Venue{
		ID:   "test-venue",
		Name: "Test Venue",
	}

	venueWithDistance := &VenueWithDistance{
		Venue:      venue,
		DistanceKm: 3.5,
	}

	if venueWithDistance.Venue.ID != "test-venue" {
		t.Errorf("Venue ID = %v, want test-venue", venueWithDistance.Venue.ID)
	}
	if venueWithDistance.DistanceKm != 3.5 {
		t.Errorf("DistanceKm = %v, want 3.5", venueWithDistance.DistanceKm)
	}
}
