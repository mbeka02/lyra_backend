package service

import (
	"context"
	"log"

	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
)

type DoctorService interface {
	CreateDoctor(ctx context.Context, req model.CreateDoctorRequest, userId int64) (*database.Doctor, error)
	GetDoctors(ctx context.Context, county, specialization, minPrice, maxPrice, sortBy, sortOrder string, minExperience, maxExpreinece, limit, offset int32) (model.GetDoctorsResponse, error)
	GetDoctorIdByUserId(ctx context.Context, userId int64) (int64, error)
}
type doctorService struct {
	doctorRepo repository.DoctorRepository
}

func NewDoctorService(doctorRepo repository.DoctorRepository) DoctorService {
	return &doctorService{
		doctorRepo,
	}
}

func (s *doctorService) CreateDoctor(ctx context.Context, req model.CreateDoctorRequest, userId int64) (*database.Doctor, error) {
	return s.doctorRepo.Create(ctx, repository.CreateDoctorParams{
		Specialization:    req.Specialization,
		LicenseNumber:     req.LicenseNumber,
		Description:       req.Description,
		County:            req.County,
		PricePerHour:      req.PricePerHour,
		YearsOfExperience: req.YearsOfExperience,
		UserID:            userId,
	})
}

func (s *doctorService) GetDoctorIdByUserId(ctx context.Context, userId int64) (int64, error) {
	return s.doctorRepo.GetDoctorIdByUserId(ctx, userId)
}

func (s *doctorService) GetDoctors(ctx context.Context, county, specialization, minPrice, maxPrice, sortBy, sortOrder string, minExperience, maxExpreinece, limit, offset int32) (model.GetDoctorsResponse, error) {
	rows, err := s.doctorRepo.GetAllDoctors(ctx, repository.GetDoctorsParams{
		// Fetch the limit+1 to determine if there's more data
		Limit:          limit + 1,
		Offset:         offset,
		County:         county,
		Specialization: specialization,
		MinPrice:       minPrice,
		MaxPrice:       maxPrice,
		MinExperience:  minExperience,
		MaxExperience:  maxExpreinece,
		SortBy:         sortBy,
		SortOrder:      sortOrder,
	})
	if err != nil {
		log.Println(err)
		return model.GetDoctorsResponse{}, err
	}
	hasMore := false
	// handle a situation where there's still more data
	if len(rows) > int(limit) {
		hasMore = true
		rows = rows[:limit] // remove the extra row
	}
	return model.GetDoctorsResponse{
		Doctors: model.NewDoctorDetails(rows),
		HasMore: hasMore,
	}, nil
}
