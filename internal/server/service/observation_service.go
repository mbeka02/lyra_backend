package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
	samplyFhir "github.com/samply/golang-fhir-models/fhir-models/fhir"

	"github.com/mbeka02/lyra_backend/internal/fhir"
	"github.com/mbeka02/lyra_backend/internal/model"
)

type ObservationService interface {
	CreateObservation(ctx context.Context, req model.CreateObservationRequest) (*samplyFhir.Observation, error)
	CreateConsultationNote(ctx context.Context, req model.CreateConsultationNoteRequest, specialistDomainID int64) (*samplyFhir.Observation, error)
	SearchObservations(ctx context.Context, targetPatientID int64, categoryCode string, code string, count int) (*samplyFhir.Bundle, error)
	CreateObservationInDB(ctx context.Context, req model.CreateObservationRequestForDB, actingUserID int64, forPatientID int64 /*, specialistID *int64 if tracking*/) (database.Observation, error)
	GetObservation(ctx context.Context, observationID uuid.UUID, actingUserID int64, forPatientID int64) (database.Observation, error)
	ListObservationsForPatient(ctx context.Context, actingUserID int64, forPatientID int64) ([]database.Observation, error)
	UpdateObservation(ctx context.Context, observationID uuid.UUID, req model.UpdateObservationRequest, actingUserID int64, forPatientID int64 /*, specialistID *int64 if tracking*/) (database.Observation, error)
	DeleteObservation(ctx context.Context, observationID uuid.UUID, actingUserID int64, forPatientID int64) error
}

type observationService struct {
	observationRepo repository.ObservationRepository
	fhirClient      *fhir.FHIRClient
}

func NewObservationService(obsRepo repository.ObservationRepository, fhirClient *fhir.FHIRClient) ObservationService {
	return &observationService{
		fhirClient:      fhirClient,
		observationRepo: obsRepo,
	}
}

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

func (s *observationService) CreateObservationInDB(
	ctx context.Context,
	req model.CreateObservationRequestForDB,
	actingUserID int64,
	forPatientID int64,
	// specialistIDIfCreating *int64, // Pass specialistID explicitly if relevant
) (database.Observation, error) {
	// TODO: Authorization: Can actingUserID create observation for forPatientID?
	// E.g., is actingUser the patient themselves (if allowed), or an authorized specialist?

	params := database.CreateObservationParams{
		PatientID:         forPatientID,
		Status:            req.Status,
		CodeText:          req.CodeText,
		EffectiveDateTime: req.EffectiveDateTime, // Should be TIMESTAMPTZ
		ValueString:       req.ValueString,
		// SpecialistID:        ToNullInt64(specialistIDIfCreating), // If tracking specialist
	}
	return s.observationRepo.Create(ctx, params)
}

func (s *observationService) GetObservation(
	ctx context.Context,
	observationID uuid.UUID,
	actingUserID int64,
	forPatientID int64,
) (database.Observation, error) {
	// TODO: Authorization: Can actingUserID view observation for forPatientID?
	return s.observationRepo.GetByID(ctx, observationID, forPatientID)
}

func (s *observationService) ListObservationsForPatient(
	ctx context.Context,
	actingUserID int64,
	forPatientID int64,
) ([]database.Observation, error) {
	// TODO: Authorization: Can actingUserID list observations for forPatientID?
	return s.observationRepo.ListByPatientID(ctx, forPatientID)
}

func (s *observationService) UpdateObservation(
	ctx context.Context,
	observationID uuid.UUID,
	req model.UpdateObservationRequest,
	actingUserID int64,
	forPatientID int64,
	// specialistIDIfUpdating *int64,
) (database.Observation, error) {
	// TODO: Authorization: Can actingUserID update observation for forPatientID?

	// Optional: Fetch existing to ensure it belongs to forPatientID before update
	// _, err := s.observationRepo.GetByID(ctx, observationID, forPatientID)
	// if err != nil {
	// 	return database.Observation{}, fmt.Errorf("observation not found or access denied: %w", err)
	// }

	params := database.UpdateObservationParams{
		ID:                observationID,
		PatientID:         forPatientID, // Used in WHERE clause for safety
		Status:            req.Status,
		CodeText:          req.CodeText,
		EffectiveDateTime: req.EffectiveDateTime,
		ValueString:       req.ValueString,
		// SpecialistID:        ToNullInt64(specialistIDIfUpdating), // If tracking specialist
		// UpdatedAt is handled by default in the query or trigger
	}
	return s.observationRepo.Update(ctx, params)
}

func (s *observationService) DeleteObservation(
	ctx context.Context,
	observationID uuid.UUID,
	actingUserID int64,
	forPatientID int64,
) error {
	// TODO: Authorization: Can actingUserID delete observation for forPatientID?

	// Optional: Fetch existing to ensure it belongs to forPatientID before delete
	// _, err := s.observationRepo.GetByID(ctx, observationID, forPatientID)
	// if err != nil {
	// 	return fmt.Errorf("observation not found or access denied: %w", err)
	// }
	return s.observationRepo.Delete(ctx, observationID, forPatientID)
}

// Helper for nullable int64 (if using specialist_id)
func ToNullInt64(val *int64) sql.NullInt64 {
	if val == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *val, Valid: true}
}
