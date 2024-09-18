package provider

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
)

var (
	terminalID                string
	bankUrl                   string
	banksendTokenUrl          string
	bankVerifyTransactionUrl  string
	bankReverseTransactionUrl string
)

type SamanGateway struct {
	client *http.Client
}

func NewSamanGateway() *SamanGateway {
	terminalID = os.Getenv("BANK_TERMINAL_ID")
	bankUrl = os.Getenv("BANK_URL")
	banksendTokenUrl = os.Getenv("BANK_SEND_TOKEN_URL")
	bankVerifyTransactionUrl = os.Getenv("BANK_VERIFY_TRANSACTION_URL")
	bankReverseTransactionUrl = os.Getenv("BANK_REVERSE_TRANSACTION_URL")

	tr := &http.Transport{}
	cl := &http.Client{
		Transport: tr,
		Timeout:   httpTimeout,
	}
	return &SamanGateway{
		client: cl,
	}
}

func (m *SamanGateway) CreatePayment(tickets *models.Tickets, phoneNumber string) (models.Response, error) {
	var response models.Response
	reqBody := models.PaymentRequest{
		Action:      "token",
		TerminalId:  terminalID,
		Amount:      tickets.TotalPrice,
		ResNum:      tickets.ID.String(),
		RedirectUrl: os.Getenv("HOST") + "/payment/success",
		CellNumber:  phoneNumber,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequest(http.MethodPost, bankUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := m.client.Do(req)
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (m *SamanGateway) VerifyTransaction(receivedPaymentRequest *models.PaymentReceipt) (models.Transaction, error) {
	var transaction models.Transaction
	reqBody := models.VerifyTransactionRequest{
		MID:    receivedPaymentRequest.MID,
		RefNum: receivedPaymentRequest.RefNum,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return transaction, err
	}

	req, err := http.NewRequest(http.MethodPost, bankVerifyTransactionUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return transaction, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := m.client.Do(req)
	if err != nil {
		return transaction, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&transaction)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (m *SamanGateway) ReverseTransaction(receivedPaymentRequest *models.PaymentReceipt) (models.Transaction, error) {
	var transaction models.Transaction
	reqBody := models.VerifyTransactionRequest{
		MID:    receivedPaymentRequest.MID,
		RefNum: receivedPaymentRequest.RefNum,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return transaction, err
	}

	req, err := http.NewRequest(http.MethodPost, bankReverseTransactionUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return transaction, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := m.client.Do(req)
	if err != nil {
		return transaction, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&transaction)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
