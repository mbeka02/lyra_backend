-- name: CreatePatient :one
INSERT INTO patients(user_id,date_of_birth,allergies) VALUES ($1,$2,$3)RETURNING *;

