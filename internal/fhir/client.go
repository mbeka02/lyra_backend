package fhir

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

func (f *FHIRClient) UpsertPatient(ctx context.Context, patient *samplyFhir.Patient) (*samplyFhir.Patient, error) {
	payload, err := json.Marshal(patient)
	if err != nil {
		return nil, fmt.Errorf("marshal Patient: %w", err)
	}
	// hb := &healthcare.HttpBody{Data: string(payload), ContentType: "application/fhir+json"}
	// var resp *healthcare.HttpBody
	path := f.basePath + "/Patient/" + *patient.Id
	var resp *http.Response
	// If the version Id is emtpy create otherwise update
	if patient.Meta == nil || patient.Meta.VersionId == nil || *patient.Meta.VersionId == "" {
		call := f.svc.Projects.Locations.Datasets.FhirStores.Fhir.Create(f.basePath+"/Patient", "Patient", bytes.NewReader(payload))
		call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
		resp, err = call.Do()
		if err != nil {
			return nil, fmt.Errorf("Create resource error: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode > 299 {
			return nil, fmt.Errorf("Create resource error: status %d : %s", resp.StatusCode, resp.Status)
		}
	} else {
		// Update with PUT; include If-Match header for optimistic concurrency
		call := f.svc.Projects.Locations.Datasets.FhirStores.Fhir.
			Update(path, bytes.NewReader(payload))
		call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
		call.Header().Set("If-Match", fmt.Sprintf(`W/"%s"`, *patient.Meta.VersionId))
		resp, err = call.Do()
	}
	var p samplyFhir.Patient
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, fmt.Errorf("error decoding patient response: %w", err)
	}
	return &p, nil
}
