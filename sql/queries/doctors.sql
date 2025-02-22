-- name: CreateDoctor :one
INSERT INTO doctors(user_id,specialization,license_number,description) VALUES ($1,$2,$3,$4)RETURNING *;

-- name: GetDoctors :many
SELECT full_name,specialization,doctor_id,profile_image_url,description FROM doctors INNER JOIN users ON doctors.user_id=users.user_id LIMIT $1 OFFSET $2;

