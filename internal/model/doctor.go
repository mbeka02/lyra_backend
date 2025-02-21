package model

type CreateDoctorRequest struct {
	Specialization string `json:"specialization" validate:"required"`
	LicenseNumber  string `json:"license_number" validate:"required"`
	Description    string `json:"description" validate:"required"`
}
