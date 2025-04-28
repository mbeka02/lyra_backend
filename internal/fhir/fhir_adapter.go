package fhir

import (
	"fmt"

	samplyFhir "github.com/samply/golang-fhir-models/fhir-models/fhir"

	"github.com/mbeka02/lyra_backend/internal/database"
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

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
