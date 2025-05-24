-- name: CreateDoctor :one
INSERT INTO doctors(user_id,specialization,license_number,description , years_of_experience , county , price_per_hour) VALUES ($1,$2,$3,$4,$5,$6,$7)RETURNING *;
-- name: GetDoctorIdByUserId :one
SELECT doctor_id FROM doctors WHERE user_id=$1;

-- name: ListPatientsUnderDoctorCare :many 
WITH DoctorPatientIds AS (
SELECT DISTINCT a.patient_id
FROM appointments a 
WHERE a.doctor_id=$1
AND a.current_status NOT IN ('cancelled','pending_payment')
)
SELECT
    p.patient_id,
    u.user_id,
    u.date_of_birth,
    u.full_name,
    u.profile_image_url as profile_picture
FROM DoctorPatientIDs dp_ids
JOIN patients p ON dp_ids.patient_id = p.patient_id
JOIN users u ON p.user_id = u.user_id
ORDER BY
    u.full_name ASC;
-- name: GetDoctors :many
SELECT 
    users.full_name, 
    doctors.specialization, 
    doctors.doctor_id, 
    users.profile_image_url, 
    doctors.description, 
    doctors.county, 
    doctors.price_per_hour, 
    doctors.years_of_experience
FROM doctors
INNER JOIN users ON doctors.user_id = users.user_id
WHERE 
    -- Filter for county - returns all results when empty
    (TRIM(@set_county::text) = '' OR doctors.county ILIKE '%' || @set_county::text || '%')
    AND (TRIM(@set_specialization::text) = '' OR doctors.specialization ILIKE '%' || @set_specialization::text || '%')
AND (NULLIF(@set_min_price::text, '')::numeric IS NULL OR doctors.price_per_hour >= NULLIF(@set_min_price::text, '')::numeric)
    AND (NULLIF(@set_max_price::text, '')::numeric IS NULL OR doctors.price_per_hour <= NULLIF(@set_max_price::text, '')::numeric)
    AND (@set_min_experience::int IS NULL OR doctors.years_of_experience >= @set_min_experience::int)
    AND (@set_max_experience::int IS NULL OR doctors.years_of_experience <= @set_max_experience::int)
ORDER BY 
    CASE 
        WHEN @set_sort_by::text = 'price' AND @set_sort_order::text = 'asc' THEN doctors.price_per_hour
        WHEN @set_sort_by::text = 'price' AND @set_sort_order::text = 'desc' THEN doctors.price_per_hour * -1
        WHEN @set_sort_by::text = 'experience' AND @set_sort_order::text = 'asc' THEN doctors.years_of_experience
        WHEN @set_sort_by::text = 'experience' AND @set_sort_order::text = 'desc' THEN doctors.years_of_experience * -1
        ELSE NULL
    END,
    CASE 
        WHEN @set_sort_by::text = 'newest' AND @set_sort_order::text = 'asc' THEN doctors.created_at
        WHEN @set_sort_by::text = 'newest' AND @set_sort_order::text = 'desc' OR @set_sort_by::text NOT IN ('price', 'experience', 'newest') THEN doctors.created_at
        ELSE NULL
    END DESC
LIMIT @set_limit::int OFFSET @set_offset::int;
