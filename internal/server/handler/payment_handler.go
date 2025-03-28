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
func (h *PaymentHandler) PaymentCallback(w http.ResponseWriter, r *http.Request) {
	params := NewQueryParamExtractor(r)
	reference := params.GetString("reference")
	if reference == "" {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("missing payment reference"))
		return
	}
	status, err := h.paymentService.UpdateStatusCallback(r.Context(), reference)
	if err != nil {
		http.Redirect(w, r, "lyra://payment?status=failed&reference="+reference, http.StatusSeeOther)
		return
	}

	// Redirect the user back to lyra  via deep link.
	// Determine the redirect URL based on status.
	var redirectURL string
	switch status {
	case "completed":
		redirectURL = "lyra://payment?status=success&reference=" + reference
	case "pending":
		redirectURL = "lyra://payment?status=pending&reference=" + reference
	default:
		redirectURL = "lyra://payment?status=failed&reference=" + reference
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
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
