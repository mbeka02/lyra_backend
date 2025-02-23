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
	Offset    int32
	Limit     int32
	County    string // Optional county filter
	SortBy    string // Sorting field (price, experience)
	SortOrder string // Sorting order (asc, desc)
}
type DoctorRepository interface {
	Create(context.Context, CreateDoctorParams) (database.Doctor, error)
	GetAll(context.Context, GetDoctorsParams) ([]database.GetDoctorsRow, error)
}

type doctorRepository struct {
	store *database.Store
}

func NewDoctorRepository(store *database.Store) DoctorRepository {
	return &doctorRepository{
		store,
	}
}

func (s *doctorRepository) Create(ctx context.Context, params CreateDoctorParams) (database.Doctor, error) {
	return s.store.CreateDoctor(ctx, database.CreateDoctorParams{
		UserID:            params.UserID,
		LicenseNumber:     params.LicenseNumber,
		Specialization:    params.Specialization,
		Description:       params.Description,
		County:            params.County,
		YearsOfExperience: params.YearsOfExperience,
		PricePerHour:      params.PricePerHour,
	})
}

func (s *doctorRepository) GetAll(ctx context.Context, params GetDoctorsParams) ([]database.GetDoctorsRow, error) {
	return s.store.GetDoctors(ctx, database.GetDoctorsParams{
		Column1: params.County,
		Column2: params.SortBy,
		Column3: params.SortOrder,
		Limit:   params.Limit,
		Offset:  params.Offset,
	})
}
