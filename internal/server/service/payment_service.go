package service

import (
	"context"
	"fmt"
	"log"

	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/payment"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
)

type PaymentService interface {
	VerifyPayment(ctx context.Context, req model.PaymentStatusRequest) (*model.VerifyTransactionResponse, error)
	UpdateStatus(ctx context.Context, req model.PaystackWebhookPayload) error
}

type paymentService struct {
	paymentProcessor *payment.PaymentProcessor
	paymentRepo      repository.PaymentRepository
}

func NewPaymentService(paymentProcessor *payment.PaymentProcessor, repo repository.PaymentRepository) PaymentService {
	return &paymentService{paymentProcessor, repo}
}

func (s *paymentService) VerifyPayment(ctx context.Context, req model.PaymentStatusRequest) (*model.VerifyTransactionResponse, error) {
	return s.paymentProcessor.VerifyTransaction(req.Reference)
	// TODO: ADD REPO LAYER FOR UPDATING THE DB
}

func (s *paymentService) UpdateStatus(ctx context.Context, req model.PaystackWebhookPayload) error {
	var (
		paymentStatus     string
		appointmentStatus string
	)
	// refactor this , you need to handle other events properly
	if req.Event != "transaction.success" {
		log.Println("paystack event=>", req.Event)
		return fmt.Errorf("error wrong event type for this endpoint")
	}
	paymentStatus = "completed"
	appointmentStatus = "scheduled"
	return s.paymentRepo.UpdatePaymentAndAppointmentStatus(ctx, repository.UpdatePaymentAndAppointmentStatusParams{
		Reference:         req.Data.Reference,
		PaymentStatus:     paymentStatus,
		AppointmentStatus: appointmentStatus,
	})
}
