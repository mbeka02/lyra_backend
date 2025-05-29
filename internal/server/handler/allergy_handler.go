package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv" // For patientId from path

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/service"
)

type AllergyHandler struct {
	allergyService service.AllergyService
	patientService service.PatientService
}

func NewAllergyHandler(as service.AllergyService, ps service.PatientService) *AllergyHandler {
	return &AllergyHandler{allergyService: as, patientService: ps}
}

func (h *AllergyHandler) HandleCreateAllergy(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}

	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in path"))
		return
	}

	var req model.CreateAllergyIntoleranceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}
	defer r.Body.Close()

	// TODO: More robust validation of req fields

	allergy, err := h.allergyService.CreateAllergy(r.Context(), req, payload.UserID, targetPatientID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err) // Or more specific errors
		return
	}
	respondWithJSON(w, http.StatusCreated, allergy)
}

func (h *AllergyHandler) HandleListAllergies(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in path"))
		return
	}

	allergies, err := h.allergyService.ListAllergiesForPatient(r.Context(), payload.UserID, targetPatientID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, allergies)
}

func (h *AllergyHandler) HandleGetAllergy(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in path"))
		return
	}

	allergyIdStr := chi.URLParam(r, "allergyId")
	allergyID, err := uuid.Parse(allergyIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid allergyId in path"))
		return
	}

	allergy, err := h.allergyService.GetAllergy(r.Context(), allergyID, payload.UserID, targetPatientID)
	if err != nil {
		if err == sql.ErrNoRows { // Or your custom "not found" error from service
			respondWithError(w, http.StatusNotFound, fmt.Errorf("allergy not found"))
		} else {
			respondWithError(w, http.StatusInternalServerError, err)
		}
		return
	}
	respondWithJSON(w, http.StatusOK, allergy)
}

func (h *AllergyHandler) HandleUpdateAllergy(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in path"))
		return
	}

	allergyIdStr := chi.URLParam(r, "allergyId")
	allergyID, err := uuid.Parse(allergyIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid allergyId in path"))
		return
	}

	var req model.UpdateAllergyIntoleranceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}
	defer r.Body.Close()

	allergy, err := h.allergyService.UpdateAllergy(r.Context(), allergyID, req, payload.UserID, targetPatientID)
	if err != nil {
		if err == sql.ErrNoRows { // Or custom error
			respondWithError(w, http.StatusNotFound, fmt.Errorf("allergy not found or not authorized to update"))
		} else {
			respondWithError(w, http.StatusInternalServerError, err)
		}
		return
	}
	respondWithJSON(w, http.StatusOK, allergy)
}

func (h *AllergyHandler) HandleDeleteAllergy(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, ok := getAuthPayload(w, r)
	if !ok {
		return
	}
	patientIdStr := chi.URLParam(r, "patientId")
	targetPatientID, err := strconv.ParseInt(patientIdStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid patientId in path"))
		return
	}

	allergyIdStr := chi.URLParam(r, "allergyId")
	allergyID, err := uuid.Parse(allergyIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid allergyId in path"))
		return
	}

	err = h.allergyService.DeleteAllergy(r.Context(), allergyID, payload.UserID, targetPatientID)
	if err != nil {
		if err == sql.ErrNoRows { // Or custom error
			respondWithError(w, http.StatusNotFound, fmt.Errorf("allergy not found or not authorized to delete"))
		} else {
			respondWithError(w, http.StatusInternalServerError, err)
		}
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Allergy record deleted successfully"})
}
