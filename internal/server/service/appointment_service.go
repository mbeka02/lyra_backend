package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/payment"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
)

type appointmentService struct {
	appointmentRepo  repository.AppointmentRepository
	patientRepo      repository.PatientRepository
	doctorRepo       repository.DoctorRepository
	paymentProcessor *payment.PaymentProcessor
}
type GetAppointmentsParams struct {
	UserID   int64
	Interval int32
	Status   string
}
type AppointmentService interface {
	CreateAppointment(ctx context.Context, req model.CreateAppointmentRequest, userId int64) (database.Appointment, error)
	CreateAppointmentWithPayment(ctx context.Context, req model.CreateAppointmentRequest, userId int64, email string) (*model.InitializeTransactionResponse, error)

	GetPatientAppointments(ctx context.Context, params GetAppointmentsParams) ([]database.GetPatientAppointmentsRow, error)
	GetDoctorAppointments(ctx context.Context, params GetAppointmentsParams) ([]database.GetDoctorAppointmentsRow, error)
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository, patientRepo repository.PatientRepository, doctorRepo repository.DoctorRepository, paymentProcessor *payment.PaymentProcessor) AppointmentService {
	return &appointmentService{
		appointmentRepo,
		patientRepo,
		doctorRepo,
		paymentProcessor,
	}
}

func (s *appointmentService) GetDoctorAppointments(ctx context.Context, params GetAppointmentsParams) ([]database.GetDoctorAppointmentsRow, error) {
	doctorID, err := s.doctorRepo.GetDoctorIdByUserId(ctx, params.UserID)
	if err != nil {
		return nil, errors.New("unable to get the doctor details for this account")
	}
	return s.appointmentRepo.GetDoctorAppointments(ctx, repository.GetDoctorAppointmentsParams{
		DoctorID: doctorID,
		Status:   params.Status,
		Interval: params.Interval,
	})
}

func (s *appointmentService) GetPatientAppointments(ctx context.Context, params GetAppointmentsParams) ([]database.GetPatientAppointmentsRow, error) {
	patientID, err := s.patientRepo.GetPatientIdByUserId(ctx, params.UserID)
	if err != nil {
		return nil, errors.New("unable to get the user details of this account")
	}
	return s.appointmentRepo.GetPatientAppointments(ctx, repository.GetPatientAppointmentsParams{
		PatientID: patientID,
		Status:    params.Status,
		Interval:  params.Interval,
	})
}

func (s *appointmentService) CreateAppointment(ctx context.Context, req model.CreateAppointmentRequest, userId int64) (database.Appointment, error) {
	patientId, err := s.patientRepo.GetPatientIdByUserId(ctx, userId)
	if err != nil {
		return database.Appointment{}, errors.New("error getting the patient info linked to this account")
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
	// convert amount to a float64
	amountFloat, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		// Handle error
		return nil, fmt.Errorf("failed to parse payment amount as float64")
	}

	// Multiply by 100 for cents
	amountInCents := int64(amountFloat * 100)

	// send paystack  initialize payment request
	response, err := s.paymentProcessor.InitializeTransaction(model.InitializeTransactionRequest{
		Email:  email,
		Amount: amountInCents,
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
