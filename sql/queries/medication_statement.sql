-- name: CreateMedicationStatement :one
INSERT INTO medication_statements (
    patient_id,
    status,
    medication_code_system,
    medication_code_code,
    medication_code_display,
    dosage_text,
    effective_date_time
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetMedicationStatementByID :one
SELECT * FROM medication_statements
WHERE id = $1 AND patient_id = $2
LIMIT 1;

-- name: ListMedicationStatementsByPatient :many
SELECT * FROM medication_statements
WHERE patient_id = $1
ORDER BY effective_date_time DESC, created_at DESC;

-- name: UpdateMedicationStatement :one
UPDATE medication_statements
SET
    status = $2,
    medication_code_system = $3,
    medication_code_code = $4,
    medication_code_display = $5,
    dosage_text = $6,
    effective_date_time = $7,
    updated_at = NOW()
WHERE id = $1 AND patient_id = $8 -- Ensure update is for the correct patient
RETURNING *;

-- name: DeleteMedicationStatement :exec
DELETE FROM medication_statements
WHERE id = $1 AND patient_id = $2;
