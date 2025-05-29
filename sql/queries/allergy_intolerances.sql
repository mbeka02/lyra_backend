-- name: CreateAllergyIntolerance :one
INSERT INTO allergy_intolerances (
    patient_id,
    clinical_status_code,
    clinical_status_display,
    code_system,
    code_code,
    code_display,
    criticality,
    reaction_manifestation_text
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetAllergyIntoleranceByID :one
SELECT * FROM allergy_intolerances
WHERE id = $1 AND patient_id = $2
LIMIT 1;

-- name: ListAllergyIntolerancesByPatient :many
SELECT * FROM allergy_intolerances
WHERE patient_id = $1
ORDER BY created_at DESC;

-- name: UpdateAllergyIntolerance :one
UPDATE allergy_intolerances
SET
    clinical_status_code = $2,
    clinical_status_display = $3,
    code_system = $4,
    code_code = $5,
    code_display = $6,
    criticality = $7,
    reaction_manifestation_text = $8,
    updated_at = NOW()
WHERE id = $1 AND patient_id = $9 -- Ensures update is for the correct patient
RETURNING *;

-- name: DeleteAllergyIntolerance :exec
DELETE FROM allergy_intolerances
WHERE id = $1 AND patient_id = $2;
