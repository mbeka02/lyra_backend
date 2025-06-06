package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mbeka02/lyra_backend/internal/model"
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

func (h *AvailabilityHandler) HandleGetSlots(w http.ResponseWriter, r *http.Request) {
	request := model.GetSlotsRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	slots, err := h.availabilityService.GetSlots(r.Context(), request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusCreated, slots)
}

func (h *AvailabilityHandler) HandleCreateAvailability(w http.ResponseWriter, r *http.Request) {
	var request model.CreateAvailabilityRequest
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	availabilitySlot, err := h.availabilityService.CreateAvailability(r.Context(), request, payload.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusCreated, availabilitySlot)
}

func (h *AvailabilityHandler) HandleDeleteById(w http.ResponseWriter, r *http.Request) {
	availabilityParam := chi.URLParam(r, "availabilityId")
	availabilityId, err := strconv.ParseInt(availabilityParam, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf(" error invalid availability id:%w", err))

		return
	}
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	err = h.availabilityService.DeleteById(r.Context(), availabilityId, payload.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, "removed the slot from your schedule")
}

func (h *AvailabilityHandler) HandleDeleteByDay(w http.ResponseWriter, r *http.Request) {
	weekParam := chi.URLParam(r, "dayOfWeek")
	dayOfWeek, err := strconv.ParseInt(weekParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	err = h.availabilityService.DeleteByDay(r.Context(), int32(dayOfWeek), payload.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, "removed the slots from your schedule")
}

func (h *AvailabilityHandler) HandleGetAvailabilityByDoctor(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	response, err := h.availabilityService.GetAvailabilityByDoctor(r.Context(), payload.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, response)
}
