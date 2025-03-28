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
	UpdateStatusCallback(ctx context.Context, reference string) error
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
		return fmt.Errorf("unable to update status for reference %s", reference)
	}
	return nil
}

func (s *paymentService) UpdateStatusCallback(ctx context.Context, reference string) error {
	verification, err := s.paymentProcessor.VerifyTransaction(reference)
	// debugging log
	log.Println("payment verification body=>", verification)
	if err != nil {
		// If verification fails, mark payment as failed.
		if repoErr := s.updateStatus(ctx, reference, "failed", "pending_payment"); repoErr != nil {
			return repoErr
		}
		return err
	}
	switch verification.Data.Status {
	case "success":
		if err := s.updateStatus(ctx, reference, "completed", "scheduled"); err != nil {
			return err
		}
	case "pending":
		if err := s.updateStatus(ctx, reference, "pending", "pending_payment"); err != nil {
			return err
		}
	default:
		if err := s.updateStatus(ctx, reference, "failed", "pending_payment"); err != nil {
			return err
		}
	}
	return nil
}

func (s *paymentService) UpdateStatusWebhook(ctx context.Context, req model.PaystackWebhookPayload) error {
	event := req.Event
	if event != "charge.success" {
		log.Println("paystack event=>", req.Event)
		return fmt.Errorf("error wrong event type for this endpoint")
	}
	return s.updateStatus(ctx, req.Data.Reference, "completed", "scheduled")
}
