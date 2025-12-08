package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/crowdunlocked/services/bookings/internal/domain"
	"github.com/crowdunlocked/services/bookings/internal/service"
	"github.com/go-chi/chi/v5"
)

type VenueHandler struct {
	service *service.VenueService
}

func NewVenueHandler(service *service.VenueService) *VenueHandler {
	return &VenueHandler{service: service}
}

// LocationRequest represents location in API requests
type LocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// AddressRequest represents address in API requests
type AddressRequest struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// CreateVenueRequest represents the request body for creating a venue
type CreateVenueRequest struct {
	Name        string          `json:"name"`
	Location    LocationRequest `json:"location"`
	Address     AddressRequest  `json:"address"`
	VenueTypes  []string        `json:"venue_types"`
	Capacity    int             `json:"capacity,omitempty"`
	Genres      []string        `json:"genres,omitempty"`
	Description string          `json:"description,omitempty"`
}

// UpdateVenueRequest represents the request body for updating a venue
type UpdateVenueRequest struct {
	Name        *string         `json:"name,omitempty"`
	Location    *LocationRequest `json:"location,omitempty"`
	Address     *AddressRequest  `json:"address,omitempty"`
	VenueTypes  []string        `json:"venue_types,omitempty"`
	Capacity    *int            `json:"capacity,omitempty"`
	Genres      []string        `json:"genres,omitempty"`
	Description *string         `json:"description,omitempty"`
}

// Search handles venue search requests
// GET /api/v1/venues/search
func (h *VenueHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	criteria := &domain.VenueSearchCriteria{
		Limit:  10, // Default limit
		Offset: 0,
	}

	// Parse location parameters
	if latStr := query.Get("lat"); latStr != "" {
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			http.Error(w, "invalid latitude", http.StatusBadRequest)
			return
		}

		lngStr := query.Get("lng")
		lng, err := strconv.ParseFloat(lngStr, 64)
		if err != nil {
			http.Error(w, "invalid longitude", http.StatusBadRequest)
			return
		}

		radiusStr := query.Get("radius")
		radius, err := strconv.ParseFloat(radiusStr, 64)
		if err != nil {
			http.Error(w, "invalid radius", http.StatusBadRequest)
			return
		}

		criteria.Location = &domain.GeoPoint{
			Latitude:  lat,
			Longitude: lng,
		}
		criteria.RadiusKm = radius
	}

	// Parse city/state
	if city := query.Get("city"); city != "" {
		criteria.City = city
	}
	if state := query.Get("state"); state != "" {
		criteria.State = state
	}
	if country := query.Get("country"); country != "" {
		criteria.Country = country
	}

	// Parse venue types
	if venueTypesStr := query.Get("venue_types"); venueTypesStr != "" {
		types := strings.Split(venueTypesStr, ",")
		criteria.VenueTypes = make([]domain.VenueType, len(types))
		for i, t := range types {
			criteria.VenueTypes[i] = domain.VenueType(strings.TrimSpace(t))
		}
	}

	// Parse capacity filters
	if minCapStr := query.Get("min_capacity"); minCapStr != "" {
		minCap, err := strconv.Atoi(minCapStr)
		if err == nil {
			criteria.MinCapacity = minCap
		}
	}
	if maxCapStr := query.Get("max_capacity"); maxCapStr != "" {
		maxCap, err := strconv.Atoi(maxCapStr)
		if err == nil {
			criteria.MaxCapacity = maxCap
		}
	}

	// Parse genres
	if genresStr := query.Get("genres"); genresStr != "" {
		criteria.Genres = strings.Split(genresStr, ",")
		for i := range criteria.Genres {
			criteria.Genres[i] = strings.TrimSpace(criteria.Genres[i])
		}
	}

	// Parse payment filters
	if minPayStr := query.Get("min_pay"); minPayStr != "" {
		minPay, err := strconv.Atoi(minPayStr)
		if err == nil {
			criteria.MinPay = minPay
		}
	}
	if maxPayStr := query.Get("max_pay"); maxPayStr != "" {
		maxPay, err := strconv.Atoi(maxPayStr)
		if err == nil {
			criteria.MaxPay = maxPay
		}
	}

	// Parse rating filter
	if minRatingStr := query.Get("min_rating"); minRatingStr != "" {
		minRating, err := strconv.ParseFloat(minRatingStr, 64)
		if err == nil {
			criteria.MinRating = minRating
		}
	}

	// Parse boolean filters
	if verifiedStr := query.Get("verified_only"); verifiedStr == "true" {
		criteria.VerifiedOnly = true
	}
	if activeStr := query.Get("active_only"); activeStr == "true" {
		criteria.ActiveOnly = true
	}

	// Parse pagination
	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			criteria.Limit = limit
		}
	}
	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err == nil && offset >= 0 {
			criteria.Offset = offset
		}
	}

	// Parse sorting
	if sortBy := query.Get("sort_by"); sortBy != "" {
		criteria.SortBy = domain.VenueSortField(sortBy)
	}
	if sortOrder := query.Get("sort_order"); sortOrder != "" {
		criteria.SortOrder = domain.SortOrder(sortOrder)
	}

	// Execute search
	result, err := h.service.Search(r.Context(), criteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetByID retrieves a venue by ID
// GET /api/v1/venues/{id}
func (h *VenueHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	venue, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "venue not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(venue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Create creates a new venue
// POST /api/v1/venues
func (h *VenueHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateVenueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if req.Location.Latitude == 0 || req.Location.Longitude == 0 {
		http.Error(w, "location is required", http.StatusBadRequest)
		return
	}
	if len(req.VenueTypes) == 0 {
		http.Error(w, "at least one venue type is required", http.StatusBadRequest)
		return
	}

	// Convert venue types
	venueTypes := make([]domain.VenueType, len(req.VenueTypes))
	for i, t := range req.VenueTypes {
		venueTypes[i] = domain.VenueType(t)
	}

	// Create venue
	venue := domain.NewVenue(
		req.Name,
		domain.GeoPoint{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		},
		domain.Address{
			Street:     req.Address.Street,
			City:       req.Address.City,
			State:      req.Address.State,
			PostalCode: req.Address.PostalCode,
			Country:    req.Address.Country,
		},
		venueTypes,
		domain.SourceUserSubmitted,
	)

	venue.Capacity = req.Capacity
	venue.Genres = req.Genres
	venue.Description = req.Description

	if err := h.service.Create(r.Context(), venue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(venue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Update updates an existing venue
// PUT /api/v1/venues/{id}
func (h *VenueHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	venue, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "venue not found", http.StatusNotFound)
		return
	}

	var req UpdateVenueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Apply updates
	if req.Name != nil {
		venue.Name = *req.Name
	}
	if req.Location != nil {
		venue.Location.Latitude = req.Location.Latitude
		venue.Location.Longitude = req.Location.Longitude
		// Geohash will be regenerated by service
		venue.Location.Geohash = ""
	}
	if req.Address != nil {
		venue.Address.Street = req.Address.Street
		venue.Address.City = req.Address.City
		venue.Address.State = req.Address.State
		venue.Address.PostalCode = req.Address.PostalCode
		venue.Address.Country = req.Address.Country
	}
	if len(req.VenueTypes) > 0 {
		venueTypes := make([]domain.VenueType, len(req.VenueTypes))
		for i, t := range req.VenueTypes {
			venueTypes[i] = domain.VenueType(t)
		}
		venue.VenueTypes = venueTypes
	}
	if req.Capacity != nil {
		venue.Capacity = *req.Capacity
	}
	if len(req.Genres) > 0 {
		venue.Genres = req.Genres
	}
	if req.Description != nil {
		venue.Description = *req.Description
	}

	if err := h.service.Update(r.Context(), venue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(venue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Delete deletes a venue
// DELETE /api/v1/venues/{id}
func (h *VenueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
