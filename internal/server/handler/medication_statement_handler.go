package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid" // For auth.Payload
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/service"
)

type MedicationHandler struct {
	medicationService service.MedicationService
}

func NewMedicationHandler(ms service.MedicationService) *MedicationHandler {
	return &MedicationHandler{
		medicationService: ms,
	}
}

func (h *MedicationHandler) HandleCreateMedication(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}

	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in URL path"))
		return
	}

	var req model.CreateMedicationStatementRequest
	if err := parseAndValidateRequest(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	// TODO: Add validation for req fields (e.g., status value from enum)

	// TODO: Authorization check: Can payload.UserID create a medication statement for targetPatientID?

	med, err := h.medicationService.CreateMedication(r.Context(), req, payload.UserID, targetPatientID)
	if err != nil {
		// Consider more specific error codes based on service error types
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to create medication statement: %w", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, med)
}

func (h *MedicationHandler) HandleListMedications(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in URL path"))
		return
	}

	// TODO: Authorization check: Can payload.UserID list medication statements for targetPatientID?

	meds, err := h.medicationService.ListMedicationsForPatient(r.Context(), payload.UserID, targetPatientID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to list medication statements: %w", err))
		return
	}
	respondWithJSON(w, http.StatusOK, meds)
}

func (h *MedicationHandler) HandleGetMedication(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in URL path"))
		return
	}

	medicationIdStr := chi.URLParam(r, "medicationId")
	medicationID, err := uuid.Parse(medicationIdStr) // Assuming UUID for IDs
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid medicationId in URL path: %w", err))
		return
	}

	// TODO: Authorization check: Can payload.UserID get this specific medication statement for targetPatientID?

	med, err := h.medicationService.GetMedication(r.Context(), medicationID, payload.UserID, targetPatientID)
	if err != nil {

		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to get medication statement: %w", err))
		return
	}
	respondWithJSON(w, http.StatusOK, med)
}

func (h *MedicationHandler) HandleUpdateMedication(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}

	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in URL path"))
		return
	}

	medicationIdStr := chi.URLParam(r, "medicationId")
	medicationID, err := uuid.Parse(medicationIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid medicationId in URL path: %w", err))
		return
	}

	var req model.UpdateMedicationStatementRequest
	if err := parseAndValidateRequest(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	// TODO: Authorization check

	updatedMed, err := h.medicationService.UpdateMedication(r.Context(), medicationID, req, payload.UserID, targetPatientID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to update medication statement: %w", err))
		return
	}
	respondWithJSON(w, http.StatusOK, updatedMed)
}

func (h *MedicationHandler) HandleDeleteMedication(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in URL path"))
		return
	}

	medicationIdStr := chi.URLParam(r, "medicationId")
	medicationID, err := uuid.Parse(medicationIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid medicationId in URL path: %w", err))
		return
	}

	// TODO: Authorization check

	err = h.medicationService.DeleteMedication(r.Context(), medicationID, payload.UserID, targetPatientID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete medication statement: %w", err))
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Medication statement deleted successfully"})
}
