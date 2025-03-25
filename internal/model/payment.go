package model

import "time"

type InitializeTransactionRequest struct {
	Amount            string   `json:"amount" validate:"required"`
	Email             string   `json:"email" validate:"required,email" `
	CallBackURL       string   `json:"callback_url,omitempty"` //  Use this to override the callback url provided on the dashboard for this transaction
	Channels          []string `json:"channels,omitempty"`     // An array of payment channels to control what channels you want to make available to the user to make a payment with.
	Currency          string   `json:"currency,omitempty"`
	Plan              string   `json:"plan,omitempty"`
	InvoiceLimit      int      `json:"invoice_limit,omitempty"`
	Metadata          string   `json:"metadata,omitempty"`
	SplitCode         string   `json:"split_code,omitempty"`
	SubAccount        string   `json:"subaccount,omitempty"`
	TransactionCharge int      `json:"transaction_charge,omitempty"`
	Bearer            string   `json:"bearer,omitempty"`
	Reference         string   `json:"reference,omitempty"`
}

// Represents the paystack API response for initializing a transaction

type InitializeTransactionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

// The request a client makes when it wants to know the payment status
type PaymentStatusRequest struct {
	Reference string `json:"reference" validate:"required"`
}

// Represents the paystack API response for fetching a specific transaction
type FetchTransactionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID              uint64 `json:"id"`
		Domain          string `json:"domain"`
		Status          string `json:"status"`
		Reference       string `json:"reference"`
		Amount          int    `json:"amount"`
		Message         string `json:"message"`
		GatewayResponse string `json:"gateway_response"`
		PaidAt          string `json:"paid_at"`
		CreatedAt       string `json:"created_at"`
		Channel         string `json:"channel"`
		Currency        string `json:"currency"`
		IPAddress       string `json:"ip_address"`
		Log             struct {
			StartTime int           `json:"start_time"`
			TimeSpent int           `json:"time_spent"`
			Attempts  int           `json:"attempts"`
			Errors    int           `json:"errors"`
			Success   bool          `json:"success"`
			Mobile    bool          `json:"mobile"`
			Input     []interface{} `json:"input"`
			History   []struct {
				Type    string `json:"type"`
				Message string `json:"message"`
				Time    int    `json:"time"`
			} `json:"history"`
		} `json:"log"`
		Fees          int         `json:"fees"`
		FeesSplit     interface{} `json:"fees_split"`
		Authorization struct {
			AuthorizationCode string `json:"authorization_code"`
			Bin               string `json:"bin"`
			Last4             string `json:"last4"`
			ExpMonth          string `json:"exp_month"`
			ExpYear           string `json:"exp_year"`
			Channel           string `json:"channel"`
			CardType          string `json:"card_type"`
			Bank              string `json:"bank"`
			CountryCode       string `json:"country_code"`
			Brand             string `json:"brand"`
			Reusable          bool   `json:"reusable"`
			Signature         string `json:"signature"`
			AccountName       string `json:"account_name"`
		} `json:"authorization"`
		Customer struct {
			ID                       uint64 `json:"id"`
			FirstName                string `json:"first_name"`
			LastName                 string `json:"last_name"`
			Email                    string `json:"email"`
			CustomerCode             string `json:"customer_code"`
			Phone                    string `json:"phone"`
			Metadata                 string `json:"metadata"`
			RiskAction               string `json:"risk_action"`
			InternationalFormatPhone string `json:"international_format_phone"`
		} `json:"customer"`
		Plan            interface{} `json:"plan"`
		Subaccount      interface{} `json:"subaccount"`
		OrderID         interface{} `json:"order_id"`
		RequestedAmount int         `json:"requested_amount"`
	} `json:"data"`
}

// Represents the paystack API response for verifying a transaction
type VerifyTransactionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID              uint64 `json:"id"`
		Domain          string `json:"domain"`
		Status          string `json:"status"`
		Reference       string `json:"reference"`
		Amount          int    `json:"amount"`
		Message         string `json:"message"`
		GatewayResponse string `json:"gateway_response"`
		Channel         string `json:"channel"`
		Currency        string `json:"currency"`
		IPAddress       string `json:"ip_address"`
		Metadata        string `json:"metadata"`
		Log             struct {
			StartTime int           `json:"start_time"`
			TimeSpent int           `json:"time_spent"`
			Attempts  int           `json:"attempts"`
			Errors    int           `json:"errors"`
			Success   bool          `json:"success"`
			Mobile    bool          `json:"mobile"`
			Input     []interface{} `json:"input"`
			History   []struct {
				Type    string `json:"type"`
				Message string `json:"message"`
				Time    int    `json:"time"`
			} `json:"history"`
		} `json:"log"`
		Fees          int `json:"fees"`
		Authorization struct {
			AuthorizationCode string `json:"authorization_code"`
			Bin               string `json:"bin"`
			Last4             string `json:"last4"`
			ExpMonth          string `json:"exp_month"`
			ExpYear           string `json:"exp_year"`
			Channel           string `json:"channel"`
			CardType          string `json:"card_type"`
			Bank              string `json:"bank"`
			CountryCode       string `json:"country_code"`
			Brand             string `json:"brand"`
			Reusable          bool   `json:"reusable"`
			Signature         string `json:"signature"`
			AccountName       string `json:"account_name"`
		} `json:"authorization"`
		Customer struct {
			ID                       uint64 `json:"id"`
			FirstName                string `json:"first_name"`
			LastName                 string `json:"last_name"`
			Email                    string `json:"email"`
			CustomerCode             string `json:"customer_code"`
			Phone                    string `json:"phone"`
			Metadata                 string `json:"metadata"`
			RiskAction               string `json:"risk_action"`
			InternationalFormatPhone string `json:"international_format_phone"`
		} `json:"customer"`
		Plan               interface{} `json:"plan"`
		Subaccount         interface{} `json:"subaccount"`
		OrderID            interface{} `json:"order_id"`
		PaidAt             string      `json:"paid_at"`
		CreatedAt          string      `json:"created_at"`
		RequestedAmount    int         `json:"requested_amount"`
		PosTransactionData interface{} `json:"pos_transaction_data"`
		Source             interface{} `json:"source"`
		FeesBreakdown      interface{} `json:"fees_breakdown"`
		TransactionDate    string      `json:"transaction_date"`
		PlanObject         interface{} `json:"plan_object"`
		Split              interface{} `json:"split"`
	} `json:"data"`
}

// represents the payload sent to the webhook endpoint
type PaystackWebhookPayload struct {
	Event string             `json:"event"`
	Data  PaystackChargeData `json:"data"`
}

type PaystackChargeData struct {
	ID              int64                  `json:"id"`
	Domain          string                 `json:"domain"`
	Status          string                 `json:"status"`
	Reference       string                 `json:"reference"`
	Amount          int64                  `json:"amount"`
	Message         interface{}            `json:"message"`
	GatewayResponse string                 `json:"gateway_response"`
	PaidAt          time.Time              `json:"paid_at"`
	CreatedAt       time.Time              `json:"created_at"`
	Channel         string                 `json:"channel"`
	Currency        string                 `json:"currency"`
	IPAddress       string                 `json:"ip_address"`
	Metadata        interface{}            `json:"metadata"`
	Log             PaystackLog            `json:"log"`
	Fees            interface{}            `json:"fees"`
	Customer        Customer               `json:"customer"`
	Authorization   Authorization          `json:"authorization"`
	Plan            map[string]interface{} `json:"plan"`
}

type PaystackLog struct {
	TimeSpent      int           `json:"time_spent"`
	Attempts       int           `json:"attempts"`
	Authentication string        `json:"authentication"`
	Errors         int           `json:"errors"`
	Success        bool          `json:"success"`
	Mobile         bool          `json:"mobile"`
	Input          []interface{} `json:"input"`
	Channel        interface{}   `json:"channel"`
	History        []LogHistory  `json:"history"`
}

type LogHistory struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Time    int    `json:"time"`
}

type Customer struct {
	ID           int64       `json:"id"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Email        string      `json:"email"`
	CustomerCode string      `json:"customer_code"`
	Phone        interface{} `json:"phone"`
	Metadata     interface{} `json:"metadata"`
	RiskAction   string      `json:"risk_action"`
}

type Authorization struct {
	AuthorizationCode string `json:"authorization_code"`
	Bin               string `json:"bin"`
	Last4             string `json:"last4"`
	ExpMonth          string `json:"exp_month"`
	ExpYear           string `json:"exp_year"`
	CardType          string `json:"card_type"`
	Bank              string `json:"bank"`
	CountryCode       string `json:"country_code"`
	Brand             string `json:"brand"`
	AccountName       string `json:"account_name"`
}
