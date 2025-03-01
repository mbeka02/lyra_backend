package handler

import (
	"net/http"

	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/middleware"
	"github.com/mbeka02/lyra_backend/internal/server/service"
)

type AvailabilityHandler struct {
	availabilityService service.AvailabilityService
}

func NewAvailabilityHandler(availabilityService service.AvailabilityService) *AvailabilityHandler {
	return &AvailabilityHandler{
		availabilityService,
	}
}

func (h *AvailabilityHandler) HandleCreateAvailability(w http.ResponseWriter, r *http.Request) {
	request := model.CreateAvailabilityRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	// ensure auth payload is present
	payload, err := middleware.GetAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	response, err := h.availabilityService.CreateAvailability(r.Context(), request, payload.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	if err := respondWithJSON(w, http.StatusCreated, response); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *AvailabilityHandler) HandleGetAvailabilityByDoctor(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, err := middleware.GetAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	response, err := h.availabilityService.GetAvailabilityByDoctor(r.Context(), payload.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	if err := respondWithJSON(w, http.StatusCreated, response); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}
