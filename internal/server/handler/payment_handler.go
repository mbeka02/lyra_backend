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

type paymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) *paymentHandler {
	return &paymentHandler{service}
}

// this is the callback endpoint that paystack will use
func (h *paymentHandler) HandlePaymentCallback(w http.ResponseWriter, r *http.Request) {
}

// this is the webhook endpoint that paystack will use
func (h *paymentHandler) PaymentWebhook(w http.ResponseWriter, r *http.Request) {
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
	// DO STUFF WITH THE REQUEST
}

// this endpoint is used for polling the payment status (fallback route incase the callback or webhook don't work)
func (h *paymentHandler) HandlePaymentStatusPolling(w http.ResponseWriter, r *http.Request) {
	request := model.PaymentStatusRequest{}
	if err := parseAndValidateRequest(r, &request); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	resp, err := h.paymentService.VerifyPayment(request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("errror, unable to verify payment:%v", err))
		return
	}
	if err := respondWithJSON(w, http.StatusOK, resp); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}
