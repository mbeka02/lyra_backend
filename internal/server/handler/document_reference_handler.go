package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/middleware"
	"github.com/mbeka02/lyra_backend/internal/server/service"
)

type DocumentReferenceHandler struct {
	patientService  service.PatientService
	doctorService   service.DoctorService
	documentService service.DocumentReferenceService
}

func NewDocumentReferenceHandler(patientService service.PatientService, doctorService service.DoctorService, documentService service.DocumentReferenceService) *DocumentReferenceHandler {
	return &DocumentReferenceHandler{patientService, doctorService, documentService}
}

func (h *DocumentReferenceHandler) HandleCreateDocumentReference(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, err := middleware.GetAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// Get the file
	_, fileHeader, err := r.FormFile("document") // Key for the file part
	if err != nil {
		if err == http.ErrMissingFile {
			respondWithError(w, http.StatusBadRequest, fmt.Errorf("missing required file field 'document'"))
		} else {
			respondWithError(w, http.StatusBadRequest, fmt.Errorf("error retrieving file 'document': %w", err))
		}
		return
	}

	// get the metadata JSON string
	metadataJSON := r.FormValue("metadata") // Key for the JSON metadata part
	if metadataJSON == "" {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("missing required metadata field"))
		return
	}
	// validate the metadata
	var reqMetadata model.CreateDocumentReferenceRequest

	if err := json.Unmarshal([]byte(metadataJSON), &reqMetadata); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("error parsing metadata JSON: %w", err))
		return
	}

	// determine Actual PatientID and SpecialistID using Services
	var targetPatientID int64
	var uploaderSpecialistID *int64 // Pointer because it's optional

	if payload.Role == "patient" {
		// Patient uploads for themselves. Get their PatientID.
		pID, err := h.patientService.GetPatientIdByUserId(r.Context(), payload.UserID) // Assume service method exists
		if err != nil {
			// Consider specific errors, e.g., sql.ErrNoRows -> Patient not found
			respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to find patient record for user: %w", err))
			return
		}
		targetPatientID = pID
		// SpecialistID remains nil

	} else if payload.Role == "specialist" {
		// Specialist uploads FOR a patient.
		// Target PatientID *must* be in the metadata provided by the specialist.
		if reqMetadata.PatientID == 0 { // Assuming 0 indicates not provided in JSON
			respondWithError(w, http.StatusBadRequest, fmt.Errorf("metadata must include target patientId when uploaded by a specialist"))
			return
		}
		targetPatientID = reqMetadata.PatientID // Trust the ID provided by the specialist (add validation/authz later)

		// get the SpecialistID of the uploader (for Author field).
		spID, err := h.doctorService.GetDoctorIdByUserId(r.Context(), payload.UserID) // Assume service method exists
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to find specialist record for user: %w", err))
			return
		}
		tempSpID := spID                 // Create temp var to take address
		uploaderSpecialistID = &tempSpID // Assign the pointer

	} else {
		respondWithError(w, http.StatusForbidden, fmt.Errorf("user role not permitted to upload documents"))
		return
	}

	// Authorization Check
	// E.g., If specialist, check if they are allowed to access/modify targetPatientID's records.

	// prepare input for the DocumentReferenceService
	// update the metadata struct with the *verified/derived* IDs
	reqMetadata.PatientID = targetPatientID
	reqMetadata.SpecialistID = uploaderSpecialistID

	serviceInput := model.CreateDocumentReferenceServiceInput{
		Metadata:   reqMetadata,
		FileHeader: fileHeader,
	}

	// call the DocumentReference Service
	savedFhirDocRef, err := h.documentService.CreateDocumentReference(r.Context(), serviceInput)
	if err != nil {
		// log the internal error for debugging
		fmt.Printf("ERROR: CreateDocumentReference failed: %v\n", err)
		// respond with a generic server error
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to create document reference"))
		return
	}

	// respond with success
	if err := respondWithJSON(w, http.StatusCreated, savedFhirDocRef); err != nil {
		// log this error, as headers might already be sent
		fmt.Printf("ERROR: Failed to write JSON response: %v\n", err)
	}
}

// handleListPatientDocuments handles GET requests for a patient's documents.
func (h *DocumentReferenceHandler) HandleListPatientDocuments(w http.ResponseWriter, r *http.Request) {
	// ensure auth payload is present
	payload, err := middleware.GetAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	// Extract Target Patient ID
	var targetPatientID int64
	// params
	params := NewQueryParamExtractor(r)
	patientIdStr := params.GetString("patientId")

	// TODO: Authorization Check (CRUCIAL!)
	// Can the authenticated user (payload.UserID, payload.Role) view documents for targetPatientID?
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
		// This requires logic in the specialistService/Repo to check the relationship.
		if patientIdStr == "" {
			respondWithError(w, http.StatusBadRequest, fmt.Errorf("missing required patientId parameter"))
			return
		}
		targetPatientID, err = strconv.ParseInt(patientIdStr, 10, 64)
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
		respondWithError(w, http.StatusForbidden, fmt.Errorf("you are not authorized to view documents for this patient"))
		return
	}

	// Extract Pagination Parameters
	// Standard FHIR param
	count := params.GetInt("_count", 0) // Ignore error, default to 0 (service might have its own default)
	if count <= 0 || count > 100 {      // Apply a reasonable max limit
		count = 20 // Default page size
	}

	pageToken := params.GetString("_page_token") // Or other token param if used

	// Call the Service
	bundle, err := h.documentService.ListPatientDocuments(r.Context(), targetPatientID, count, pageToken)
	if err != nil {
		log.Printf("ERROR: ListPatientDocuments failed: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to retrieve documents"))
		return
	}

	// Respond with the FHIR Bundle
	// Ensure response helpers set Content-Type: application/fhir+json
	w.Header().Set("Content-Type", "application/fhir+json") // Explicitly set here or in helper
	if err := respondWithJSON(w, http.StatusOK, bundle); err != nil {
		fmt.Printf("ERROR: Failed to write JSON response: %v\n", err)
	}
}
