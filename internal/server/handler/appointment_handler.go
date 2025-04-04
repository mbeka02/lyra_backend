package handler

import (
	"net/http"

	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/service"
)

type AppointmentHandler struct {
	appointmentService service.AppointmentService
}

func NewAppointmentHandler(appointmentService service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService,
	}
}

func (h *AppointmentHandler) HandleGetPatientAppointments(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	params := NewQueryParamExtractor(r)

	response, err := h.appointmentService.GetPatientAppointments(r.Context(), service.GetPatientAppointmentsParams{
		UserId:   payload.UserID,
		Status:   params.GetString("status"),
		Interval: params.GetInt32("interval"),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	if err := respondWithJSON(w, http.StatusOK, response); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *AppointmentHandler) HandleCreateAppointment(w http.ResponseWriter, r *http.Request) {
	request := model.CreateAppointmentRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	response, err := h.appointmentService.CreateAppointmentWithPayment(r.Context(), request, payload.UserID, payload.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	if err := respondWithJSON(w, http.StatusCreated, response); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}
