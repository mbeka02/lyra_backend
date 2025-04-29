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
	"golang.org/x/oauth2/google"
	healthcare "google.golang.org/api/healthcare/v1"
	"google.golang.org/api/option"
)

// FHIRClient wraps Google Cloud Healthcare FHIR API calls.
// It handles JSON marshalling/unmarshalling of FHIR resources.
type FHIRClient struct {
	svc        *healthcare.Service
	basePath   string
	baseApiUrl string // Base API endpoint:
	client     *http.Client
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
	// 1. Get the default authenticated HTTP client using ADC
	// Ensure the context has the necessary credentials available (e.g., running on GCP, GOOGLE_APPLICATION_CREDENTIALS env var)
	httpClient, err := google.DefaultClient(ctx, healthcare.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient error: %w", err)
	}

	// 2. Create the healthcare service, passing the authenticated client
	// Prepend the authenticated client option to any user-provided options
	finalOpts := append([]option.ClientOption{option.WithHTTPClient(httpClient)}, opts...)
	svc, err := healthcare.NewService(ctx, finalOpts...)
	if err != nil {
		return nil, fmt.Errorf("healthcare.NewService error: %w", err)
	}

	// 3. Construct path segments
	basePath := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s",
		config.ProjectID, config.DatasetLocation, config.DatasetID, config.FHIRStoreID)
	// Base URL for the Healthcare API v1
	baseApiUrl := "https://healthcare.googleapis.com/v1" // Adjust if using a different endpoint/version

	// 4. Return the client struct containing both service and http client
	return &FHIRClient{
		svc:        svc,
		client:     httpClient, // Store the client
		basePath:   basePath,
		baseApiUrl: baseApiUrl,
	}, nil
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
	if err != nil {
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
	parentPath := fmt.Sprintf("%s", f.basePath) // Path for Create operation includes /fhir

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
	resourcePath := fmt.Sprintf("%s/%s/%s", f.basePath, resourceType, resourceID)

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

// This method performs a search for DocumentReference resources based on query parameters.
// NOTE: It uses a http client to make requests instead of the google package
func (f *FHIRClient) SearchDocumentReferences(ctx context.Context, queryParams string) (*samplyFhir.Bundle, error) {
	resourceType := "DocumentReference"

	// Construct the full API endpoint URL for the GET search
	// Example: https://healthcare.googleapis.com/v1/projects/p/locations/l/datasets/d/fhirStores/f/fhir/DocumentReference?subject=Patient/123
	fullApiPath := fmt.Sprintf("%s/%s/fhir/%s", f.baseApiUrl, f.basePath, resourceType)

	// Append query parameters correctly
	if queryParams != "" {
		queryParams = strings.TrimPrefix(queryParams, "?")
		queryParams = strings.TrimPrefix(queryParams, "&")
		fullApiPath = fmt.Sprintf("%s?%s", fullApiPath, queryParams)
	}

	// Create the HTTP GET request object
	req, err := http.NewRequestWithContext(ctx, "GET", fullApiPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	// Set required headers
	req.Header.Set("Accept", "application/fhir+json")
	// Authentication headers are handled automatically by the httpClient obtained from google.DefaultClient

	//  Execute the request using the stored authenticated client
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search documentreferences (GET %s) API call failed: %w", fullApiPath, err)
	}
	defer resp.Body.Close()

	//  Handle the response (status check)
	if resp.StatusCode != http.StatusOK {
		// Use the existing helper, passing the operation description
		return nil, f.readErrorResponse(resp, fmt.Sprintf("search documentreferences (GET %s)", fullApiPath))
	}

	//  Decode the response body (should be a Bundle)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading search response body: %w", err)
	}

	var bundle samplyFhir.Bundle
	if err := json.Unmarshal(bodyBytes, &bundle); err != nil {
		// Include raw body in error for debugging decode failures
		return nil, fmt.Errorf("error decoding search result bundle: %w. Body: %s", err, string(bodyBytes))
	}

	//  Validate bundle type (optional but recommended)
	if bundle.Type != samplyFhir.BundleTypeSearchset {
		fmt.Printf("Warning: Expected searchset bundle, got type %s\n", bundle.Type)
	}

	return &bundle, nil
}

//OLD FUNCTION
/*
func (f *FHIRClient) SearchDocumentReferences(ctx context.Context, queryParams string) (*samplyFhir.Bundle, error) {
	resourceType := "DocumentReference"

	// Prepare the query parameters
	if queryParams == "" {
		queryParams = "_type=" + resourceType
	} else if !strings.Contains(queryParams, "_type=") {
		queryParams += "&_type=" + resourceType
	}
	// Construct the base URL path for DocumentReference resources
	// Add query parameters to the base path if provided
	// Construct the URL with the query parameters
	parentPath := fmt.Sprintf("%s", f.basePath)
	// Create the search request
	req := &healthcare.SearchResourcesRequest{
		ResourceType: resourceType,
	}

	// Create the search call
	call := f.svc.Projects.Locations.Datasets.FhirStores.Fhir.SearchType(parentPath, resourceType, req)

	// Set necessary headers
	call.Header().Set("Accept", "application/fhir+json")

	// Execute the request
	resp, err := call.Do()
	if err != nil {
		// Consider adding retry logic here for transient network errors
		return nil, fmt.Errorf("search documentreferences API call failed: %w", err)
	}
	defer resp.Body.Close()

	// Search usually returns 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, f.readErrorResponse(resp, "search documentreferences")
	}

	// Decode the response body (should be a Bundle of type searchset)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading search response body: %w", err)
	}

	var bundle samplyFhir.Bundle
	if err := json.Unmarshal(bodyBytes, &bundle); err != nil {
		return nil, fmt.Errorf("error decoding search result bundle: %w. Body: %s", err, string(bodyBytes))
	}

	// Validate that it's a searchset bundle
	if bundle.Type != samplyFhir.BundleTypeSearchset {
		fmt.Printf("Warning: Expected searchset bundle, got type %v\n", bundle.Type)
	}

	return &bundle, nil
}
*/
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
			issues.WriteString(fmt.Sprintf("Severity: %s, Code: %s, Details: %s", issue.Severity, issue.Code, *issue.Diagnostics))
		}
		return fmt.Errorf("fhir client error during '%s': status %d, Outcome: [%s]", operation, resp.StatusCode, issues.String())
	}

	// Fallback to raw body if not an OperationOutcome
	return fmt.Errorf("fhir client error during '%s': status %d, body: %s", operation, resp.StatusCode, string(body))
}
