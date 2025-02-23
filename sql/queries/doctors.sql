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
    ($1 = '' OR doctors.county ILIKE $1) -- County filter
ORDER BY 
    CASE 
        WHEN $2 = 'price' AND $3 = 'asc' THEN doctors.price_per_hour 
        WHEN $2 = 'price' AND $3 = 'desc' THEN doctors.price_per_hour * -1
        WHEN $2 = 'experience' AND $3 = 'asc' THEN doctors.years_of_experience
        WHEN $2 = 'experience' AND $3 = 'desc' THEN doctors.years_of_experience * -1
    END
LIMIT $4  OFFSET $5;
