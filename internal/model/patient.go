package model

import "time"

type CreatePatientRequest struct {
	Allergies   string    `json:"allergies"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
}
