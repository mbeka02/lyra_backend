package service

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	samplyFhir "github.com/samply/golang-fhir-models/fhir-models/fhir"

	"github.com/mbeka02/lyra_backend/internal/fhir"
	"github.com/mbeka02/lyra_backend/internal/model"
)

type ObservationService interface {
	CreateObservation(ctx context.Context, req model.CreateObservationRequest) (*samplyFhir.Observation, error)
	CreateConsultationNote(ctx context.Context, req model.CreateConsultationNoteRequest, specialistDomainID int64) (*samplyFhir.Observation, error)
	SearchObservations(ctx context.Context, targetPatientID int64, categoryCode string, code string, count int) (*samplyFhir.Bundle, error)
}

type observationService struct {
	fhirClient *fhir.FHIRClient
}

func NewObservationService(fhirClient *fhir.FHIRClient) ObservationService {
	return &observationService{
		fhirClient: fhirClient,
	}
}

// CreateConsultationNote handles the business logic for creating a consultation note.
// specialistDomainID is the ID of the specialist (e.g., from your 'doctors' table),
// which will be used as the FHIR Practitioner ID.
func (s *observationService) CreateConsultationNote(
	ctx context.Context,
	req model.CreateConsultationNoteRequest,
	specialistDomainID int64,
) (*samplyFhir.Observation, error) {
	effectiveTime := time.Now() // Or derive from appointment context

	// Build the FHIR Observation resource for the note.
	fhirObs, err := fhir.BuildFHIRObservationFromNote(
		req.PatientID,      // This is the target patient's *domain* ID
		specialistDomainID, // Specialist's *domain* ID (used as Practitioner ID)
		req.NoteText,
		effectiveTime,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build FHIR Observation for note: %w", err)
	}

	// Call the FHIR client to create the Observation in the FHIR Store.
	savedFhirObs, err := s.fhirClient.CreateObservation(ctx, fhirObs) // Assumes fhirClient.CreateObservation exists
	if err != nil {
		// Log the internal error for debugging
		fmt.Printf("ERROR: CreateObservation (note) failed in FHIR client: %v\n", err)
		// Return a more generic error to the handler/user
		return nil, fmt.Errorf("failed to save consultation note")
	}

	// Return the saved resource (which includes server-assigned ID, meta, etc.)
	return savedFhirObs, nil
}

// SearchObservations fetches observations based on query parameters.
func (s *observationService) SearchObservations(
	ctx context.Context,
	targetPatientID int64,
	categoryCode string, // e.g., "notes"
	code string, // e.g., "urn:lyra:codesystem:observation-type|CONSULTATION_NOTE"
	count int,
) (*samplyFhir.Bundle, error) {
	queryValues := url.Values{}
	queryValues.Set("subject", fmt.Sprintf("Patient/%d", targetPatientID))
	queryValues.Set("_sort", "-effective-date") // Sort by effective date, most recent first

	if categoryCode != "" {
		queryValues.Set("category", categoryCode)
	}
	if code != "" {
		queryValues.Set("code", code) // Search by specific code if provided
	}

	if count > 0 {
		queryValues.Set("_count", strconv.Itoa(count))
	} else {
		queryValues.Set("_count", "20") // Default count
	}

	queryString := queryValues.Encode()

	// Call FHIR Client Search (ensure SearchObservations exists in FHIRClient)
	bundle, err := s.fhirClient.SearchForResorces(ctx, "Observation", queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to search observations in FHIR store: %w", err)
	}
	return bundle, nil
}

// Ensure your FHIRClient has this method:
// // internal/fhir/client.go - Snippet
//
//	func (f *FHIRClient) SearchObservations(ctx context.Context, queryParams string) (*samplyFhir.Bundle, error) {
//		resourceType := "Observation"
//		// ... rest is identical to SearchDocumentReferences, just change resourceType ...
//		// Ensure _type=Observation or resourceType in path for Execute GET is correct.
//	}
func (s *observationService) CreateObservation(ctx context.Context, req model.CreateObservationRequest) (*samplyFhir.Observation, error) {
	// Build FHIR Observation
	fhirObs, err := fhir.BuildFHIRObservation(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build FHIR Observation: %w", err)
	}

	// Create in FHIR Store using the dedicated client method
	savedFhirObs, err := s.fhirClient.CreateObservation(ctx, fhirObs)
	if err != nil {
		return nil, fmt.Errorf("failed to create Observation in FHIR store: %w", err)
	}

	// Return saved resource
	return savedFhirObs, nil
}
