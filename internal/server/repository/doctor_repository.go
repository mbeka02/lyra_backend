package repository

import (
	"context"

	"github.com/mbeka02/lyra_backend/internal/database"
)

type CreateDoctorParams struct {
	Specialization    string
	LicenseNumber     string
	Description       string
	County            string
	PricePerHour      string
	YearsOfExperience int32
	UserID            int64
}
type GetDoctorsParams struct {
	Offset         int32
	Limit          int32
	Specialization string
	MinPrice       string
	MaxPrice       string
	MinExperience  int32
	MaxExperience  int32

	County    string // Optional county filter
	SortBy    string // Sorting field (price, experience)
	SortOrder string // Sorting order (asc, desc)
}

type DoctorRepository interface {
	Create(ctx context.Context, params CreateDoctorParams) (*database.Doctor, error)
	GetAllDoctors(ctx context.Context, params GetDoctorsParams) ([]database.GetDoctorsRow, error)
	GetDoctorIdByUserId(ctx context.Context, userId int64) (int64, error)
	ListPatientsUnderCare(ctx context.Context, doctorId int64) ([]database.ListPatientsUnderDoctorCareRow, error)
}

type doctorRepository struct {
	store *database.Store
}

func NewDoctorRepository(store *database.Store) DoctorRepository {
	return &doctorRepository{
		store,
	}
}

func (r *doctorRepository) ListPatientsUnderCare(ctx context.Context, doctorID int64) ([]database.ListPatientsUnderDoctorCareRow, error) {
	return r.store.ListPatientsUnderDoctorCare(ctx, doctorID)
}

func (r *doctorRepository) Create(ctx context.Context, params CreateDoctorParams) (*database.Doctor, error) {
	var doctor database.Doctor
	err := r.store.ExecTx(ctx, func(q *database.Queries) error {
		var err error
		// create record
		doctor, err = q.CreateDoctor(ctx, database.CreateDoctorParams{
			UserID:            params.UserID,
			LicenseNumber:     params.LicenseNumber,
			Specialization:    params.Specialization,
			Description:       params.Description,
			County:            params.County,
			YearsOfExperience: params.YearsOfExperience,
			PricePerHour:      params.PricePerHour,
		})
		if err != nil {
			return err
		}
		// mark onboarding as completed
		err = q.CompleteOnboarding(ctx, params.UserID)
		return err
	})
	return &doctor, err
}

func (r *doctorRepository) GetDoctorIdByUserId(ctx context.Context, UserID int64) (int64, error) {
	return r.store.GetDoctorIdByUserId(ctx, UserID)
}

func (r *doctorRepository) GetAllDoctors(ctx context.Context, params GetDoctorsParams) ([]database.GetDoctorsRow, error) {
	return r.store.GetDoctors(ctx, database.GetDoctorsParams{
		// filters
		SetCounty:         params.County,
		SetSpecialization: params.Specialization,
		SetMinPrice:       params.MinPrice,
		SetMaxPrice:       params.MaxPrice,
		SetMinExperience:  params.MinExperience,
		SetMaxExperience:  params.MaxExperience,
		// ordering
		SetSortBy:    params.SortBy,
		SetSortOrder: params.SortOrder,
		// pagination
		SetLimit:  params.Limit,
		SetOffset: params.Offset,
	})
}
