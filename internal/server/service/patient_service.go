package service

import (
	"context"
	"fmt"

	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/fhir"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/objstore"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
)

type patientService struct {
	patientRepo repository.PatientRepository
	fhirClient  *fhir.FHIRClient
	fileStorage objstore.Storage
}

type PatientService interface {
	CreatePatient(ctx context.Context, req model.CreatePatientRequest, userId int64) (*database.Patient, error)
	GetPatientIdByUserId(ctx context.Context, userId int64) (int64, error)
}

func NewPatientService(patientRepo repository.PatientRepository, fhirClient *fhir.FHIRClient, fileStorage objstore.Storage) PatientService {
	return &patientService{
		patientRepo,
		fhirClient,
		fileStorage,
	}
}

func (s *patientService) GetPatientIdByUserId(ctx context.Context, userId int64) (int64, error) {
	return s.patientRepo.GetPatientIdByUserId(ctx, userId)
}

func (s *patientService) CreatePatient(ctx context.Context, req model.CreatePatientRequest, userId int64) (*database.Patient, error) {
	// create patient in the DB
	txResult, err := s.patientRepo.Create(ctx, repository.CreatePatientParams{
		UserID:                userId,
		Allergies:             req.Allergies,
		CurrentMedication:     req.CurrentMedication,
		PastMedicalHistory:    req.PastMedicalHistory,
		FamilyMedicalHistory:  req.FamilyMedicalHistory,
		InsurancePolicyNumber: req.InsurancePolicyNumber,
		InsuranceProvider:     req.InsuranceProvider,
		Address:               req.Address,
		EmergencyContactName:  req.EmergencyContactName,
		EmergencyContactPhone: req.EmergencyContactPhone,
	})
	if err != nil {
		return nil, err
	}
	// Build the FHIR Resource
	fhirPatient, err := fhir.BuildFHIRPatientFromDB(&txResult.Patient, &txResult.User)
	if err != nil {
		return nil, err
	}
	// Save resource to the  Google HealthCare API
	savedFhirPatient, err := s.fhirClient.UpsertPatient(ctx, fhirPatient)
	if err != nil {
		return nil, err
	}
	// update the version ID
	if savedFhirPatient.Meta.VersionId == nil {
		return nil, fmt.Errorf("patient resource version id error")
	}
	err = s.patientRepo.UpdateFHIRVersion(ctx, txResult.Patient.PatientID, *savedFhirPatient.Meta.VersionId)
	if err != nil {
		return nil, err
	}
	return &txResult.Patient, nil
}
