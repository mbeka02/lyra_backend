// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: patients.sql

package database

import (
	"context"
)

const createPatient = `-- name: CreatePatient :one
INSERT INTO patients(user_id,allergies,current_medication,past_medical_history,family_medical_history,insurance_provider,insurance_policy_number, address,emergency_contact_name , emergency_contact_phone) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)RETURNING patient_id, user_id, address, emergency_contact_name, emergency_contact_phone, allergies, current_medication, past_medical_history, family_medical_history, insurance_provider, insurance_policy_number, created_at, updated_at
`

type CreatePatientParams struct {
	UserID                int64  `json:"user_id"`
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

func (q *Queries) CreatePatient(ctx context.Context, arg CreatePatientParams) (Patient, error) {
	row := q.db.QueryRowContext(ctx, createPatient,
		arg.UserID,
		arg.Allergies,
		arg.CurrentMedication,
		arg.PastMedicalHistory,
		arg.FamilyMedicalHistory,
		arg.InsuranceProvider,
		arg.InsurancePolicyNumber,
		arg.Address,
		arg.EmergencyContactName,
		arg.EmergencyContactPhone,
	)
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

const getPatientIdByUserId = `-- name: GetPatientIdByUserId :one
SELECT patient_id FROM patients WHERE user_id=$1
`

func (q *Queries) GetPatientIdByUserId(ctx context.Context, userID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getPatientIdByUserId, userID)
	var patient_id int64
	err := row.Scan(&patient_id)
	return patient_id, err
}
