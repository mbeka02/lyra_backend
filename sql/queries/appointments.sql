-- name: CreateAppointment :one
INSERT INTO appointments(patient_id,doctor_id,appointment_date) VALUES ($1,$2,$3) RETURNING *;
-- name: GetPatientAppointments :many
SELECT * FROM appointments WHERE patient_id=$1 LIMIT $2 OFFSET $3;
-- name: UpdateAppointmentStatus :exec
UPDATE appointments SET current_status=$1 WHERE appointment_id=$2;
