-- name: CreateDoctor :one
INSERT INTO doctors(user_id,specialization,license_number) VALUES ($1,$2,$3)RETURNING *;

-- name: GetDoctors :many
SELECT full_name,specialization,doctor_id,profile_image_url FROM doctors INNER JOIN users ON doctors.user_id=users.user_id LIMIT $1 OFFSET $2;

