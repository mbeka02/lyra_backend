package model

import "time"

type CreateMedicationStatementRequest struct {
	PatientID             int64      `json:"patient_id"`
	Status                string     `json:"status" validate:"required"` // e.g., "active", "completed"
	MedicationCodeSystem  *string    `json:"medication_code_system"`
	MedicationCodeCode    string     `json:"medication_code_code" validate:"required"`
	MedicationCodeDisplay string     `json:"medication_code_display" validate:"required"`
	DosageText            *string    `json:"dosage_text"`
	EffectiveDateTime     *time.Time `json:"effective_date_time"` // Can be null if ongoing
}

type UpdateMedicationStatementRequest struct {
	Status                string     `json:"status" validate:"required"`
	MedicationCodeSystem  *string    `json:"medication_code_system"`
	MedicationCodeCode    string     `json:"medication_code_code" validate:"required"`
	MedicationCodeDisplay string     `json:"medication_code_display" validate:"required"`
	DosageText            *string    `json:"dosage_text"`
	EffectiveDateTime     *time.Time `json:"effective_date_time"`
}
