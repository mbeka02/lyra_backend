package handler

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/service"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service}
}

// this is the callback endpoint that paystack will use
func (h *PaymentHandler) HandlePaymentCallback(w http.ResponseWriter, r *http.Request) {
}

// this is the webhook endpoint that paystack will use
func (h *PaymentHandler) PaymentWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("error , unable to read request body:%v", err))
		return
	}
	defer r.Body.Close()
	secretKey := os.Getenv("PAYSTACK_API_KEY")
	// create HMAC hash
	mac := hmac.New(sha512.New, []byte(secretKey))
	mac.Write(body)
	expectedHash := hex.EncodeToString(mac.Sum(nil))
	signature := r.Header.Get("x-paystack-signature")
	if expectedHash != signature {
		respondWithError(w, http.StatusUnauthorized, fmt.Errorf("invalid signature"))
		return
	}
	request := model.PaystackWebhookPayload{}
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	if err = h.paymentService.UpdateStatusWebhook(r.Context(), request); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	var nilValue interface{}
	if err := respondWithJSON(w, http.StatusOK, nilValue); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}
