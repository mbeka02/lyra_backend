-- name: CreateAppointment :one
INSERT INTO appointments(patient_id,doctor_id,start_time,end_time, reason) VALUES ($1,$2,$3,$4,$5) RETURNING *;
-- name: GetPatientAppointments :many
SELECT
a.*,
d.specialization,
u.full_name AS doctor_name,
u.profile_image_url AS doctor_profile_image_url
FROM 
appointments a
JOIN 
doctors d ON a.doctor_id = d.doctor_id
JOIN 
users u ON d.user_id = u.user_id
WHERE a.patient_id=$1
AND (@status::text = '' OR a.current_status::text = @status::text)
AND DATE(a.start_time) BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '1 day'* @set_interval::integer;

-- name: GetDoctorAppointments :many
SELECT 
a.*, 
u.full_name AS patient_name,
u.profile_image_url AS patient_profile_image_url
FROM appointments a 
JOIN 
patients p ON a.patient_id=p.patient_id
JOIN
users u ON p.user_id = u.user_id
WHERE a.doctor_id=$1
AND (@status::text = '' OR a.current_status::text = @status::text)
AND DATE(a.start_time) BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '1 day'* @set_interval::integer;

-- name: GetAppointmentIDs :many
WITH params AS (
  SELECT
    @id::bigint AS id,
    @role::text   AS role
)
SELECT appointment_id
FROM appointments
CROSS JOIN params
WHERE current_status = 'completed'
  AND (
    (params.role = 'doctor'  AND doctor_id  = params.id)
    OR
    (params.role = 'patient' AND patient_id = params.id)
  )
ORDER BY appointment_id;
-- name: UpdateAppointmentStatus :exec
UPDATE appointments SET current_status=$1 WHERE appointment_id=$2;

-- name: DeleteAppointment :exec
DELETE FROM appointments WHERE appointment_id=$1;

