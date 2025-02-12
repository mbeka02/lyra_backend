package model

type CreateSpecialistRequest struct {
	Specialization string `json:"specialization" validate:"required"`
	LicenseNumber  string `json:"license_number" validate:"required"`
}
