-- name: CreateAppointment :one
INSERT INTO appointments(patient_id,doctor_id,start_time,end_time, reason) VALUES ($1,$2,$3,$4,$5) RETURNING *;
-- name: GetPatientAppointments :many
SELECT * FROM appointments WHERE patient_id=$1 LIMIT $2 OFFSET $3;
-- name: UpdateAppointmentStatus :exec
UPDATE appointments SET current_status=$1 WHERE appointment_id=$2;
-- name: DeleteAppointment :exec
DELETE FROM appointments WHERE appointment_id=$1;

