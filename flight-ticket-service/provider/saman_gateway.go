package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
)

var (
	terminalID               string
	bankUrl                  string
	banksendTokenUrl         string
	bankVerifyTransactionUrl string
	// bankReverseTransactionUrl string
)

type SamanGateway struct {
	client *http.Client
}

func NewSamanGateway() *SamanGateway {
	terminalID = os.Getenv("BANK_TERMINAL_ID")
	bankUrl = os.Getenv("BANK_URL")
	banksendTokenUrl = os.Getenv("BANK_SEND_TOKEN_URL")
	bankVerifyTransactionUrl = os.Getenv("BANK_VERIFY_TRANSACTION_URL")
	// bankReverseTransactionUrl = os.Getenv("BANK_REVERSE_TRANSACTION_URL")

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

func (m *SamanGateway) VerifyTransaction(receivedPaymentRequest *models.ReceivedPaymentRequest) (models.Transaction, error) {
	var transaction models.Transaction
	reqBody := models.VerifyTransactionRequest{
		MID: "134756366",
		// RefNum: "f0e05013-1ebc-4d8c-8356-922fa4f954e5",
		// MID:    receivedPaymentRequest.MID,
		RefNum: receivedPaymentRequest.RefNum,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("1")
		return transaction, err
	}

	req, err := http.NewRequest(http.MethodPost, bankVerifyTransactionUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("2")
		return transaction, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := m.client.Do(req)
	if err != nil {
		fmt.Println("3")
		return transaction, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&transaction)
	if err != nil {
		fmt.Println("4")
		return transaction, err
	}
	fmt.Println(transaction)

	return transaction, nil
}

func (m *SamanGateway) ReverseTransaction(receivedPaymentRequest *models.ReceivedPaymentRequest) error {

	return nil
}
