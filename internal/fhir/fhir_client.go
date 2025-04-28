package fhir

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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
		return nil, f.readErrorResponse(resp, "upsert patient")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("decoding error:%v", err)
	}
	p, err := samplyFhir.UnmarshalPatient(bodyBytes)
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
	call := f.svc.Projects.Locations.Datasets.FhirStores.Fhir.Create(parentPath, "Observation", bytes.NewReader(payload))
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
	resp, err = call.Do()
	if err != nil {
		return nil, fmt.Errorf("create observation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, f.readErrorResponse(resp, "create observation")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("decoding error:%v", err)
	}
	o, err := samplyFhir.UnmarshalObservation(bodyBytes)
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return nil, fmt.Errorf("decode observation response: %w", err)
	}
	return &o, nil
}

// This method creates a new DocumentReference resource in the FHIR store.
func (f *FHIRClient) CreateDocumentReference(ctx context.Context, docRef *samplyFhir.DocumentReference) (*samplyFhir.DocumentReference, error) {
	payload, err := json.Marshal(docRef)
	if err != nil {
		return nil, fmt.Errorf("create: marshal documentreference: %w", err)
	}

	resourceType := "DocumentReference"
	parentPath := fmt.Sprintf("%s/fhir", f.basePath) // Path for Create operation includes /fhir

	// Call the Create API
	call := f.svc.Projects.Locations.Datasets.FhirStores.Fhir.Create(parentPath, resourceType, bytes.NewReader(payload))
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")

	resp, err := call.Do()
	if err != nil {
		// TODO: Consider adding retry logic here for transient network errors
		return nil, fmt.Errorf("create documentreference API call failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code specifically for created
	if resp.StatusCode != http.StatusCreated {
		return nil, f.readErrorResponse(resp, "create documentreference") // Pass operation name
	}

	// Decode the response body (the created resource with server-assigned ID/Meta)
	return f.decodeDocumentReferenceResponse(resp.Body)
}

// This method updates an existing DocumentReference resource in the FHIR store.
// The input docRef MUST have Id and Meta.VersionId populated correctly.
func (f *FHIRClient) UpdateDocumentReference(ctx context.Context, docRef *samplyFhir.DocumentReference) (*samplyFhir.DocumentReference, error) {
	if docRef.Id == nil || *docRef.Id == "" {
		return nil, fmt.Errorf("update documentreference error: missing required Id field")
	}
	if docRef.Meta == nil || docRef.Meta.VersionId == nil || *docRef.Meta.VersionId == "" {
		return nil, fmt.Errorf("update documentreference error: missing required Meta.VersionId field for If-Match header")
	}

	payload, err := json.Marshal(docRef)
	if err != nil {
		return nil, fmt.Errorf("update: marshal documentreference: %w", err)
	}

	resourceType := "DocumentReference"
	resourceID := *docRef.Id
	versionID := *docRef.Meta.VersionId

	// Path for Update operation includes resource type and ID
	resourcePath := fmt.Sprintf("%s/fhir/%s/%s", f.basePath, resourceType, resourceID)

	// Call the Update API
	call := f.svc.Projects.Locations.Datasets.FhirStores.Fhir.Update(resourcePath, bytes.NewReader(payload))
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
	// Set If-Match header for optimistic locking
	call.Header().Set("If-Match", fmt.Sprintf(`W/"%s"`, versionID))

	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("update documentreference API call failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code specifically for OK (Update returns 200 OK)
	if resp.StatusCode != http.StatusOK {
		return nil, f.readErrorResponse(resp, "update documentreference") // Pass operation name
	}

	return f.decodeDocumentReferenceResponse(resp.Body)
}

// Helper function to decode response
func (f *FHIRClient) decodeDocumentReferenceResponse(body io.ReadCloser) (*samplyFhir.DocumentReference, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var dr samplyFhir.DocumentReference
	if err := json.Unmarshal(bodyBytes, &dr); err != nil {
		return nil, fmt.Errorf("error decoding documentreference response: %w. Body: %s", err, string(bodyBytes))
	}
	return &dr, nil
}

func (f *FHIRClient) readErrorResponse(resp *http.Response, operation string) error {
	body, _ := io.ReadAll(resp.Body)

	// Attempt to parse as OperationOutcome for more detailed errors
	var opOutcome samplyFhir.OperationOutcome
	if err := json.Unmarshal(body, &opOutcome); err == nil && len(opOutcome.Issue) > 0 {
		// Format issues nicely
		var issues strings.Builder
		for i, issue := range opOutcome.Issue {
			if i > 0 {
				issues.WriteString("; ")
			}
			issues.WriteString(fmt.Sprintf("Severity: %s, Code: %s, Details: %v", issue.Severity, issue.Code, issue.Diagnostics))
		}
		return fmt.Errorf("fhir client error during '%s': status %d, Outcome: [%s]", operation, resp.StatusCode, issues.String())
	}

	// Fallback to raw body if not an OperationOutcome
	return fmt.Errorf("fhir client error during '%s': status %d, body: %s", operation, resp.StatusCode, string(body))
}
