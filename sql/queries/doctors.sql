-- name: CreateDoctor :one
INSERT INTO doctors(user_id,specialization,license_number,description , years_of_experience , county , price_per_hour) VALUES ($1,$2,$3,$4,$5,$6,$7)RETURNING *;

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
--FIXME:get this to return all results when county is an empty string or null
    (NULLIF(@set_county::text, '') IS NULL OR doctors.county ILIKE @set_county::text) -- County filter
ORDER BY 
    CASE 
        WHEN @set_sort_by::text = 'price' AND @set_sort_order::text = 'asc' THEN doctors.price_per_hour  
        WHEN @set_sort_by::text = 'price' AND @set_sort_order::text = 'desc' THEN doctors.price_per_hour * -1
        WHEN @set_sort_by::text = 'experience' AND @set_sort_order::text = 'asc' THEN doctors.years_of_experience
        WHEN @set_sort_by::text = 'experience' AND @set_sort_order::text = 'desc' THEN doctors.years_of_experience * -1
    END
LIMIT @set_limit::int OFFSET @set_offset::int;
