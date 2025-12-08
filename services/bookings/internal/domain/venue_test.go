package domain

import (
	"testing"
	"time"
)

func TestNewVenue(t *testing.T) {
	location := GeoPoint{
		Latitude:  37.7749,
		Longitude: -122.4194,
		Geohash:   "9q8yyk",
	}
	address := Address{
		Street:     "123 Main St",
		City:       "San Francisco",
		State:      "CA",
		PostalCode: "94102",
		Country:    "US",
	}
	venueTypes := []VenueType{VenueTypeBrewery}

	venue := NewVenue("Test Brewery", location, address, venueTypes, SourceGooglePlaces)

	if venue.ID == "" {
		t.Error("NewVenue() should generate an ID")
	}
	if venue.Name != "Test Brewery" {
		t.Errorf("NewVenue() name = %v, want %v", venue.Name, "Test Brewery")
	}
	if venue.Location.Latitude != location.Latitude {
		t.Errorf("NewVenue() location = %v, want %v", venue.Location, location)
	}
	if venue.Address.City != "San Francisco" {
		t.Errorf("NewVenue() city = %v, want %v", venue.Address.City, "San Francisco")
	}
	if len(venue.VenueTypes) != 1 || venue.VenueTypes[0] != VenueTypeBrewery {
		t.Errorf("NewVenue() venue types = %v, want [brewery]", venue.VenueTypes)
	}
	if venue.Source != SourceGooglePlaces {
		t.Errorf("NewVenue() source = %v, want %v", venue.Source, SourceGooglePlaces)
	}
	if !venue.Active {
		t.Error("NewVenue() should be active by default")
	}
	if venue.Verified {
		t.Error("NewVenue() should not be verified by default")
	}
	if venue.CreatedAt.IsZero() {
		t.Error("NewVenue() should set CreatedAt")
	}
	if venue.UpdatedAt.IsZero() {
		t.Error("NewVenue() should set UpdatedAt")
	}
}

func TestVenue_Update(t *testing.T) {
	venue := &Venue{
		ID:        "test-id",
		Name:      "Test Venue",
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	oldUpdatedAt := venue.UpdatedAt
	time.Sleep(10 * time.Millisecond)
	venue.Update()

	if !venue.UpdatedAt.After(oldUpdatedAt) {
		t.Error("Update() should update UpdatedAt timestamp")
	}
}

func TestVenue_MarkSynced(t *testing.T) {
	venue := &Venue{
		ID:           "test-id",
		Name:         "Test Venue",
		LastSyncedAt: nil,
	}

	venue.MarkSynced()

	if venue.LastSyncedAt == nil {
		t.Error("MarkSynced() should set LastSyncedAt")
	}
	if venue.LastSyncedAt.IsZero() {
		t.Error("MarkSynced() should set LastSyncedAt to current time")
	}
}

func TestVenue_Verify(t *testing.T) {
	venue := &Venue{
		ID:       "test-id",
		Name:     "Test Venue",
		Verified: false,
	}

	venue.Verify()

	if !venue.Verified {
		t.Error("Verify() should set Verified to true")
	}
}

func TestVenue_Deactivate(t *testing.T) {
	venue := &Venue{
		ID:     "test-id",
		Name:   "Test Venue",
		Active: true,
	}

	venue.Deactivate()

	if venue.Active {
		t.Error("Deactivate() should set Active to false")
	}
}

func TestVenueTypes(t *testing.T) {
	// Test that all venue type constants are defined
	types := []VenueType{
		VenueTypeClub,
		VenueTypeTheater,
		VenueTypeBrewery,
		VenueTypeWinery,
		VenueTypeCoffeehouse,
		VenueTypeFestival,
		VenueTypeBar,
		VenueTypeRestaurant,
		VenueTypeArena,
		VenueTypeOther,
	}

	for _, vt := range types {
		if string(vt) == "" {
			t.Errorf("VenueType should not be empty: %v", vt)
		}
	}
}

func TestPaymentTypes(t *testing.T) {
	// Test that all payment type constants are defined
	types := []PaymentType{
		PaymentGuarantee,
		PaymentDoorSplit,
		PaymentBarTab,
		PaymentTicketSales,
		PaymentNone,
	}

	for _, pt := range types {
		if string(pt) == "" {
			t.Errorf("PaymentType should not be empty: %v", pt)
		}
	}
}

func TestAmenities(t *testing.T) {
	// Test that all amenity constants are defined
	amenities := []Amenity{
		AmenitySoundSystem,
		AmenityBackline,
		AmenityGreenRoom,
		AmenityParking,
		AmenityLoadingDock,
		AmenityLighting,
		AmenityRecording,
		AmenityLiveStream,
		AmenityMerchTable,
		AmenityAccessible,
	}

	for _, a := range amenities {
		if string(a) == "" {
			t.Errorf("Amenity should not be empty: %v", a)
		}
	}
}

func TestDataSources(t *testing.T) {
	// Test that all data source constants are defined
	sources := []DataSource{
		SourceSongkick,
		SourceBandsintown,
		SourceGooglePlaces,
		SourceUserSubmitted,
		SourceManual,
	}

	for _, s := range sources {
		if string(s) == "" {
			t.Errorf("DataSource should not be empty: %v", s)
		}
	}
}
