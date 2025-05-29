package model

// CreateAllergyIntoleranceRequest for creating a new allergy record for a patient.
// Assumes patient_id will be populated by the service/handler based on context or request.
type CreateAllergyIntoleranceRequest struct {
	PatientID                 int64   `json:"patient_id"`                               // Required if admin/doctor is creating for a patient
	ClinicalStatusCode        string  `json:"clinical_status_code" validate:"required"` // e.g., "active", "inactive"
	ClinicalStatusDisplay     *string `json:"clinical_status_display"`
	CodeSystem                *string `json:"code_system"`
	CodeCode                  string  `json:"code_code" validate:"required"`
	CodeDisplay               string  `json:"code_display" validate:"required"`
	Criticality               *string `json:"criticality"`                 // e.g., "low", "high"
	ReactionManifestationText *string `json:"reaction_manifestation_text"` // Simplified text
}

// UpdateAllergyIntoleranceRequest for updating an existing allergy record.
type UpdateAllergyIntoleranceRequest struct {
	// PatientID not needed here as it's usually part of URL or checked for auth
	ClinicalStatusCode        string  `json:"clinical_status_code" validate:"required"`
	ClinicalStatusDisplay     *string `json:"clinical_status_display"`
	CodeSystem                *string `json:"code_system"`
	CodeCode                  string  `json:"code_code" validate:"required"`
	CodeDisplay               string  `json:"code_display" validate:"required"`
	Criticality               *string `json:"criticality"`
	ReactionManifestationText *string `json:"reaction_manifestation_text"`
}

// AllergyIntoleranceResponse (could be your database.AllergyIntolerance struct directly if fields align)
// For now, let's assume it's the same as database.AllergyIntolerance
// type AllergyIntoleranceResponse database.AllergyIntolerance
