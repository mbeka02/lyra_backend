-- name: CreateAvailability :one
INSERT INTO availability (
  doctor_id, day_of_week, start_time, end_time, is_recurring,interval_minutes
) VALUES (
  $1, $2, $3, $4, $5,$6
) RETURNING *;
-- name: GetAvailabilityByDoctor :many
SELECT * FROM availability WHERE doctor_id=$1;
-- name: GetAvailabilityByDoctorAndDay :many
SELECT * FROM availability WHERE doctor_id=$1 AND day_of_week=$2;
-- name: DeleteAvailabityById :exec
DELETE FROM availability WHERE availability_id=$1 AND doctor_id=$2;
-- name: DeleteAvailabityByDay :exec
DELETE  FROM availability WHERE day_of_week=$1 AND doctor_id=$2;
