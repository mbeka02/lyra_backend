package repository

import (
	"fmt"

	samplyFhir "github.com/samply/golang-fhir-models/fhir-models/fhir"

	"github.com/mbeka02/lyra_backend/internal/database"
)

func buildFHIRPatientFromDB(p *database.Patient) *samplyFhir.Patient {
	patientID := fmt.Sprintf("%d", p.PatientID)
	return &samplyFhir.Patient{
		Id: &patientID,
		Meta: &samplyFhir.Meta{
			Profile: []string{"http://hl7.org/fhir/StructureDefinition/Patient"},
		},
		Address: []samplyFhir.Address{{
			Text: &p.Address,
		}},
		Contact: []samplyFhir.PatientContact{{
			Name: &samplyFhir.HumanName{
				Text: &p.EmergencyContactName,
			},
		}},
		Extension: []samplyFhir.Extension{
			{
				Url:         "http://lyra.com/fhir/StructureDefinition/allergies",
				ValueString: &p.Allergies,
			},
			{
				Url:         "http://lyra.com/fhir/StructureDefinition/currentMedication",
				ValueString: &p.CurrentMedication,
			},
		},
	}
}
