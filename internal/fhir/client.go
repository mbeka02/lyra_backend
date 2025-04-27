package fhir

import (
	"context"
	"fmt"

	samplyFhir "github.com/samply/golang-fhir-models/fhir-models/fhir"
	healthcare "google.golang.org/api/healthcare/v1"
	"google.golang.org/api/option"
)

// FHIRClient wraps Google Cloud Healthcare FHIR API calls.
// It handles JSON marshalling/unmarshalling of FHIR resources.
type FHIRClient struct {
	svc      *healthcare.Service
	basePath string
}

// FHIRConfig holds the project id ,dataset location , id and the id of the fhir store
type FHIRConfig struct {
	ProjectID       string
	DatasetLocation string
	DatasetID       string
	FHIRStoreID     string
}

// NewFHIRClient initializes the Healthcare API client using ADC.
func NewFHIRClient(ctx context.Context, config FHIRConfig, opts ...option.ClientOption) (*FHIRClient, error) {
	// Default to CloudPlatformScope if not overridden
	opts = append([]option.ClientOption{option.WithScopes(healthcare.CloudPlatformScope)}, opts...)
	svc, err := healthcare.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("healthcare.NewService error: %w", err)
	}
	base := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir",
		config.ProjectID, config.DatasetLocation, config.DatasetID, config.FHIRStoreID)
	return &FHIRClient{svc: svc, basePath: base}, nil
}

func (f *FHIRClient) UpsertPatient(ctx context.Context, patient *samplyFhir.Patient) {
}
