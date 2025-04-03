package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/payment"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
)

type appointmentService struct {
	appointmentRepo  repository.AppointmentRepository
	patientRepo      repository.PatientRepository
	paymentProcessor *payment.PaymentProcessor
}

type AppointmentService interface {
	CreateAppointment(ctx context.Context, req model.CreateAppointmentRequest, userId int64) (database.Appointment, error)
	CreateAppointmentWithPayment(ctx context.Context, req model.CreateAppointmentRequest, userId int64, email string) (*model.InitializeTransactionResponse, error)

	GetPatientAppointments(ctx context.Context, userId int64) ([]database.Appointment, error)
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository, patientRepo repository.PatientRepository, paymentProcessor *payment.PaymentProcessor) AppointmentService {
	return &appointmentService{
		appointmentRepo,
		patientRepo,
		paymentProcessor,
	}
}

func (s *appointmentService) GetPatientAppointments(ctx context.Context, userId int64) ([]database.Appointment, error) {
	patientId, err := s.patientRepo.GetPatientIdByUserId(ctx, userId)
	if err != nil {
		return nil, errors.New("unable to get the user details of this account")
	}
	return s.appointmentRepo.GetPatientAppointments(ctx, patientId)
}

func (s *appointmentService) CreateAppointment(ctx context.Context, req model.CreateAppointmentRequest, userId int64) (database.Appointment, error) {
	patientId, err := s.patientRepo.GetPatientIdByUserId(ctx, userId)
	if err != nil {
		return database.Appointment{}, errors.New("unable to get the user details of this account")
	}

	return s.appointmentRepo.Create(ctx, repository.CreateAppointmentParams{
		DoctorID:  req.DoctorID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Reason:    req.Reason,
	}, patientId)
}

// TODO: CLEAN THIS UP
func (s *appointmentService) CreateAppointmentWithPayment(ctx context.Context, req model.CreateAppointmentRequest, userId int64, email string) (*model.InitializeTransactionResponse, error) {
	patientId, err := s.patientRepo.GetPatientIdByUserId(ctx, userId)
	if err != nil {
		return nil, errors.New("unable to get the user details of this account")
	}
	// send paystack  initialize payment request
	response, err := s.paymentProcessor.InitializeTransaction(model.InitializeTransactionRequest{
		Email:  email,
		Amount: req.Amount,
	})
	if err != nil {
		return nil, fmt.Errorf("payment processing error:%v", err)
	}
	// add records to db (transaction)
	_, err = s.appointmentRepo.CreateAppointmentWithPayment(ctx, repository.CreateAppointmentWithPaymentParams{
		// appointment details
		PatientID: patientId,
		DoctorID:  req.DoctorID,
		Reason:    req.Reason,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		// payment details
		Reference: response.Data.Reference,
		Amount:    req.Amount,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}
