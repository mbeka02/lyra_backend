package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/service"
)

type PatientHandler struct {
	patientService service.PatientService
}

func NewPatientHandler(patientService service.PatientService) *PatientHandler {
	return &PatientHandler{
		patientService,
	}
}

func (h *PatientHandler) HandleGetPatient(w http.ResponseWriter, r *http.Request) {
	patientParam := chi.URLParam(r, "patientId")
	patientID, err := strconv.ParseInt(patientParam, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf(" error invalid patient id:%w", err))
		return
	}
	patientDetails, err := h.patientService.GetPatientAccountDetails(r.Context(), patientID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to get patient details"))
		return
	}
	respondWithJSON(w, http.StatusOK, patientDetails)
}

func (h *PatientHandler) HandleCreatePatient(w http.ResponseWriter, r *http.Request) {
	request := model.CreatePatientRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}

	patient, err := h.patientService.CreatePatient(r.Context(), request, payload.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusCreated, patient)
}
