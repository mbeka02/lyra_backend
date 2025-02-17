-- name: CreateSpecialist :one
INSERT INTO specialists(user_id,specialization,license_number) VALUES ($1,$2,$3)RETURNING *;

-- name: GetSpecialists :many
SELECT * FROM specialists LIMIT $1 OFFSET $2;

