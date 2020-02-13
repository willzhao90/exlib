package payment

import (
	"github.com/mitchellh/mapstructure"
)

type InitTransactionRequest struct {
	Amount              string `json:"Amount"`
	CurrencyCode        string `json:"CurrencyCode"`      // AUD
	MerchantReference   string `json:"MerchantReference"` // lets use deposit id
	MerchantHomepageURL string `json:"MerchantHomepageURL"`
	SuccessURL          string `json:"SuccessURL"`
	NotificationURL     string `json:"NotificationURL"`
	CancellationURL     string `json:"CancellationURL"`
	FailureURL          string `json:"FailureURL"`
}

type InitTransactionResponse struct {
	TransactionRefNo string `json:"TransactionRefNo"`
	NavigateURL      string `json:"NavigateURL"`
	ErrorCode        int    `json:"ErrorCode"`
	ErrorMessage     string `json:"ErrorMessage"`
}

type GetTransactionResponse struct {
	TransactionRefNo         string  `json:"TransactionRefNo,omitempty"`
	CurrencyCode             string  `json:"CurrencyCode,omitempty"`
	TransactionStatusCode    string  `json:"TransactionStatusCode,omitempty"`
	ErrorCode                string  `json:"ErrorCode,omitempty"`
	ErrorMessage             string  `json:"ErrorMessage,omitempty"`
	StartDateTime            string  `json:"StartDateTime,omitempty"`
	EndDateTime              string  `json:"EndDateTime,omitempty"`
	CountryCode              string  `json:"CountryCode,omitempty"`
	PaymentAmount            float64 `json:"PaymentAmount,omitempty"`
	AmountPaid               float64 `json:"AmountPaid,omitempty"`
	BankReceipt              string  `json:"BankReceipt,omitempty"`
	BankReceiptDateTime      string  `json:"BankReceiptDateTime,omitempty"`
	FinancialInstitutionCode string  `json:"FinancialInstitutionCode,omitempty"`
	MerchantReference        string  `json:"MerchantReference,omitempty"`
	PayerFirstName           string  `json:"PayerFirstName,omitempty"`
	PayerFamilyName          string  `json:"PayerFamilyName,omitempty"`
	PayerAccountSortCode     string  `json:"PayerAccountSortCode,omitempty"`
	PayerAccountNumber       string  `json:"PayerAccountNumber,omitempty"`
	PayerAccountSuffix       string  `json:"PayerAccountSuffix,omitempty"`
}

func InitTxResponseFromMap(values map[string]interface{}) (*InitTransactionResponse, error) {
	var trans InitTransactionResponse
	err := mapstructure.Decode(values, &trans)
	if err != nil {
		return nil, nil
	}
	return &trans, nil
}

func GetTransactionResponseFromMap(values map[string]interface{}) (*GetTransactionResponse, error) {
	var trans GetTransactionResponse
	err := mapstructure.Decode(values, &trans)
	if err != nil {
		return nil, err
	}
	return &trans, nil
}
