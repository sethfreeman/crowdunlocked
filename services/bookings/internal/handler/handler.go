package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/crowdunlocked/services/bookings/internal/domain"
	"github.com/crowdunlocked/services/bookings/internal/repository"
	"github.com/go-chi/chi/v5"
)

type BookingHandler struct {
	repo *repository.BookingRepository
}

func NewBookingHandler(repo *repository.BookingRepository) *BookingHandler {
	return &BookingHandler{repo: repo}
}

type CreateBookingRequest struct {
	ArtistID  string    `json:"artist_id"`
	VenueID   string    `json:"venue_id"`
	EventDate time.Time `json:"event_date"`
	Fee       float64   `json:"fee"`
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	booking := domain.NewBooking(req.ArtistID, req.VenueID, req.EventDate, req.Fee)
	if err := h.repo.Create(r.Context(), booking); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	booking, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if booking == nil {
		http.Error(w, "booking not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *BookingHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	booking, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if booking == nil {
		http.Error(w, "booking not found", http.StatusNotFound)
		return
	}

	booking.Confirm()
	if err := h.repo.Update(r.Context(), booking); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
