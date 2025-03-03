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
-- name: GetAppointmentSlots :many
WITH aggregated_availability AS (
  SELECT
    a.doctor_id,
    a.start_time, 
    a.end_time, 
    a.interval_minutes 
  FROM availability a
  WHERE a.doctor_id = $1
  AND a.day_of_week = $2
),
doctor_slots AS (
  SELECT 
    slot_timestamp AS slot_start_time,
    slot_timestamp + (aa.interval_minutes * interval '1 minute') AS slot_end_time,
    aa.doctor_id
  FROM aggregated_availability aa,
  LATERAL generate_series(
    ($3::date + aa.start_time)::timestamp,
    ($3::date + aa.end_time - (aa.interval_minutes * interval '1 minute'))::timestamp,
    (aa.interval_minutes * interval '1 minute')
  ) AS slot_timestamp
)
SELECT
  ds.slot_start_time::time AS slot_start_time, 
  ds.slot_end_time::time AS slot_end_time, 
  CASE 
    WHEN appt.appointment_id IS NOT NULL THEN 'booked'
    ELSE 'available'
  END AS slot_status
FROM doctor_slots ds
LEFT JOIN appointments appt
  ON appt.doctor_id = ds.doctor_id
  AND appt.start_time = ds.slot_start_time;
