-- name: CreateAvailability :one
INSERT INTO availability (
  doctor_id, day_of_week, start_time, end_time, is_recurring
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;
-- GetAvailabilityByDocctor :one
SELECT * FROM availability WHERE doctor_id=$1;
