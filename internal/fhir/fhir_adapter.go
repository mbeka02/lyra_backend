package fhir

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	samplyFhir "github.com/samply/golang-fhir-models/fhir-models/fhir"

	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/model"
)

func BuildFHIRPatientFromDB(p *database.Patient, user *database.User) (*samplyFhir.Patient, error) {
	if p == nil || user == nil {
		return nil, fmt.Errorf("Cannot use nil pointers to build the resource")
	}
	patientID := fmt.Sprintf("%d", p.PatientID)

	// Create enum values as pointers
	homeAddressUse := samplyFhir.AddressUseHome
	phoneSystem := samplyFhir.ContactPointSystemPhone
	homeUse := samplyFhir.ContactPointUseHome
	// seealsoLinkType := samplyFhir.LinkTypeSeealso
	contactPointSystemPhone := samplyFhir.ContactPointSystemPhone
	contactPointSystemEmail := samplyFhir.ContactPointSystemEmail
	contactPointUseMobile := samplyFhir.ContactPointUseMobile
	// Create base patient with ID and meta information
	patient := &samplyFhir.Patient{
		Id: &patientID,
		Identifier: []samplyFhir.Identifier{{
			System: stringPtr("http://lyra.com/fhir/Patient/id"),
			Value:  &patientID,
		}},
		Meta: &samplyFhir.Meta{
			// TODO: Add Version ID
			// VersionId: stringPtr(p.FhirVersion),
			Profile: []string{"http://hl7.org/fhir/StructureDefinition/Patient"},
			// LastUpdated: &p.UpdatedAt,
		},
		Name: []samplyFhir.HumanName{{
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
		// Address information
		Address: []samplyFhir.Address{
			{
				Text: &p.Address,
				Use:  &homeAddressUse,
			},
		},
		// Emergency contact information
		Contact: []samplyFhir.PatientContact{
			{
				Name: &samplyFhir.HumanName{
					Text: &p.EmergencyContactName,
				},
				Telecom: []samplyFhir.ContactPoint{
					{
						System: &phoneSystem,
						Value:  &p.EmergencyContactPhone,
						Use:    &homeUse,
					},
				},
				Relationship: []samplyFhir.CodeableConcept{
					{
						Text: stringPtr("Emergency Contact"),
					},
				},
			},
		},
	}

	// Add extension elements for medical and insurance information
	patient.Extension = []samplyFhir.Extension{
		{
			Url:         "http://lyra.com/fhir/StructureDefinition/allergies",
			ValueString: &p.Allergies,
		},
		{
			Url:         "http://lyra.com/fhir/StructureDefinition/currentMedication",
			ValueString: &p.CurrentMedication,
		},
		{
			Url:         "http://lyra.com/fhir/StructureDefinition/pastMedicalHistory",
			ValueString: &p.PastMedicalHistory,
		},
		{
			Url:         "http://lyra.com/fhir/StructureDefinition/familyMedicalHistory",
			ValueString: &p.FamilyMedicalHistory,
		},
	}

	// Add insurance information as coverage reference
	if p.InsuranceProvider != "" || p.InsurancePolicyNumber != "" {
		insuranceText := fmt.Sprintf("Provider: %s, Policy: %s", p.InsuranceProvider, p.InsurancePolicyNumber)
		patient.Extension = append(patient.Extension, samplyFhir.Extension{
			Url:         "http://lyra.com/fhir/StructureDefinition/insuranceInformation",
			ValueString: &insuranceText,
		})
	}

	// // Add link to user account
	// userReference := fmt.Sprintf("User/%d", p.UserID)
	// patient.Link = []samplyFhir.PatientLink{
	// 	{
	// 		Other: samplyFhir.Reference{
	// 			Reference: &userReference,
	// 		},
	// 		Type: seealsoLinkType,
	// 	},
	// }

	return patient, nil
}

// BuildFHIRDocumentReference constructs a DocumentReference FHIR resource.
// It now takes metadata and GCS/file details separately.
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
