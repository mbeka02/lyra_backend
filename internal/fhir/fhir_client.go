package fhir

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	base := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s",
		config.ProjectID, config.DatasetLocation, config.DatasetID, config.FHIRStoreID)
	return &FHIRClient{svc: svc, basePath: base}, nil
}

// This method creates or updates a document reference resource
func (f *FHIRClient) UpsertDocumentReference(ctx context.Context, docRef *samplyFhir.DocumentReference) (*samplyFhir.DocumentReference, error) {
	payload, err := json.Marshal(docRef)
	if err != nil {
		return nil, fmt.Errorf("marshal documentreference: %w", err)
	}
	resourceType := "DocumentReference"
	resourceID := *docRef.Id

	parentPath := fmt.Sprintf("%s", f.basePath)                                   // for Create
	resourcePath := fmt.Sprintf("%s/%s/%s", f.basePath, resourceType, resourceID) // for Update

	var resp *http.Response

	if docRef.Meta == nil || docRef.Meta.VersionId == nil || *docRef.Meta.VersionId == "" {
		// Create
		call := f.svc.Projects.Locations.Datasets.FhirStores.Fhir.Create(parentPath, "DocumentReference", bytes.NewReader(payload))
		call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
		resp, err = call.Do()
	} else {
		// Update
		call := f.svc.Projects.Locations.Datasets.FhirStores.Fhir.Update(resourcePath, bytes.NewReader(payload))
		call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
		call.Header().Set("If-Match", fmt.Sprintf(`W/"%s"`, *docRef.Meta.VersionId))
		resp, err = call.Do()
	}

	if err != nil {
		return nil, fmt.Errorf("upsert documentreference: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, f.readErrorResponse(resp)
	}

	var dr samplyFhir.DocumentReference
	if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
		return nil, fmt.Errorf("decode documentreference response: %w", err)
	}
	return &dr, nil
}

func (f *FHIRClient) UpsertPatient(ctx context.Context, patient *samplyFhir.Patient) (*samplyFhir.Patient, error) {
	payload, err := json.Marshal(patient)
	if err != nil {
		return nil, fmt.Errorf("marshal patient: %w", err)
	}

	resourceType := "Patient"
	resourceID := *patient.Id

	parentPath := fmt.Sprintf("%s", f.basePath)                                   // for Create
	resourcePath := fmt.Sprintf("%s/%s/%s", f.basePath, resourceType, resourceID) // for Update

	var resp *http.Response

	if patient.Meta == nil || patient.Meta.VersionId == nil || *patient.Meta.VersionId == "" {
		// CREATE
		call := f.svc.Projects.Locations.Datasets.FhirStores.Fhir.Create(parentPath, resourceType, bytes.NewReader(payload))
		call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
		resp, err = call.Do()
	} else {
		// UPDATE
		call := f.svc.Projects.Locations.Datasets.FhirStores.Fhir.Update(resourcePath, bytes.NewReader(payload))
		call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
		call.Header().Set("If-Match", fmt.Sprintf(`W/"%s"`, *patient.Meta.VersionId))
		resp, err = call.Do()
	}

	if err != nil {
		return nil, fmt.Errorf("upsert patient: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, f.readErrorResponse(resp)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	p, err := samplyFhir.UnmarshalPatient(bodyBytes)
	fmt.Println("patient", patient)
	if err != nil {
		return nil, fmt.Errorf("unmarshal patient response: %w", err)
	}
	return &p, nil
}

// This method creates or updates an Observation resource.
func (f *FHIRClient) CreateObservation(ctx context.Context, obs *samplyFhir.Observation) (*samplyFhir.Observation, error) {
	payload, err := json.Marshal(obs)
	if err != nil {
		return nil, fmt.Errorf("marshal Observation: %w", err)
	}
	var resp *http.Response

	parentPath := fmt.Sprintf("%s", f.basePath)
	var o samplyFhir.Observation
	call := f.svc.Projects.Locations.Datasets.FhirStores.Fhir.Create(parentPath, "Observation", bytes.NewReader(payload))
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
	resp, err = call.Do()
	if err != nil {
		return nil, fmt.Errorf("create observation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, f.readErrorResponse(resp)
	}
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return nil, fmt.Errorf("decode observation response: %w", err)
	}
	return &o, nil
}

// a utility function for reading the error response
func (f *FHIRClient) readErrorResponse(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	return fmt.Errorf("fhir client error: status %d, body: %s", resp.StatusCode, string(body))
}
