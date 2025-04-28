package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/middleware"
	"github.com/mbeka02/lyra_backend/internal/server/service"
)

// Assume maxDocumentSize is defined elsewhere or define it here
const maxDocumentSize = 50 * 1024 * 1024 // 50 MB
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

	// Optional but Recommended: Authorization Check
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
