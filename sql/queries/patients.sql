-- name: CreatePatient :one
INSERT INTO patients(user_id,allergies,current_medication,past_medical_history,family_medical_history,insurance_provider,insurance_policy_number, address,emergency_contact_name , emergency_contact_phone) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)RETURNING *;
-- name: GetPatientIdByUserId :one
SELECT patient_id FROM patients WHERE user_id=$1;

-- name: UpdateFhirVersionId :exec
UPDATE  patients SET fhir_version=$1 WHERE patient_id=$2; 
