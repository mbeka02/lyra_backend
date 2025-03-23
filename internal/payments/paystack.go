package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mbeka02/lyra_backend/internal/model"
)

var baseURL = "https://api.paystack.co"

type PaymentProcessor struct {
	apiKey string
	client *http.Client
}

func NewPaymentProcessor(apiKey string) *PaymentProcessor {
	return &PaymentProcessor{apiKey, &http.Client{}}
}

func (p *PaymentProcessor) FetchTransaction(transactionId uint64) (*model.FetchTransactionResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/transaction/%v", baseURL, transactionId), nil)
	// Add content type header
	req.Header.Add("Content-Type", "application/json")
	// Add Authorization Header
	req.Header.Add("Authorization", "Bearer "+p.apiKey)
	// send request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respBody model.FetchTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	if !respBody.Status {
		return nil, fmt.Errorf("paystack error:%s", respBody.Message)
	}
	return &respBody, nil
}

func (p *PaymentProcessor) VerifyTransaction(reference string) (*model.VerifyTransactionResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/transaction/verify/%s", baseURL, reference), nil)
	// Add content type header
	req.Header.Add("Content-Type", "application/json")
	// Add Authorization Header
	req.Header.Add("Authorization", "Bearer "+p.apiKey)
	// send request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respBody model.VerifyTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	if !respBody.Status {
		return nil, fmt.Errorf("paystack error:%s", respBody.Message)
	}
	return &respBody, nil
}

func (p *PaymentProcessor) InitializeTransaction(request model.InitializeTransactionRequest) (*model.InitializeTransactionResponse, error) {
	buff, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/transaction/initialize", baseURL), bytes.NewBuffer(buff))
	// Add content type header
	req.Header.Add("Content-Type", "application/json")
	// Add Authorization Header
	req.Header.Add("Authorization", "Bearer "+p.apiKey)
	// send request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respBody model.InitializeTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	if !respBody.Status {
		return nil, fmt.Errorf("paystack error:%s", respBody.Message)
	}
	return &respBody, nil
}
