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

func (h *AppointmentHandler) HandleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	var request model.UpdateAppointmentStatusRequest
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	err := h.appointmentService.UpdateAppointmentStatus(r.Context(), request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, "updated appointment status successfully")
}

func (h *AppointmentHandler) HandleGetCompletedAppointments(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	completedAppointments, err := h.appointmentService.GetAppointmentIDs(r.Context(), service.GetAppointmentIDsParams{
		Role:   payload.Role,
		UserID: payload.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, completedAppointments)
}

func (h *AppointmentHandler) HandleGetDoctorAppointments(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	params := NewQueryParamExtractor(r)
	defaultInterval := 21
	appointments, err := h.appointmentService.GetDoctorAppointments(r.Context(), service.GetAppointmentsParams{
		UserID:   payload.UserID,
		Status:   params.GetString("status"),
		Interval: params.GetInt32("interval", int32(defaultInterval)),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, appointments)
}

func (h *AppointmentHandler) HandleGetPatientAppointments(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	params := NewQueryParamExtractor(r)
	defaultInterval := 21
	appointments, err := h.appointmentService.GetPatientAppointments(r.Context(), service.GetAppointmentsParams{
		UserID:   payload.UserID,
		Status:   params.GetString("status"),
		Interval: params.GetInt32("interval", int32(defaultInterval)),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, appointments)
}

func (h *AppointmentHandler) HandleCreateAppointment(w http.ResponseWriter, r *http.Request) {
	var request model.CreateAppointmentRequest
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}

	appointment, err := h.appointmentService.CreateAppointmentWithPayment(r.Context(), request, payload.UserID, payload.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusCreated, appointment)
}
