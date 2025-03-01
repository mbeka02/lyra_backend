package model

type CreateAvailabilityRequest struct {
	DayOfWeek       int32  `json:"day_of_week" default:"0" `
	StartTime       string `json:"start_time" validate:"required" `
	EndTime         string `json:"end_time"  validate:"required"`
	IntervalMinutes int32  `json:"interval_minutes"`

	// IsRecurring bool      `json:"is_recurring"`
}
