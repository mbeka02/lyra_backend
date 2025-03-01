package model

import "time"

type CreateAvailabilityRequest struct {
	DayOfWeek int32     `json:"day_of_week" validate:"required" `
	StartTime time.Time `json:"start_time" validate:"required" `
	EndTime   time.Time `json:"end_time"  validate:"required"`
	// IsRecurring bool      `json:"is_recurring"`
}
