-- name: CreateAvailability :one
INSERT INTO availability (
  doctor_id, day_of_week, start_time, end_time, is_recurring,interval_minutes
) VALUES (
  $1, $2, $3, $4, $5,$6
) RETURNING *;
-- name: GetAvailabilityByDoctor :many
SELECT doctor_id , day_of_week, start_time , end_time , is_recurring , interval_minutes FROM availability WHERE doctor_id=$1;
-- name: GetAvailabilityByDoctorAndDay :many
SELECT doctor_id , day_of_week, start_time , end_time , is_recurring , interval_minutes FROM availability WHERE doctor_id=$1 AND day_of_week=$2;
