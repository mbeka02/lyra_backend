package model

import "time"

type CreateAppointmentRequest struct {
	DoctorID  int64     `json:"doctor_id" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required" `
	Reason    string    `json:"reason" validate:"required"`
}
