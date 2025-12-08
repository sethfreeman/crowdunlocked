package domain

// VenueSearchCriteria represents search filters for venues
type VenueSearchCriteria struct {
	// Geographic filters
	Location      *GeoPoint
	RadiusKm      float64
	City          string
	State         string
	Country       string
	
	// Venue characteristics
	VenueTypes    []VenueType
	MinCapacity   int
	MaxCapacity   int
	Genres        []string
	Amenities     []Amenity
	
	// Payment filters
	MinPay        int
	MaxPay        int
	PaymentTypes  []PaymentType
	
	// Availability
	AvailableFrom *DateRange
	
	// Quality filters
	MinRating     float64
	VerifiedOnly  bool
	ActiveOnly    bool
	
	// Pagination
	Limit         int
	Offset        int
	
	// Sorting
	SortBy        VenueSortField
	SortOrder     SortOrder
}

type VenueSortField string

const (
	SortByDistance   VenueSortField = "distance"
	SortByRating     VenueSortField = "rating"
	SortByCapacity   VenueSortField = "capacity"
	SortByPay        VenueSortField = "pay"
	SortByName       VenueSortField = "name"
	SortByCreatedAt  VenueSortField = "created_at"
)

type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// VenueSearchResult represents search results with metadata
type VenueSearchResult struct {
	Venues     []*VenueWithDistance `json:"venues"`
	Total      int                  `json:"total"`
	Limit      int                  `json:"limit"`
	Offset     int                  `json:"offset"`
	HasMore    bool                 `json:"has_more"`
}

// VenueWithDistance includes distance from search point
type VenueWithDistance struct {
	*Venue
	DistanceKm float64 `json:"distance_km"`
}
