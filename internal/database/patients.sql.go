// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: patients.sql

package database

import (
	"context"
)

const createPatient = `-- name: CreatePatient :one
INSERT INTO patients(user_id,allergies) VALUES ($1,$2)RETURNING patient_id, user_id, address, emergency_contact_name, emergency_contact_phone, allergies, current_medication, past_medical_history, family_medical_history, insurance_provider, insurance_policy_number, created_at, updated_at
`

type CreatePatientParams struct {
	UserID    int64
	Allergies string
}

func (q *Queries) CreatePatient(ctx context.Context, arg CreatePatientParams) (Patient, error) {
	row := q.db.QueryRowContext(ctx, createPatient, arg.UserID, arg.Allergies)
	var i Patient
	err := row.Scan(
		&i.PatientID,
		&i.UserID,
		&i.Address,
		&i.EmergencyContactName,
		&i.EmergencyContactPhone,
		&i.Allergies,
		&i.CurrentMedication,
		&i.PastMedicalHistory,
		&i.FamilyMedicalHistory,
		&i.InsuranceProvider,
		&i.InsurancePolicyNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
