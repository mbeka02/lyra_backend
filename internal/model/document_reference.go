package model

import "mime/multipart"

// CreateDocumentReferenceRequest defines the input for creating a DocumentReference.
// The actual file is passed separately.
type CreateDocumentReferenceRequest struct {
	PatientID      int64   `json:"patientId" validate:"required"` // To build the Subject reference
	SpecialistID   *int64  `json:"specialistId"`                  // Optional: Uploader if a specialist
	Title          *string `json:"title"`                         // Optional: User-provided title for the document
	DocTypeCode    *string `json:"docTypeCode"`                   // Optional: Code for document type (e.g., LOINC)
	DocTypeDisplay *string `json:"docTypeDisplay"`                // Optional: Display name for document type
}

// CreateDocumentReferenceServiceInput combines the request metadata and the file.
// This internal struct helps pass data within the service call.
type CreateDocumentReferenceServiceInput struct {
	Metadata   CreateDocumentReferenceRequest
	FileHeader *multipart.FileHeader // The actual file to upload
}
