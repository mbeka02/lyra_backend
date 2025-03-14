-- name: CreateAvailability :one
INSERT INTO availability (
  doctor_id, day_of_week, start_time, end_time,interval_minutes
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;
-- name: GetAvailabilityByDoctor :many
SELECT * FROM availability WHERE doctor_id=$1;
-- name: GetAvailabilityByDoctorAndDay :many
SELECT * FROM availability WHERE doctor_id=$1 AND day_of_week=$2;
-- name: DeleteAvailabityById :exec
DELETE FROM availability WHERE availability_id=$1 AND doctor_id=$2;
-- name: DeleteAvailabityByDay :exec
DELETE  FROM availability WHERE day_of_week=$1 AND doctor_id=$2;
-- name: GetAppointmentSlots :many
WITH time_slots AS (
  SELECT 
    a.doctor_id,
    slot_time::time AS slot_start_time,
    (slot_time + (a.interval_minutes * interval '1 minute'))::time AS slot_end_time
  FROM availability a,
  LATERAL generate_series(
    ($3::date + a.start_time)::timestamp,
    ($3::date + a.end_time - (a.interval_minutes * interval '1 minute'))::timestamp,
    (a.interval_minutes * interval '1 minute')
  ) AS slot_time
  WHERE a.doctor_id = $1
  AND a.day_of_week = $2
)
SELECT
  ts.slot_start_time,
  ts.slot_end_time,
  CASE 
    WHEN appt.appointment_id IS NOT NULL THEN 'booked'
    ELSE 'available'
  END AS slot_status
FROM time_slots ts
LEFT JOIN appointments appt
  ON appt.doctor_id = ts.doctor_id
  AND appt.start_time::time = ts.slot_start_time
  AND appt.start_time::date = $3::date;
