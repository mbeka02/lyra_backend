package service

import (
	"context"
	"fmt"

	samplyFhir "github.com/samply/golang-fhir-models/fhir-models/fhir"

	"github.com/mbeka02/lyra_backend/internal/fhir"
	"github.com/mbeka02/lyra_backend/internal/model"
)

type ObservationService interface {
	CreateObservation(ctx context.Context, req model.CreateObservationRequest) (*samplyFhir.Observation, error)
}

type observationService struct {
	fhirClient *fhir.FHIRClient
}

func NewObservationService(fhirClient *fhir.FHIRClient) ObservationService {
	return &observationService{
		fhirClient: fhirClient,
	}
}

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
