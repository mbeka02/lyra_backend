package service

import (
	"context"
	"log"

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
}

func NewPatientService(patientRepo repository.PatientRepository, fhirClient *fhir.FHIRClient, fileStorage objstore.Storage) PatientService {
	return &patientService{
		patientRepo,
		fhirClient,
		fileStorage,
	}
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
	// Save resource to the API
	_, err = s.fhirClient.UpsertPatient(ctx, fhirPatient)
	if err != nil {
		return nil, err
	}
	log.Printf("fhir patient:%v", fhirPatient)
	// update the version ID
	// err = s.patientRepo.UpdateFHIRVersion(ctx, txResult.Patient.PatientID, *fhirPatient.Meta.VersionId)
	// if err != nil {
	// 	return nil, err
	// }
	return &txResult.Patient, nil
}
