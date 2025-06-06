package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/service"
)

type ObservationHandler struct {
	patientService     service.PatientService
	doctorService      service.DoctorService
	observationService service.ObservationService
}

func NewObservationHandler(patientService service.PatientService, doctorService service.DoctorService, observationService service.ObservationService) *ObservationHandler {
	return &ObservationHandler{patientService, doctorService, observationService}
}

// HandleCreateConsultationNote handles POST requests to create a new consultation note.
// This endpoint should be restricted to specialists.
func (h *ObservationHandler) HandleCreateConsultationNote(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	// Ensure the user is a specialist
	// TODO: Move these role checks to a middleware function
	if payload.Role != "specialist" {
		respondWithError(w, http.StatusForbidden, fmt.Errorf("access denied: only specialists can create notes"))
		return
	}

	doctorID, err := h.doctorService.GetDoctorIdByUserId(r.Context(), payload.UserID)
	if err != nil {
		fmt.Printf("Error getting specialist ID for user %d: %v\n", payload.UserID, err)
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("could not retrieve specialist details for user"))
		return
	}

	var request model.CreateConsultationNoteRequest
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)

		return
	}
	// TODO: Make a utility function for these authorization checks
	// Authorization Check: Is this doctor allowed to create notes for this patient?

	authorized := false
	isUnderCare, err := h.doctorService.IsPatientUnderCare(r.Context(), doctorID, request.PatientID)
	if err != nil {
		// Log the error from the care check
		log.Printf("Auth check error for doctor %d viewing patient %d: %v\n", doctorID, request.PatientID, err)
	} else {
		authorized = isUnderCare // Authorize if the service confirms care relationship
	}

	if !authorized {
		respondWithError(w, http.StatusForbidden, fmt.Errorf("you are not authorized to modify the details for this patient"))
		return
	}

	savedObservation, err := h.observationService.CreateConsultationNote(r.Context(), request, doctorID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/fhir+json") // Important for FHIR responses
	respondWithJSON(w, http.StatusCreated, savedObservation)
}

func (h *ObservationHandler) HandleListPatientObservations(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	// Extract Target Patient ID
	var targetPatientID int64
	// params
	params := NewQueryParamExtractor(r)
	patientIdStr := params.GetString("patientId")

	// Authorization Check
	// Can the authenticated user (payload.UserID, payload.Role) view observations for targetPatientID?
	authorized := false
	if payload.Role == "patient" {
		// Is the patient viewing their own documents?
		pID, err := h.patientService.GetPatientIdByUserId(r.Context(), payload.UserID)
		if err == nil /*&& pID == targetPatientID*/ {
			authorized = true
		}
		targetPatientID = pID
	} else if payload.Role == "specialist" {
		// Is the specialist viewing documents for a patient under their care?
		if patientIdStr == "" {
			respondWithError(w, http.StatusBadRequest, fmt.Errorf("missing required patientId parameter"))
			return
		}
		targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId parameter: %v", err))
			return
		}
		// First, get the DoctorID for the logged-in user
		doctorID, err := h.doctorService.GetDoctorIdByUserId(r.Context(), payload.UserID)
		if err != nil {
			log.Printf("Auth check failed for specialist user %d: %v\n", payload.UserID, err)
			respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to retrieve doctor information"))
			return
		}
		// Now, check the care relationship using the service method
		isUnderCare, err := h.doctorService.IsPatientUnderCare(r.Context(), doctorID, targetPatientID)
		if err != nil {
			// Log the error from the care check
			log.Printf("Auth check error for doctor %d viewing patient %d: %v\n", doctorID, targetPatientID, err)
		} else {
			authorized = isUnderCare // Authorize if the service confirms care relationship
		}
	}
	if !authorized {
		respondWithError(w, http.StatusForbidden, fmt.Errorf("not authorized to view observations for this patient"))
		return
	}

	// 4. Get optional query parameters for filtering observations
	categoryCode := params.GetString("category") // e.g., "notes"
	codeParam := params.GetString("code")        // e.g., "urn:lyra:codesystem:observation-type|CONSULTATION_NOTE"
	count := params.GetInt("_count", 0)
	if count <= 0 || count > 100 { // Apply a reasonable max limit
		count = 20 // Default page size
	}

	// Call the Observation Service
	bundle, err := h.observationService.SearchObservations(r.Context(), targetPatientID, categoryCode, codeParam, count)
	if err != nil {
		fmt.Printf("ERROR: ListObservations failed: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to retrieve observations"))
		return
	}

	// 6. Respond with the FHIR Bundle
	w.Header().Set("Content-Type", "application/fhir+json")
	respondWithJSON(w, http.StatusOK, bundle)
}

func (h *ObservationHandler) HandleCreateObservationInDB(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in path: %w", err))
		return
	}

	var req model.CreateObservationRequestForDB

	if err := parseAndValidateRequest(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	observation, err := h.observationService.CreateObservationInDB(r.Context(), req, payload.UserID, targetPatientID /*, specialistIDForDB */)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to create observation: %w", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, observation)
}

func (h *ObservationHandler) HandleListObservations(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}

	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in path: %w", err))
		return
	}

	observations, err := h.observationService.ListObservationsForPatient(r.Context(), payload.UserID, targetPatientID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to list observations: %w", err))
		return
	}
	respondWithJSON(w, http.StatusOK, observations)
}

func (h *ObservationHandler) HandleGetObservation(w http.ResponseWriter, r *http.Request) {
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in path: %w", err))
		return
	}

	observationIdStr := chi.URLParam(r, "observationId")
	observationID, err := uuid.Parse(observationIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid observationId in path: %w", err))
		return
	}

	observation, err := h.observationService.GetObservation(r.Context(), observationID, payload.UserID, targetPatientID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Errorf("observation not found"))
		} else {
			respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to get observation: %w", err))
		}
		return
	}
	respondWithJSON(w, http.StatusOK, observation)
}

func (h *ObservationHandler) HandleUpdateObservation(w http.ResponseWriter, r *http.Request) {
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in path: %w", err))
		return
	}

	observationIdStr := chi.URLParam(r, "observationId")
	observationID, err := uuid.Parse(observationIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid observationId in path: %w", err))
		return
	}

	var req model.UpdateObservationRequest
	if err := parseAndValidateRequest(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	updatedObservation, err := h.observationService.UpdateObservation(r.Context(), observationID, req, payload.UserID, targetPatientID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Errorf("observation not found or not authorized to update"))
		} else {
			respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to update observation: %w", err))
		}
		return
	}
	respondWithJSON(w, http.StatusOK, updatedObservation)
}

func (h *ObservationHandler) HandleDeleteObservation(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}

	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in path: %w", err))
		return
	}

	observationIdStr := chi.URLParam(r, "observationId")
	observationID, err := uuid.Parse(observationIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid observationId in path: %w", err))
		return
	}

	// TODO: Authorization check

	err = h.observationService.DeleteObservation(r.Context(), observationID, payload.UserID, targetPatientID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Errorf("observation not found or not authorized to delete"))
		} else {
			respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete observation: %w", err))
		}
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Observation record deleted successfully"})
}
