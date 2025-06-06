package model

import "time"

// internal/model/ehr_observation.go

// CreateObservationRequest for creating a new observation (note).
// Assumes patient_id will be populated by the service/handler based on context or request.
type CreateObservationRequestForDB struct {
	PatientID         int64     `json:"patient_id"`                    // Required if admin/doctor is creating for a patient
	Status            string    `json:"status" validate:"required"`    // e.g., "final", "amended"
	CodeText          string    `json:"code_text" validate:"required"` // Description (e.g., "Consultation Note")
	EffectiveDateTime time.Time `json:"effective_date_time" validate:"required"`
	ValueString       string    `json:"value_string" validate:"required"` // The actual note content
	// SpecialistID   *int64     `json:"specialist_id"` // Optional: if explicitly passing from frontend
}

// UpdateObservationRequest for updating an existing observation.
type UpdateObservationRequest struct {
	Status            string    `json:"status" validate:"required"`
	CodeText          string    `json:"code_text" validate:"required"`
	EffectiveDateTime time.Time `json:"effective_date_time" validate:"required"`
	ValueString       string    `json:"value_string" validate:"required"`
	// SpecialistID   *int64    `json:"specialist_id"`
}

// ObservationResponse can be your database.Observation struct directly if it aligns well.
// type ObservationResponse database.Observation
// CreateObservationRequest defines the input for creating an Observation.
// Example structure for a simple vital sign or quantitative result.
type CreateObservationRequest struct {
	PatientID         int64      `json:"patientId" validate:"required"` // Subject
	SpecialistID      *int64     `json:"specialistId"`                  // Optional Performer (using your ID name)
	Status            int        `json:"status" validate:"required"`    // e.g., "final", "preliminary" (FHIR ObservationStatus code)
	CategoryCode      *string    `json:"categoryCode"`                  // e.g., "vital-signs"
	CategorySystem    *string    `json:"categorySystem"`                // e.g., "http://terminology.hl7.org/CodeSystem/observation-category"
	CategoryDisplay   *string    `json:"categoryDisplay"`
	CodeCode          string     `json:"codeCode" validate:"required"`          // e.g., LOINC code like "8302-2" for Body Height
	CodeSystem        *string    `json:"codeSystem"`                            // e.g., "http://loinc.org"
	CodeDisplay       *string    `json:"codeDisplay"`                           // e.g., "Body height"
	EffectiveDateTime *time.Time `json:"effectiveDateTime" validate:"required"` // When the observation was taken
	ValueQuantity     *float64   `json:"valueQuantity"`                         // The numeric value
	ValueUnit         *string    `json:"valueUnit"`                             // e.g., "cm", "kg", "mmHg"
	ValueSystem       *string    `json:"valueSystem"`                           // e.g., "http://unitsofmeasure.org"
	ValueCode         *string    `json:"ValueCode"`                             // e.g., UCUM code like "cm"
	Note              *string    `json:"note"`                                  // Optional text note
}

type CreateConsultationNoteRequest struct {
	PatientID int64  `json:"patient_id" validate:"required"` // The ID of the patient this note is for
	NoteText  string `json:"note_text" validate:"required"`  // The actual text of the consultation note
	// SpecialistID will be derived from the authenticated user in the handler
}
