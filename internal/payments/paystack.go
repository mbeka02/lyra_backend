package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type paymentsProcessor struct {
	apiKey string
	client *http.Client
}

type InitializePaymentRequest struct {
	Amount float64 `json:"amount"`
	Email  string  `json:"email"`
}

type InitializePaymentResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

func NewPaymentsProcessor(apiKey string) *paymentsProcessor {
	return &paymentsProcessor{apiKey, &http.Client{}}
}

func (p *paymentsProcessor) InitializePayment(request InitializePaymentRequest) (InitializePaymentResponse, error) {
	buff, err := json.Marshal(request)
	if err != nil {
		return InitializePaymentResponse{}, err
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.paystack.co/transaction/initialize", bytes.NewBuffer(buff))
	// Add content type header
	req.Header.Add("Content-Type", "application/json")
	// Add Authorization Header
	req.Header.Add("Authorization", "Bearer "+p.apiKey)
	// send request
	res, err := p.client.Do(req)
	if err != nil {
		return InitializePaymentResponse{}, err
	}
	defer res.Body.Close()

	var body InitializePaymentResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return InitializePaymentResponse{}, err
	}
	if !body.Status {
		return InitializePaymentResponse{}, fmt.Errorf("paystack error:%s", body.Message)
	}
	return body, nil
}
