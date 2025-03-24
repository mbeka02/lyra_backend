package service

import (
	"context"

	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/payment"
)

type PaymentService interface {
	VerifyPayment(ctx context.Context, req model.PaymentStatusRequest) (*model.VerifyTransactionResponse, error)
}

type paymentService struct {
	paymentProcessor *payment.PaymentProcessor
}

func NewPaymentService(paymentProcessor *payment.PaymentProcessor) PaymentService {
	return &paymentService{paymentProcessor}
}

func (s *paymentService) VerifyPayment(ctx context.Context, req model.PaymentStatusRequest) (*model.VerifyTransactionResponse, error) {
	return s.paymentProcessor.VerifyTransaction(req.Reference)
	// TODO: ADD REPO LAYER FOR UPDATING THE DB
}
