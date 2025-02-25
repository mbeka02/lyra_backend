-- name: CreatePatient :one
INSERT INTO patients(user_id,allergies,current_medication,past_medical_history,family_medical_history,insurance_provider,insurance_policy_number) VALUES ($1,$2,$3,$4,$5,$6,$7)RETURNING *;

