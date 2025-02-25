-- name: CreatePatient :one
INSERT INTO patients(user_id,allergies) VALUES ($1,$2)RETURNING *;

