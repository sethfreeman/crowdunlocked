package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/crowdunlocked/services/bookings/internal/domain"
	"github.com/crowdunlocked/services/bookings/internal/repository"
)

// VenueService provides business logic for venue operations
type VenueService struct {
	repo repository.VenueRepository
}

// NewVenueService creates a new venue service
func NewVenueService(repo repository.VenueRepository) *VenueService {
	return &VenueService{
		repo: repo,
	}
}

// Search searches for venues based on criteria
func (s *VenueService) Search(ctx context.Context, criteria *domain.VenueSearchCriteria) (*domain.VenueSearchResult, error) {
	var venues []*domain.Venue
	var err error

	// Determine search strategy based on criteria
	if criteria.Location != nil && criteria.RadiusKm > 0 {
		// Geospatial search
		venues, err = s.searchByLocation(ctx, criteria)
	} else if criteria.City != "" && criteria.State != "" {
		// City search
		venues, err = s.repo.SearchByCity(ctx, criteria.City, criteria.State, criteria.Limit*2) // Get more for filtering
	} else if len(criteria.VenueTypes) > 0 {
		// Type search
		venues, err = s.searchByTypes(ctx, criteria)
	} else {
		return nil, fmt.Errorf("search criteria must include location, city, or venue type")
	}

	if err != nil {
		return nil, err
	}

	// Apply filters
	venues = s.applyFilters(venues, criteria)

	// Calculate distances if location provided
	venuesWithDistance := s.calculateDistances(venues, criteria.Location)

	// Filter by radius if location provided
	if criteria.Location != nil && criteria.RadiusKm > 0 {
		venuesWithDistance = s.filterByRadius(venuesWithDistance, criteria.RadiusKm)
	}

	// Sort results
	s.sortResults(venuesWithDistance, criteria)

	// Apply pagination
	total := len(venuesWithDistance)
	start := criteria.Offset
	end := start + criteria.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedVenues := venuesWithDistance[start:end]

	return &domain.VenueSearchResult{
		Venues:  paginatedVenues,
		Total:   total,
		Limit:   criteria.Limit,
		Offset:  criteria.Offset,
		HasMore: end < total,
	}, nil
}

// searchByLocation performs geospatial search using geohash
func (s *VenueService) searchByLocation(ctx context.Context, criteria *domain.VenueSearchCriteria) ([]*domain.Venue, error) {
	// Get geohash prefixes for the search area
	geohashPrefixes := domain.GetGeohashPrefixes(
		criteria.Location.Latitude,
		criteria.Location.Longitude,
		criteria.RadiusKm,
	)

	// Query by geohash prefixes
	return s.repo.SearchByGeohash(ctx, geohashPrefixes, criteria.Limit*3) // Get more for filtering
}

// searchByTypes searches by multiple venue types
func (s *VenueService) searchByTypes(ctx context.Context, criteria *domain.VenueSearchCriteria) ([]*domain.Venue, error) {
	allVenues := make([]*domain.Venue, 0)
	seen := make(map[string]bool)

	for _, venueType := range criteria.VenueTypes {
		venues, err := s.repo.SearchByType(ctx, venueType, criteria.Limit*2)
		if err != nil {
			return nil, err
		}

		// Deduplicate
		for _, venue := range venues {
			if !seen[venue.ID] {
				allVenues = append(allVenues, venue)
				seen[venue.ID] = true
			}
		}
	}

	return allVenues, nil
}

// applyFilters applies all filter criteria to venues
func (s *VenueService) applyFilters(venues []*domain.Venue, criteria *domain.VenueSearchCriteria) []*domain.Venue {
	filtered := make([]*domain.Venue, 0)

	for _, venue := range venues {
		if s.matchesFilters(venue, criteria) {
			filtered = append(filtered, venue)
		}
	}

	return filtered
}

// matchesFilters checks if a venue matches all filter criteria
func (s *VenueService) matchesFilters(venue *domain.Venue, criteria *domain.VenueSearchCriteria) bool {
	// Capacity filter
	if criteria.MinCapacity > 0 && venue.Capacity < criteria.MinCapacity {
		return false
	}
	if criteria.MaxCapacity > 0 && venue.Capacity > criteria.MaxCapacity {
		return false
	}

	// Genre filter
	if len(criteria.Genres) > 0 && !s.hasAnyGenre(venue.Genres, criteria.Genres) {
		return false
	}

	// Amenities filter
	if len(criteria.Amenities) > 0 && !s.hasAllAmenities(venue.Amenities, criteria.Amenities) {
		return false
	}

	// Payment filter
	if criteria.MinPay > 0 && venue.PayRange != nil && venue.PayRange.Max < criteria.MinPay {
		return false
	}
	if criteria.MaxPay > 0 && venue.PayRange != nil && venue.PayRange.Min > criteria.MaxPay {
		return false
	}

	// Rating filter
	if criteria.MinRating > 0 && venue.Rating < criteria.MinRating {
		return false
	}

	// Verified filter
	if criteria.VerifiedOnly && !venue.Verified {
		return false
	}

	// Active filter
	if criteria.ActiveOnly && !venue.Active {
		return false
	}

	return true
}

// hasAnyGenre checks if venue has any of the requested genres
func (s *VenueService) hasAnyGenre(venueGenres, requestedGenres []string) bool {
	genreMap := make(map[string]bool)
	for _, g := range venueGenres {
		genreMap[g] = true
	}

	for _, g := range requestedGenres {
		if genreMap[g] {
			return true
		}
	}

	return false
}

// hasAllAmenities checks if venue has all requested amenities
func (s *VenueService) hasAllAmenities(venueAmenities, requestedAmenities []domain.Amenity) bool {
	amenityMap := make(map[domain.Amenity]bool)
	for _, a := range venueAmenities {
		amenityMap[a] = true
	}

	for _, a := range requestedAmenities {
		if !amenityMap[a] {
			return false
		}
	}

	return true
}

// calculateDistances calculates distance from search location to each venue
func (s *VenueService) calculateDistances(venues []*domain.Venue, location *domain.GeoPoint) []*domain.VenueWithDistance {
	result := make([]*domain.VenueWithDistance, len(venues))

	for i, venue := range venues {
		distance := 0.0
		if location != nil {
			distance = domain.CalculateDistance(
				location.Latitude,
				location.Longitude,
				venue.Location.Latitude,
				venue.Location.Longitude,
			)
		}

		result[i] = &domain.VenueWithDistance{
			Venue:      venue,
			DistanceKm: distance,
		}
	}

	return result
}

// filterByRadius filters venues within the specified radius
func (s *VenueService) filterByRadius(venues []*domain.VenueWithDistance, radiusKm float64) []*domain.VenueWithDistance {
	filtered := make([]*domain.VenueWithDistance, 0)

	for _, v := range venues {
		if v.DistanceKm <= radiusKm {
			filtered = append(filtered, v)
		}
	}

	return filtered
}

// sortResults sorts venues based on criteria
func (s *VenueService) sortResults(venues []*domain.VenueWithDistance, criteria *domain.VenueSearchCriteria) {
	if criteria.SortBy == "" {
		criteria.SortBy = domain.SortByDistance
	}
	if criteria.SortOrder == "" {
		criteria.SortOrder = domain.SortAsc
	}

	sort.Slice(venues, func(i, j int) bool {
		var less bool

		switch criteria.SortBy {
		case domain.SortByDistance:
			less = venues[i].DistanceKm < venues[j].DistanceKm
		case domain.SortByRating:
			less = venues[i].Venue.Rating < venues[j].Venue.Rating
		case domain.SortByCapacity:
			less = venues[i].Venue.Capacity < venues[j].Venue.Capacity
		case domain.SortByPay:
			iPay := 0
			if venues[i].Venue.PayRange != nil {
				iPay = venues[i].Venue.PayRange.Max
			}
			jPay := 0
			if venues[j].Venue.PayRange != nil {
				jPay = venues[j].Venue.PayRange.Max
			}
			less = iPay < jPay
		case domain.SortByName:
			less = venues[i].Venue.Name < venues[j].Venue.Name
		case domain.SortByCreatedAt:
			less = venues[i].Venue.CreatedAt.Before(venues[j].Venue.CreatedAt)
		default:
			less = venues[i].DistanceKm < venues[j].DistanceKm
		}

		if criteria.SortOrder == domain.SortDesc {
			return !less
		}
		return less
	})
}

// GetByID retrieves a venue by ID
func (s *VenueService) GetByID(ctx context.Context, id string) (*domain.Venue, error) {
	return s.repo.GetByID(ctx, id)
}

// Create creates a new venue
func (s *VenueService) Create(ctx context.Context, venue *domain.Venue) error {
	// Generate geohash if not provided
	if venue.Location.Geohash == "" {
		venue.Location.Geohash = domain.EncodeGeohash(
			venue.Location.Latitude,
			venue.Location.Longitude,
			6, // Default precision
		)
	}

	return s.repo.Create(ctx, venue)
}

// Update updates an existing venue
func (s *VenueService) Update(ctx context.Context, venue *domain.Venue) error {
	venue.Update()

	// Update geohash if location changed
	if venue.Location.Geohash == "" {
		venue.Location.Geohash = domain.EncodeGeohash(
			venue.Location.Latitude,
			venue.Location.Longitude,
			6,
		)
	}

	return s.repo.Update(ctx, venue)
}

// Delete deletes a venue
func (s *VenueService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// GetByExternalID retrieves a venue by external source ID
func (s *VenueService) GetByExternalID(ctx context.Context, source domain.DataSource, externalID string) (*domain.Venue, error) {
	return s.repo.GetByExternalID(ctx, source, externalID)
}
