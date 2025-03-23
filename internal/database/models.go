// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/sqlc-dev/pqtype"
)

type AppointmentStatus string

const (
	AppointmentStatusPendingPayment AppointmentStatus = "pending_payment"
	AppointmentStatusScheduled      AppointmentStatus = "scheduled"
	AppointmentStatusCompleted      AppointmentStatus = "completed"
	AppointmentStatusCanceled       AppointmentStatus = "canceled"
)

func (e *AppointmentStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AppointmentStatus(s)
	case string:
		*e = AppointmentStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for AppointmentStatus: %T", src)
	}
	return nil
}

type NullAppointmentStatus struct {
	AppointmentStatus AppointmentStatus `json:"appointment_status"`
	Valid             bool              `json:"valid"` // Valid is true if AppointmentStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAppointmentStatus) Scan(value interface{}) error {
	if value == nil {
		ns.AppointmentStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AppointmentStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAppointmentStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AppointmentStatus), nil
}

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
)

func (e *PaymentStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentStatus(s)
	case string:
		*e = PaymentStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentStatus: %T", src)
	}
	return nil
}

type NullPaymentStatus struct {
	PaymentStatus PaymentStatus `json:"payment_status"`
	Valid         bool          `json:"valid"` // Valid is true if PaymentStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentStatus), nil
}

type Role string

const (
	RolePatient    Role = "patient"
	RoleSpecialist Role = "specialist"
)

func (e *Role) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Role(s)
	case string:
		*e = Role(s)
	default:
		return fmt.Errorf("unsupported scan type for Role: %T", src)
	}
	return nil
}

type NullRole struct {
	Role  Role `json:"role"`
	Valid bool `json:"valid"` // Valid is true if Role is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRole) Scan(value interface{}) error {
	if value == nil {
		ns.Role, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Role.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Role), nil
}

type Appointment struct {
	AppointmentID int64             `json:"appointment_id"`
	PatientID     int64             `json:"patient_id"`
	DoctorID      int64             `json:"doctor_id"`
	CurrentStatus AppointmentStatus `json:"current_status"`
	Reason        string            `json:"reason"`
	Notes         sql.NullString    `json:"notes"`
	StartTime     time.Time         `json:"start_time"`
	EndTime       time.Time         `json:"end_time"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     sql.NullTime      `json:"updated_at"`
}

type Availability struct {
	AvailabilityID  int64        `json:"availability_id"`
	DoctorID        int64        `json:"doctor_id"`
	StartTime       string       `json:"start_time"`
	EndTime         string       `json:"end_time"`
	IsRecurring     bool         `json:"is_recurring"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       sql.NullTime `json:"updated_at"`
	DayOfWeek       int32        `json:"day_of_week"`
	IntervalMinutes int32        `json:"interval_minutes"`
}

type Doctor struct {
	DoctorID          int64        `json:"doctor_id"`
	UserID            int64        `json:"user_id"`
	Description       string       `json:"description"`
	Specialization    string       `json:"specialization"`
	YearsOfExperience int32        `json:"years_of_experience"`
	County            string       `json:"county"`
	PricePerHour      string       `json:"price_per_hour"`
	LicenseNumber     string       `json:"license_number"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         sql.NullTime `json:"updated_at"`
}

type Patient struct {
	PatientID             int64        `json:"patient_id"`
	UserID                int64        `json:"user_id"`
	Address               string       `json:"address"`
	EmergencyContactName  string       `json:"emergency_contact_name"`
	EmergencyContactPhone string       `json:"emergency_contact_phone"`
	Allergies             string       `json:"allergies"`
	CurrentMedication     string       `json:"current_medication"`
	PastMedicalHistory    string       `json:"past_medical_history"`
	FamilyMedicalHistory  string       `json:"family_medical_history"`
	InsuranceProvider     string       `json:"insurance_provider"`
	InsurancePolicyNumber string       `json:"insurance_policy_number"`
	CreatedAt             time.Time    `json:"created_at"`
	UpdatedAt             sql.NullTime `json:"updated_at"`
}

type Payment struct {
	PaymentID     int64                 `json:"payment_id"`
	Reference     string                `json:"reference"`
	CurrentStatus PaymentStatus         `json:"current_status"`
	Amount        string                `json:"amount"`
	Metadata      pqtype.NullRawMessage `json:"metadata"`
	PaymentMethod sql.NullString        `json:"payment_method"`
	Currency      string                `json:"currency"`
	AppointmentID int64                 `json:"appointment_id"`
	PatientID     int64                 `json:"patient_id"`
	DoctorID      int64                 `json:"doctor_id"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     sql.NullTime          `json:"updated_at"`
	CompletedAt   sql.NullTime          `json:"completed_at"`
}

type User struct {
	UserID            int64        `json:"user_id"`
	DateOfBirth       time.Time    `json:"date_of_birth"`
	FullName          string       `json:"full_name"`
	Password          string       `json:"password"`
	Email             string       `json:"email"`
	TelephoneNumber   string       `json:"telephone_number"`
	ProfileImageUrl   string       `json:"profile_image_url"`
	CreatedAt         time.Time    `json:"created_at"`
	UserRole          Role         `json:"user_role"`
	VerifiedAt        sql.NullTime `json:"verified_at"`
	PasswordChangedAt time.Time    `json:"password_changed_at"`
}
