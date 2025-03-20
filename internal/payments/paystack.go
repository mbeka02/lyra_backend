package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mbeka02/lyra_backend/internal/model"
)

var baseURL = "https://api.paystack.co/transaction/initialize"

type paymentsProcessor struct {
	apiKey string
	client *http.Client
}

func NewPaymentsProcessor(apiKey string) *paymentsProcessor {
	return &paymentsProcessor{apiKey, &http.Client{}}
}

func (p *paymentsProcessor) VerifyTransaction() (model.VerifyTransactionResponse, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	// Add content type header
	req.Header.Add("Content-Type", "application/json")
	// Add Authorization Header
	req.Header.Add("Authorization", "Bearer "+p.apiKey)
	// send request
	resp, err := p.client.Do(req)
	if err != nil {
		return model.VerifyTransactionResponse{}, err
	}
	defer resp.Body.Close()

	var respBody model.VerifyTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return model.VerifyTransactionResponse{}, err
	}
	if !respBody.Status {
		return model.VerifyTransactionResponse{}, fmt.Errorf("paystack error:%s", respBody.Message)
	}
	return respBody, nil
}

func (p *paymentsProcessor) InitializeTransaction(request model.InitializeTransactionRequest) (model.InitializeTransactionResponse, error) {
	buff, err := json.Marshal(request)
	if err != nil {
		return model.InitializeTransactionResponse{}, err
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.paystack.co/transaction/initialize", bytes.NewBuffer(buff))
	// Add content type header
	req.Header.Add("Content-Type", "application/json")
	// Add Authorization Header
	req.Header.Add("Authorization", "Bearer "+p.apiKey)
	// send request
	resp, err := p.client.Do(req)
	if err != nil {
		return model.InitializeTransactionResponse{}, err
	}
	defer resp.Body.Close()

	var respBody model.InitializeTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return model.InitializeTransactionResponse{}, err
	}
	if !respBody.Status {
		return model.InitializeTransactionResponse{}, fmt.Errorf("paystack error:%s", respBody.Message)
	}
	return respBody, nil
}
