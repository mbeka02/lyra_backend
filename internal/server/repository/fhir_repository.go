package repository

import (
	"fmt"

	samplyFhir "github.com/samply/golang-fhir-models/fhir-models/fhir"

	"github.com/mbeka02/lyra_backend/internal/database"
)

func buildFHIRPatientFromDB(p *database.Patient) *samplyFhir.Patient {
	patientID := fmt.Sprintf("%d", p.PatientID)

	// Create enum values as pointers
	homeAddressUse := samplyFhir.AddressUseHome
	phoneSystem := samplyFhir.ContactPointSystemPhone
	homeUse := samplyFhir.ContactPointUseHome
	seealsoLinkType := samplyFhir.LinkTypeSeealso

	// Create base patient with ID and meta information
	patient := &samplyFhir.Patient{
		Id: &patientID,
		Meta: &samplyFhir.Meta{
			Profile: []string{"http://hl7.org/fhir/StructureDefinition/Patient"},
			// LastUpdated: &p.UpdatedAt,
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

	// Add link to user account
	userReference := fmt.Sprintf("User/%d", p.UserID)
	patient.Link = []samplyFhir.PatientLink{
		{
			Other: samplyFhir.Reference{
				Reference: &userReference,
			},
			Type: seealsoLinkType,
		},
	}

	return patient
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
