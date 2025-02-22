package model

import "github.com/mbeka02/lyra_backend/internal/database"

type CreateDoctorRequest struct {
	Specialization string `json:"specialization" validate:"required"`
	LicenseNumber  string `json:"license_number" validate:"required"`
	Description    string `json:"description" validate:"required"`
}

type DoctorDetails struct {
	DoctorID        int64  `json:"doctor_id"`
	FullName        string `json:"full_name"`
	Specialization  string `json:"specialization"`
	ProfileImageUrl string `json:"profile_image_url"`
	Description     string `json:"description"`
}
type GetDoctorsResponse struct {
	HasMore bool            `json:"has_more"`
	Doctors []DoctorDetails `json:"doctors"`
}

func NewDoctorDetails(rows []database.GetDoctorsRow) []DoctorDetails {
	resp := make([]DoctorDetails, 0, len(rows))

	for _, row := range rows {
		resp = append(resp, DoctorDetails{
			DoctorID:        row.DoctorID,
			FullName:        row.FullName,
			Specialization:  row.Specialization,
			ProfileImageUrl: row.ProfileImageUrl,
			Description:     row.Description,
		})
	}

	return resp
}
