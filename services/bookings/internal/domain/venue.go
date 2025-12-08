package domain

import (
	"time"

	"github.com/google/uuid"
)

// VenueType represents different types of venues
type VenueType string

const (
	VenueTypeClub       VenueType = "club"
	VenueTypeTheater    VenueType = "theater"
	VenueTypeBrewery    VenueType = "brewery"
	VenueTypeWinery     VenueType = "winery"
	VenueTypeCoffeehouse VenueType = "coffeehouse"
	VenueTypeFestival   VenueType = "festival"
	VenueTypeBar        VenueType = "bar"
	VenueTypeRestaurant VenueType = "restaurant"
	VenueTypeArena      VenueType = "arena"
	VenueTypeOther      VenueType = "other"
)

// PaymentType represents how artists are compensated
type PaymentType string

const (
	PaymentGuarantee  PaymentType = "guarantee"
	PaymentDoorSplit  PaymentType = "door_split"
	PaymentBarTab     PaymentType = "bar_tab"
	PaymentTicketSales PaymentType = "ticket_sales"
	PaymentNone       PaymentType = "none"
)

// DataSource represents where venue data originated
type DataSource string

const (
	SourceSongkick     DataSource = "songkick"
	SourceBandsintown  DataSource = "bandsintown"
	SourceGooglePlaces DataSource = "google_places"
	SourceUserSubmitted DataSource = "user_submitted"
	SourceManual       DataSource = "manual"
)

// GeoPoint represents a geographic coordinate
type GeoPoint struct {
	Latitude  float64 `dynamodbav:"latitude" json:"latitude"`
	Longitude float64 `dynamodbav:"longitude" json:"longitude"`
	Geohash   string  `dynamodbav:"geohash" json:"geohash"` // For efficient spatial queries
}

// Address represents a physical address
type Address struct {
	Street     string `dynamodbav:"street" json:"street"`
	City       string `dynamodbav:"city" json:"city"`
	State      string `dynamodbav:"state" json:"state"`
	PostalCode string `dynamodbav:"postal_code" json:"postal_code"`
	Country    string `dynamodbav:"country" json:"country"`
}

// PayRange represents compensation range
type PayRange struct {
	Min      int         `dynamodbav:"min" json:"min"`
	Max      int         `dynamodbav:"max" json:"max"`
	Currency string      `dynamodbav:"currency" json:"currency"`
	Type     PaymentType `dynamodbav:"type" json:"type"`
	Notes    string      `dynamodbav:"notes,omitempty" json:"notes,omitempty"`
}

// ContactInfo represents venue contact details
type ContactInfo struct {
	Email       string `dynamodbav:"email,omitempty" json:"email,omitempty"`
	Phone       string `dynamodbav:"phone,omitempty" json:"phone,omitempty"`
	Website     string `dynamodbav:"website,omitempty" json:"website,omitempty"`
	BookingURL  string `dynamodbav:"booking_url,omitempty" json:"booking_url,omitempty"`
	ContactName string `dynamodbav:"contact_name,omitempty" json:"contact_name,omitempty"`
}

// DateRange represents availability windows
type DateRange struct {
	Start time.Time `dynamodbav:"start" json:"start"`
	End   time.Time `dynamodbav:"end" json:"end"`
}

// Amenity represents venue features
type Amenity string

const (
	AmenitySoundSystem  Amenity = "sound_system"
	AmenityBackline     Amenity = "backline"
	AmenityGreenRoom    Amenity = "green_room"
	AmenityParking      Amenity = "parking"
	AmenityLoadingDock  Amenity = "loading_dock"
	AmenityLighting     Amenity = "lighting"
	AmenityRecording    Amenity = "recording"
	AmenityLiveStream   Amenity = "live_stream"
	AmenityMerchTable   Amenity = "merch_table"
	AmenityAccessible   Amenity = "accessible"
)

// Venue represents a performance venue
type Venue struct {
	ID           string      `dynamodbav:"id" json:"id"`
	Name         string      `dynamodbav:"name" json:"name"`
	Location     GeoPoint    `dynamodbav:"location" json:"location"`
	Address      Address     `dynamodbav:"address" json:"address"`
	VenueTypes   []VenueType `dynamodbav:"venue_types" json:"venue_types"`
	Capacity     int         `dynamodbav:"capacity" json:"capacity"`
	Genres       []string    `dynamodbav:"genres" json:"genres"`
	PayRange     *PayRange   `dynamodbav:"pay_range,omitempty" json:"pay_range,omitempty"`
	Amenities    []Amenity   `dynamodbav:"amenities" json:"amenities"`
	Photos       []string    `dynamodbav:"photos" json:"photos"`
	ContactInfo  ContactInfo `dynamodbav:"contact_info" json:"contact_info"`
	Availability []DateRange `dynamodbav:"availability,omitempty" json:"availability,omitempty"`
	Rating       float64     `dynamodbav:"rating" json:"rating"`
	ReviewCount  int         `dynamodbav:"review_count" json:"review_count"`
	Description  string      `dynamodbav:"description,omitempty" json:"description,omitempty"`
	
	// External IDs for syncing
	SongkickID     string     `dynamodbav:"songkick_id,omitempty" json:"songkick_id,omitempty"`
	BandsintownID  string     `dynamodbav:"bandsintown_id,omitempty" json:"bandsintown_id,omitempty"`
	GooglePlaceID  string     `dynamodbav:"google_place_id,omitempty" json:"google_place_id,omitempty"`
	
	Source         DataSource `dynamodbav:"source" json:"source"`
	Verified       bool       `dynamodbav:"verified" json:"verified"`
	Active         bool       `dynamodbav:"active" json:"active"`
	
	CreatedAt      time.Time  `dynamodbav:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `dynamodbav:"updated_at" json:"updated_at"`
	LastSyncedAt   *time.Time `dynamodbav:"last_synced_at,omitempty" json:"last_synced_at,omitempty"`
}

// NewVenue creates a new venue
func NewVenue(name string, location GeoPoint, address Address, venueTypes []VenueType, source DataSource) *Venue {
	now := time.Now()
	return &Venue{
		ID:          uuid.New().String(),
		Name:        name,
		Location:    location,
		Address:     address,
		VenueTypes:  venueTypes,
		Amenities:   []Amenity{},
		Photos:      []string{},
		Genres:      []string{},
		Source:      source,
		Active:      true,
		Verified:    false,
		Rating:      0.0,
		ReviewCount: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates venue information
func (v *Venue) Update() {
	v.UpdatedAt = time.Now()
}

// MarkSynced marks the venue as synced from external source
func (v *Venue) MarkSynced() {
	now := time.Now()
	v.LastSyncedAt = &now
	v.UpdatedAt = now
}

// Verify marks the venue as verified
func (v *Venue) Verify() {
	v.Verified = true
	v.UpdatedAt = time.Now()
}

// Deactivate deactivates the venue
func (v *Venue) Deactivate() {
	v.Active = false
	v.UpdatedAt = time.Now()
}
