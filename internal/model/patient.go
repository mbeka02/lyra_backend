package model

type CreatePatientRequest struct {
	Allergies             string `json:"allergies"`
	CurrentMedication     string `json:"current_medication"`
	PastMedicalHistory    string `json:"past_medical_history"`
	FamilyMedicalHistory  string `json:"family_medical_history"`
	InsuranceProvider     string `json:"insurance_provider"`
	InsurancePolicyNumber string `json:"insurance_policy_number"`
	Address               string `json:"address"`
	EmergencyContactName  string `json:"emergency_contact_name"`
	EmergencyContactPhone string `json:"emergency_contact_phone"`
}
