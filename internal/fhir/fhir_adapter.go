package fhir

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	samplyFhir "github.com/samply/golang-fhir-models/fhir-models/fhir"

	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/model"
)

// BuildFHIRObservationFromNote constructs an Observation FHIR resource for a consultation note.
func BuildFHIRObservationFromNote(
	patientID int64, // The FHIR Patient ID (e.g., "25" which becomes "Patient/25")
	specialistID int64, // The FHIR Practitioner ID (e.g., "10" which becomes "Practitioner/10")
	noteText string,
	effectiveTime time.Time, // The time the observation (note) was made or is effective
) (*samplyFhir.Observation, error) {
	if noteText == "" {
		return nil, fmt.Errorf("noteText cannot be empty")
	}

	obs := &samplyFhir.Observation{}

	// Identifier ( Recommended for resource instance uniqueness)
	obs.Identifier = []samplyFhir.Identifier{{
		System: stringPtr("urn:ietf:rfc:3986"), // UUID URN
		Value:  stringPtr("urn:uuid:" + uuid.NewString()),
	}}

	// Status (Required)
	// For a finalized note, 'final' is appropriate.
	obs.Status = samplyFhir.ObservationStatusFinal

	// Category (Optional, but useful for grouping)
	obs.Category = []samplyFhir.CodeableConcept{
		{
			Coding: []samplyFhir.Coding{
				{
					System:  stringPtr("http://terminology.hl7.org/CodeSystem/observation-category"),
					Code:    stringPtr("notes"), // Standard category for notes
					Display: stringPtr("Notes"),
				},
			},
			Text: stringPtr("Clinical Notes"),
		},
	}

	// Code (Required - What kind of observation is this?)
	// Use a consistent code for "Consultation Note" within your Lyra system.
	obs.Code = samplyFhir.CodeableConcept{
		Coding: []samplyFhir.Coding{
			{
				System:  stringPtr("urn:lyra:codesystem:observation-type"), // Your custom system URI
				Code:    stringPtr("CONSULTATION_NOTE"),
				Display: stringPtr("Consultation Note"),
			},
		},
		Text: stringPtr("Consultation Note"),
	}

	// Subject (Required - Who is this observation about?)
	obs.Subject = &samplyFhir.Reference{
		Reference: stringPtr(fmt.Sprintf("Patient/%d", patientID)),
		Type:      stringPtr("Patient"),
	}

	// EffectiveDateTime (Required - When was this observation made or relevant?)
	effectiveDateTimeStr := effectiveTime.Format(time.RFC3339Nano) // FHIR instant format
	obs.EffectiveDateTime = &effectiveDateTimeStr

	// Issued (Optional - When was this version of the observation created?)
	// The FHIR server often sets meta.lastUpdated.
	// issuedTimeStr := time.Now().Format(time.RFC3339Nano)
	// obs.Issued = &issuedTimeStr

	// Performer (Optional but Recommended - Who made the observation/note?)
	obs.Performer = []samplyFhir.Reference{
		{
			Reference: stringPtr(fmt.Sprintf("Practitioner/%d", specialistID)), // Assuming specialistID is the Practitioner ID
			Type:      stringPtr("Practitioner"),
		},
	}

	// Value[x] (Required - The actual observation content)
	// For a textual note, we use valueString.
	obs.ValueString = stringPtr(noteText)

	return obs, nil
}

func BuildFHIRPatientFromDB(p *database.Patient, user *database.User) (*samplyFhir.Patient, error) {
	if p == nil || user == nil {
		return nil, fmt.Errorf("Cannot use nil pointers to build the resource")
	}
	patientIDStr := fmt.Sprintf("%d", p.PatientID) // Changed from p.PatientID to patientIDStr for consistency

	// Create enum values as pointers
	homeAddressUse := samplyFhir.AddressUseHome
	contactPointSystemPhone := samplyFhir.ContactPointSystemPhone
	contactPointSystemEmail := samplyFhir.ContactPointSystemEmail
	contactPointUseMobile := samplyFhir.ContactPointUseMobile
	// Emergency contact specifics
	emergencyPhoneSystem := samplyFhir.ContactPointSystemPhone
	emergencyContactUse := samplyFhir.ContactPointUseHome

	patient := &samplyFhir.Patient{
		Id: &patientIDStr,
		Identifier: []samplyFhir.Identifier{{
			System: stringPtr("http://lyra.com/fhir/Patient/id"),
			Value:  &patientIDStr,
		}},
		Meta: &samplyFhir.Meta{
			// VersionId: stringPtr(p.FhirVersion), // Assuming p has FhirVersion
			Profile: []string{"http://hl7.org/fhir/StructureDefinition/Patient"},
			// LastUpdated: stringPtr(user.UpdatedAt.Format(time.RFC3339Nano)), // Example
		},
		Name: []samplyFhir.HumanName{{
			Text:  stringPtr(user.FullName),
			Given: []string{user.FullName},
		}},
		Telecom: []samplyFhir.ContactPoint{
			{
				System: &contactPointSystemPhone,
				Value:  stringPtr(user.TelephoneNumber),
				Use:    &contactPointUseMobile,
			},
			{
				System: &contactPointSystemEmail,
				Value:  stringPtr(user.Email),
			},
		},
		Address: []samplyFhir.Address{
			{
				Text: stringPtr(p.Address),
				Use:  &homeAddressUse,
			},
		},
		Contact: []samplyFhir.PatientContact{
			{
				Name: &samplyFhir.HumanName{
					Text: stringPtr(p.EmergencyContactName),
				},
				Telecom: []samplyFhir.ContactPoint{
					{
						System: &emergencyPhoneSystem,
						Value:  stringPtr(p.EmergencyContactPhone),
						Use:    &emergencyContactUse,
					},
				},
				Relationship: []samplyFhir.CodeableConcept{
					{
						Text: stringPtr("Emergency Contact"),
					},
				},
			},
		},
		Photo: []samplyFhir.Attachment{},
	}

	// --- CONDITIONALLY ADD PROFILE PICTURE ---
	// Only add the photo entry if user.ProfileImageUrl is not empty.
	if user.ProfileImageUrl != "" {
		patient.Photo = append(patient.Photo, samplyFhir.Attachment{
			Url: stringPtr(user.ProfileImageUrl), // Use the URL directly
			// Title: stringPtr(fmt.Sprintf("Profile picture for %s", user.FullName)), // Optional title
			// ContentType can be omitted if the URL itself resolves to the correct content type
		})
	}

	// Add extension elements for medical and insurance information
	// (Keeping your existing extension structure)
	patient.Extension = []samplyFhir.Extension{
		{
			Url:         "http://lyra.com/fhir/StructureDefinition/allergies",
			ValueString: stringPtr(p.Allergies),
		},
		{
			Url:         "http://lyra.com/fhir/StructureDefinition/currentMedication",
			ValueString: stringPtr(p.CurrentMedication),
		},
		{
			Url:         "http://lyra.com/fhir/StructureDefinition/pastMedicalHistory",
			ValueString: stringPtr(p.PastMedicalHistory),
		},
		{
			Url:         "http://lyra.com/fhir/StructureDefinition/familyMedicalHistory",
			ValueString: stringPtr(p.FamilyMedicalHistory),
		},
	}
	// Add insurance as extension (from your existing code)
	if p.InsuranceProvider != "" || p.InsurancePolicyNumber != "" {
		insuranceText := fmt.Sprintf("Provider: %s, Policy: %s", p.InsuranceProvider, p.InsurancePolicyNumber)
		patient.Extension = append(patient.Extension, samplyFhir.Extension{
			Url:         "http://lyra.com/fhir/StructureDefinition/insuranceInformation",
			ValueString: stringPtr(insuranceText),
		})
	}

	// Adding BirthDate
	if !user.DateOfBirth.IsZero() { // Check if DateOfBirth is a valid, non-zero date
		patient.BirthDate = stringPtr(user.DateOfBirth.Format("2006-01-02")) // FHIR date format
	}

	// if user.GenderField != "" { // Example if you add a GenderField to your user table
	// switch user.GenderField {
	// case "male":
	// patient.Gender = samplyFhir.AdministrativeGenderMale.Enum()
	// case "female":
	// patient.Gender = samplyFhir.AdministrativeGenderFemale.Enum()
	// default:
	// patient.Gender = samplyFhir.AdministrativeGenderUnknown.Enum()
	// }
	// }

	return patient, nil
}

// BuildFHIRDocumentReference constructs a DocumentReference FHIR resource.
func BuildFHIRDocumentReference(
	metadata model.CreateDocumentReferenceRequest,
	gcsUrl string,
	contentType string,
	sizeBytes int64,
	attachmentCreation *time.Time, // Optional: Actual creation time of the file if known
) (*samplyFhir.DocumentReference, error) {
	docRef := &samplyFhir.DocumentReference{}

	docRef.MasterIdentifier = &samplyFhir.Identifier{
		System: stringPtr("urn:ietf:rfc:3986"),
		Value:  stringPtr("urn:uuid:" + uuid.NewString()),
	}
	docRef.Status = samplyFhir.DocumentReferenceStatusCurrent

	if metadata.DocTypeCode != nil || metadata.DocTypeDisplay != nil {
		docType := samplyFhir.CodeableConcept{}
		if metadata.DocTypeDisplay != nil {
			docType.Text = stringPtr(*metadata.DocTypeDisplay)
		}
		if metadata.DocTypeCode != nil {
			coding := samplyFhir.Coding{
				Code: stringPtr(*metadata.DocTypeCode),
			}
			if metadata.DocTypeDisplay != nil {
				coding.Display = stringPtr(*metadata.DocTypeDisplay)
			}
			docType.Coding = []samplyFhir.Coding{coding}
		}
		docRef.Type = &docType
	}

	docRef.Subject = &samplyFhir.Reference{
		Reference: stringPtr(fmt.Sprintf("Patient/%d", metadata.PatientID)),
		Type:      stringPtr("Patient"),
	}

	// --- Author (using SpecialistID or Patient) ---
	// A document *can* be shared, but FHIR `author` usually refers to the creator
	// of the *reference* or the document itself. Access control determines who can *see* it.
	// If a specialist uploads it, they are the author. If patient uploads, they are.
	if metadata.SpecialistID != nil {
		docRef.Author = []samplyFhir.Reference{
			{
				// Assuming Practitioner resources exist or IDs are consistent
				Reference: stringPtr(fmt.Sprintf("Practitioner/%d", *metadata.SpecialistID)),
				Type:      stringPtr("Practitioner"),
			},
		}
	} else {
		// Default to Patient as author if no specialist ID provided
		docRef.Author = []samplyFhir.Reference{
			{
				Reference: stringPtr(fmt.Sprintf("Patient/%d", metadata.PatientID)),
				Type:      stringPtr("Patient"),
			},
		}
	}

	now := time.Now().Format(time.RFC3339Nano)
	docRef.Date = &now // When the reference was created

	attachment := samplyFhir.Attachment{
		ContentType: stringPtr(contentType),
		Url:         stringPtr(gcsUrl),
		Title:       metadata.Title, // Use the provided title
	}
	if sizeBytes > 0 {
		sizeBytesInt := int(sizeBytes)
		attachment.Size = &sizeBytesInt
	}
	if attachmentCreation != nil {
		creationTimeStr := attachmentCreation.Format(time.RFC3339Nano)
		attachment.Creation = &creationTimeStr
	}

	docRef.Content = []samplyFhir.DocumentReferenceContent{
		{Attachment: attachment},
	}

	return docRef, nil
}

func BuildFHIRObservation(req model.CreateObservationRequest) (*samplyFhir.Observation, error) {
	obs := &samplyFhir.Observation{}
	// ... (Identifier, Status, Category, Code - same as before) ...
	obs.Identifier = []samplyFhir.Identifier{{
		System: stringPtr("urn:ietf:rfc:3986"),
		Value:  stringPtr("urn:uuid:" + uuid.NewString()),
	}}
	obs.Status = samplyFhir.ObservationStatus(req.Status)
	if req.CategoryCode != nil || req.CategoryDisplay != nil {
		category := samplyFhir.CodeableConcept{}
		if req.CategoryDisplay != nil {
			category.Text = stringPtr(*req.CategoryDisplay)
		}
		if req.CategoryCode != nil {
			coding := samplyFhir.Coding{
				System: stringPtrIfNotEmpty(req.CategorySystem),
				Code:   stringPtr(*req.CategoryCode),
			}
			if req.CategoryDisplay != nil {
				coding.Display = stringPtr(*req.CategoryDisplay)
			}
			category.Coding = []samplyFhir.Coding{coding}
		}
		obs.Category = []samplyFhir.CodeableConcept{category}
	}
	code := samplyFhir.CodeableConcept{}
	if req.CodeDisplay != nil {
		code.Text = stringPtr(*req.CodeDisplay)
	}
	codeCoding := samplyFhir.Coding{
		System: stringPtrIfNotEmpty(req.CodeSystem),
		Code:   stringPtr(req.CodeCode),
	}
	if req.CodeDisplay != nil {
		codeCoding.Display = stringPtr(*req.CodeDisplay)
	}
	code.Coding = []samplyFhir.Coding{codeCoding}
	obs.Code = code

	// --- Subject (Link to Patient) ---
	obs.Subject = &samplyFhir.Reference{
		Reference: stringPtr(fmt.Sprintf("Patient/%d", req.PatientID)),
		Type:      stringPtr("Patient"),
	}

	// --- Effective DateTime ---
	if req.EffectiveDateTime != nil {
		effectiveTimeStr := req.EffectiveDateTime.Format(time.RFC3339Nano)
		obs.EffectiveDateTime = &effectiveTimeStr
	} else {
		return nil, fmt.Errorf("effectiveDateTime is required for Observation")
	}

	// --- Performer (using SpecialistID) ---
	if req.SpecialistID != nil {
		obs.Performer = []samplyFhir.Reference{
			{
				// Assuming Practitioner resources exist or IDs are consistent
				Reference: stringPtr(fmt.Sprintf("Practitioner/%d", *req.SpecialistID)),
				Type:      stringPtr("Practitioner"),
			},
		}
	}

	// value
	if req.ValueQuantity != nil {
		// 1. Get the float64 value
		floatVal := *req.ValueQuantity

		// 2. Format the float64 to a string representation.
		// 'g' format is generally good for preserving precision compactly.
		// -1 precision means use the smallest number of digits necessary.
		valueStr := strconv.FormatFloat(floatVal, 'g', -1, 64)

		// 3. Create a json.Number from the string.
		jsonNum := json.Number(valueStr)

		// 4. Create the Quantity struct using the *pointer* to the json.Number.
		valueQuantity := samplyFhir.Quantity{
			Value:  &jsonNum, // Assign the address of the json.Number
			Unit:   stringPtrIfNotEmpty(req.ValueUnit),
			System: stringPtrIfNotEmpty(req.ValueSystem),
			Code:   stringPtrIfNotEmpty(req.ValueCode),
		}
		obs.ValueQuantity = &valueQuantity
	}
	// ... (Note, Issued - same as before) ...
	if req.Note != nil {
		obs.Note = []samplyFhir.Annotation{
			{Text: *req.Note},
		}
	}

	return obs, nil
}

// Helper function to create string pointers only if the string is not empty (reuse from previous suggestion)
func stringPtrIfNotEmpty(s *string) *string {
	if s != nil && *s != "" {
		return s
	}
	return nil
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
