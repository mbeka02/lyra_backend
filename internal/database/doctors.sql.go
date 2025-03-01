// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: doctors.sql

package database

import (
	"context"
)

const createDoctor = `-- name: CreateDoctor :one
INSERT INTO doctors(user_id,specialization,license_number,description , years_of_experience , county , price_per_hour) VALUES ($1,$2,$3,$4,$5,$6,$7)RETURNING doctor_id, user_id, description, specialization, years_of_experience, county, price_per_hour, license_number, created_at, updated_at
`

type CreateDoctorParams struct {
	UserID            int64  `json:"user_id"`
	Specialization    string `json:"specialization"`
	LicenseNumber     string `json:"license_number"`
	Description       string `json:"description"`
	YearsOfExperience int32  `json:"years_of_experience"`
	County            string `json:"county"`
	PricePerHour      string `json:"price_per_hour"`
}

func (q *Queries) CreateDoctor(ctx context.Context, arg CreateDoctorParams) (Doctor, error) {
	row := q.db.QueryRowContext(ctx, createDoctor,
		arg.UserID,
		arg.Specialization,
		arg.LicenseNumber,
		arg.Description,
		arg.YearsOfExperience,
		arg.County,
		arg.PricePerHour,
	)
	var i Doctor
	err := row.Scan(
		&i.DoctorID,
		&i.UserID,
		&i.Description,
		&i.Specialization,
		&i.YearsOfExperience,
		&i.County,
		&i.PricePerHour,
		&i.LicenseNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getDoctorIdByUserId = `-- name: GetDoctorIdByUserId :one
SELECT doctor_id FROM doctors WHERE user_id=$1
`

func (q *Queries) GetDoctorIdByUserId(ctx context.Context, userID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getDoctorIdByUserId, userID)
	var doctor_id int64
	err := row.Scan(&doctor_id)
	return doctor_id, err
}

const getDoctors = `-- name: GetDoctors :many
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
    (TRIM($1::text) = '' OR doctors.county ILIKE '%' || $1::text || '%')
    AND (TRIM($2::text) = '' OR doctors.specialization ILIKE '%' || $2::text || '%')
AND (NULLIF($3::text, '')::numeric IS NULL OR doctors.price_per_hour >= NULLIF($3::text, '')::numeric)
    AND (NULLIF($4::text, '')::numeric IS NULL OR doctors.price_per_hour <= NULLIF($4::text, '')::numeric)
    AND ($5::int IS NULL OR doctors.years_of_experience >= $5::int)
    AND ($6::int IS NULL OR doctors.years_of_experience <= $6::int)
ORDER BY 
    CASE 
        WHEN $7::text = 'price' AND $8::text = 'asc' THEN doctors.price_per_hour
        WHEN $7::text = 'price' AND $8::text = 'desc' THEN doctors.price_per_hour * -1
        WHEN $7::text = 'experience' AND $8::text = 'asc' THEN doctors.years_of_experience
        WHEN $7::text = 'experience' AND $8::text = 'desc' THEN doctors.years_of_experience * -1
        ELSE NULL
    END,
    CASE 
        WHEN $7::text = 'newest' AND $8::text = 'asc' THEN doctors.created_at
        WHEN $7::text = 'newest' AND $8::text = 'desc' OR $7::text NOT IN ('price', 'experience', 'newest') THEN doctors.created_at
        ELSE NULL
    END DESC
LIMIT $10::int OFFSET $9::int
`

type GetDoctorsParams struct {
	SetCounty         string `json:"set_county"`
	SetSpecialization string `json:"set_specialization"`
	SetMinPrice       string `json:"set_min_price"`
	SetMaxPrice       string `json:"set_max_price"`
	SetMinExperience  int32  `json:"set_min_experience"`
	SetMaxExperience  int32  `json:"set_max_experience"`
	SetSortBy         string `json:"set_sort_by"`
	SetSortOrder      string `json:"set_sort_order"`
	SetOffset         int32  `json:"set_offset"`
	SetLimit          int32  `json:"set_limit"`
}

type GetDoctorsRow struct {
	FullName          string `json:"full_name"`
	Specialization    string `json:"specialization"`
	DoctorID          int64  `json:"doctor_id"`
	ProfileImageUrl   string `json:"profile_image_url"`
	Description       string `json:"description"`
	County            string `json:"county"`
	PricePerHour      string `json:"price_per_hour"`
	YearsOfExperience int32  `json:"years_of_experience"`
}

func (q *Queries) GetDoctors(ctx context.Context, arg GetDoctorsParams) ([]GetDoctorsRow, error) {
	rows, err := q.db.QueryContext(ctx, getDoctors,
		arg.SetCounty,
		arg.SetSpecialization,
		arg.SetMinPrice,
		arg.SetMaxPrice,
		arg.SetMinExperience,
		arg.SetMaxExperience,
		arg.SetSortBy,
		arg.SetSortOrder,
		arg.SetOffset,
		arg.SetLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDoctorsRow
	for rows.Next() {
		var i GetDoctorsRow
		if err := rows.Scan(
			&i.FullName,
			&i.Specialization,
			&i.DoctorID,
			&i.ProfileImageUrl,
			&i.Description,
			&i.County,
			&i.PricePerHour,
			&i.YearsOfExperience,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
