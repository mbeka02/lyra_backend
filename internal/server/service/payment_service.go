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
	UpdateStatusWebhook(ctx context.Context, req model.PaystackWebhookPayload) error
	UpdateStatusCallback(ctx context.Context, reference string) (currentStatus string, err error)
}

type paymentService struct {
	paymentProcessor *payment.PaymentProcessor
	paymentRepo      repository.PaymentRepository
}

func NewPaymentService(paymentProcessor *payment.PaymentProcessor, repo repository.PaymentRepository) PaymentService {
	return &paymentService{paymentProcessor, repo}
}

// updateStatus is a helper to update both payment and appointment statuses.
func (s *paymentService) updateStatus(ctx context.Context, reference, paymentStatus, appointmentStatus string) error {
	if err := s.paymentRepo.UpdatePaymentAndAppointmentStatus(ctx, repository.UpdatePaymentAndAppointmentStatusParams{
		Reference:         reference,
		PaymentStatus:     paymentStatus,
		AppointmentStatus: appointmentStatus,
	}); err != nil {
		log.Printf("Error updating status for reference %s: %v", reference, err)
		return fmt.Errorf("unable to update status for reference %s: %w", reference, err)
	}
	return nil
}

func (s *paymentService) UpdateStatusCallback(ctx context.Context, reference string) (string, error) {
	verification, err := s.paymentProcessor.VerifyTransaction(reference)
	if err != nil {
		// If verification fails, mark payment as failed.
		if repoErr := s.updateStatus(ctx, reference, "failed", "pending_payment"); repoErr != nil {
			return "failed", repoErr
		}
		return "failed", err
	}
	var paymentStatus string
	switch verification.Data.Status {
	case "success":
		paymentStatus = "completed"
		if err := s.updateStatus(ctx, reference, paymentStatus, "scheduled"); err != nil {
			return paymentStatus, err
		}
	case "pending":
		paymentStatus = "pending"
		if err := s.updateStatus(ctx, reference, paymentStatus, "pending_payment"); err != nil {
			return paymentStatus, err
		}
	default:
		paymentStatus = "failed"
		if err := s.updateStatus(ctx, reference, paymentStatus, "pending_payment"); err != nil {
			return paymentStatus, err
		}
	}
	return paymentStatus, nil
}

func (s *paymentService) UpdateStatusWebhook(ctx context.Context, req model.PaystackWebhookPayload) error {
	event := req.Event
	if event != "charge.success" {
		log.Printf("Received unsupported Paystack event: %s", req.Event)
		return fmt.Errorf("unsupported event type: %s", req.Event)
	}
	return s.updateStatus(ctx, req.Data.Reference, "completed", "scheduled")
}
