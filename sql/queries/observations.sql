-- name: CreateObservation :one
INSERT INTO observations (
    patient_id,
    -- specialist_id, 
    status,
    code_text,
    effective_date_time,
    value_string
) VALUES (
    $1, $2, $3, $4, $5 
) RETURNING *;

-- name: GetObservationByID :one
SELECT * FROM observations
WHERE id = $1 AND patient_id = $2
LIMIT 1;

-- name: ListObservationsByPatient :many
SELECT * FROM observations
WHERE patient_id = $1
ORDER BY effective_date_time DESC;

-- name: UpdateObservation :one
UPDATE observations
SET
    -- specialist_id 
    status = $2, 
    code_text = $3,
    effective_date_time = $4,
    value_string = $5,
    updated_at = NOW()
WHERE id = $1 AND patient_id = $6 
RETURNING *;

-- name: DeleteObservation :exec
DELETE FROM observations
WHERE id = $1 AND patient_id = $2;
