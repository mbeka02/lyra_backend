package model

import "mime/multipart"

type GetSignedURLRequest struct {
	UnsignedURL string `json:"unsigned_url" validate:"required"`
}

// CreateDocumentReferenceRequest defines the input for creating a DocumentReference.
// The actual file is passed separately.
type CreateDocumentReferenceRequest struct {
	PatientID      int64   `json:"patient_id"`       // To build the Subject reference
	SpecialistID   *int64  `json:"specialist_id"`    // Optional: Uploader if a specialist
	Title          *string `json:"title"`            // Optional: User-provided title for the document
	DocTypeCode    *string `json:"doc_type_code"`    // Optional: Code for document type (e.g., LOINC)
	DocTypeDisplay *string `json:"doc_type_display"` // Optional: Display name for document type
}

// CreateDocumentReferenceServiceInput combines the request metadata and the file.
// This internal struct helps pass data within the service call.
type CreateDocumentReferenceServiceInput struct {
	Metadata   CreateDocumentReferenceRequest
	FileHeader *multipart.FileHeader // The actual file to upload
}
