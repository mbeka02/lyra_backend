package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath" // For getting file extension
	"strings"

	"github.com/google/uuid" // For unique object names
	samplyFhir "github.com/samply/golang-fhir-models/fhir-models/fhir"

	"github.com/mbeka02/lyra_backend/internal/fhir"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/objstore" // Need storage
)

const maxDocumentSize = 50 * 1024 * 1024    // 50 MB limit for documents
var allowedDocumentTypes = map[string]bool{ // Example list, expand significantly
	"application/pdf":    true,
	"image/jpeg":         true,
	"image/png":          true,
	"image/tiff":         true,
	"text/plain":         true,
	"application/msword": true, // .doc
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
	// Add more MIME types as required (e.g., HL7 CDA XML, DICOM etc. if applicable)
}

type DocumentReferenceService interface {
	CreateDocumentReference(ctx context.Context, input model.CreateDocumentReferenceServiceInput) (*samplyFhir.DocumentReference, error)
	// TODO: Add Get, List methods etc
}

type documentReferenceService struct {
	fhirClient  *fhir.FHIRClient
	fileStorage objstore.Storage // Inject storage dependency
}

// NewDocumentReferenceService creates a new DocumentReferenceService.
func NewDocumentReferenceService(fhirClient *fhir.FHIRClient, fileStorage objstore.Storage) DocumentReferenceService {
	return &documentReferenceService{
		fhirClient:  fhirClient,
		fileStorage: fileStorage,
	}
}

func (s *documentReferenceService) CreateDocumentReference(ctx context.Context, input model.CreateDocumentReferenceServiceInput) (*samplyFhir.DocumentReference, error) {
	fileHeader := input.FileHeader
	metadata := input.Metadata

	// Validate File Input
	err := s.validateDocumentFile(fileHeader)
	if err != nil {
		return nil, fmt.Errorf("invalid document file: %w", err)
	}

	// Generate unique GCS Object Name
	fileExt := filepath.Ext(fileHeader.Filename)
	objectName := fmt.Sprintf("patients/%d/documents/%s%s",
		metadata.PatientID,
		uuid.NewString(),
		fileExt,
	)

	// Upload file to GCS using injected storage service
	gcsUrl, err := s.fileStorage.Upload(ctx, objectName, fileHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to upload document to storage: %w", err)
	}

	// Extract necessary file metadata from header
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream" // Default if not provided
	}
	sizeBytes := fileHeader.Size
	/* NB: You easily get the *actual file creation time* from a multipart upload.
	If the client sends it (e.g., via model.AttachmentCreationStr), I could parse it here,for now, i'll pass nil for attachmentCreation.
	*/
	// Build the FHIR DocumentReference Resource
	fhirDocRef, err := fhir.BuildFHIRDocumentReference(
		metadata,
		gcsUrl,
		contentType,
		sizeBytes,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build FHIR DocumentReference: %w", err)
	}

	// Save resource to the FHIR Store
	savedFhirDocRef, err := s.fhirClient.CreateDocumentReference(ctx, fhirDocRef)
	if err != nil {
		/*TODO:Attempt to clean up GCS file if FHIR creation fails(Could be complex)
		Consider adding cleanup logic here or documenting that orphans might occur.*/
		fmt.Printf("Error saving DocRef to FHIR Store, GCS object '%s' might be orphaned: %v\n", objectName, err)
		return nil, fmt.Errorf("failed to save DocumentReference in FHIR store: %w", err)
	}

	// Return the saved resource
	return savedFhirDocRef, nil
}

// This checks size and allowed types
func (s *documentReferenceService) validateDocumentFile(fileHeader *multipart.FileHeader) error {
	if fileHeader == nil {
		return fmt.Errorf("file header is missing")
	}
	fileContentType := fileHeader.Header.Get("Content-Type")
	if fileHeader.Size > maxDocumentSize {
		return fmt.Errorf("the document size %d exceeds the limit of %d bytes", fileHeader.Size, maxDocumentSize)
	}
	if fileContentType == "" {
		fmt.Println("Warning: Uploaded document file has no Content-Type header.")
	} else if _, ok := allowedDocumentTypes[strings.ToLower(fileContentType)]; !ok && len(allowedDocumentTypes) > 0 {
		// Only check if allowedDocumentTypes is populated
		return fmt.Errorf("this file format is not supported: %s", fileContentType)
	}
	return nil
}
