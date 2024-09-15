package service

import (
	"errors"
	"os"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/amirdashtii/Q/flight-ticket-service/ports"
	"github.com/amirdashtii/Q/flight-ticket-service/provider"
	"github.com/amirdashtii/Q/flight-ticket-service/repositories"
	"gorm.io/gorm"
)

type PaymentService struct {
	db   ports.RepositoryContracts
	pr   ports.FlightProviderContract
	pgpr ports.PaymentGatewayProviderContract
}

func NewPaymentService() *PaymentService {
	db := repositories.NewPostgres()
	pr := provider.NewProviderClient()
	pgpr := provider.NewSamanGateway()

	return &PaymentService{
		db:   db,
		pr:   pr,
		pgpr: pgpr,
	}
}

var (
// terminalID       string
// bankUrl          string
// banksendTokenUrl string
// bankVerifyTransactionUrl  string
// bankReverseTransactionUrl string
)

func (s *PaymentService) PayTicketBySaman(tickets *models.Tickets) (string, error) {
	var flight models.ProviderFlight
	// var user models.User

	seats := len(tickets.TicketItems)

	if err := s.db.GetReservationByID(tickets); err != nil {
		return "", err
	}
	flightID := tickets.FlightID.String()

	if err := s.pr.RequestFlight(&flightID, &flight); err != nil {
		return "", err
	}

	if flight.RemainingSeat < seats {
		return "", errors.New("not enough seats")
	}

	if err := s.pr.ReserveTicketWithProvider(seats, flightID); err != nil {
		return "", err
	}
	// TODO: change status
	// tickets.Status = "reserved"

	// TODO: get user
	phoneNumber := "09123456789"

	response, err := s.pgpr.CreatePayment(tickets, phoneNumber)

	if err != nil {
		return "", err
	}

	if response.Status != 1 {
		return "", errors.New(models.BankErrorMessage[response.ErrorCode])
	}
	paymentLink := os.Getenv("BANK_SEND_TOKEN_URL") + "?token=" + response.Token
	// pay
	// PayTicket(tickets)
	// cancel if not pay

	// save to database
	return paymentLink, nil

}
func (s *PaymentService) PayTicket(tickets *models.Tickets, paymentGateway string) (string, error) {
	if paymentGateway == "Samsn Geteway" {
		return s.PayTicketBySaman(tickets)
	}
	return "", errors.New("payment gateway not found")
}

func (s *PaymentService) VerifyTransaction(receivedPaymentRequest *models.ReceivedPaymentRequest) (models.Transaction, error) {

	if receivedPaymentRequest.State != "OK" {
		return models.Transaction{}, errors.New(models.BankErrorMessage[receivedPaymentRequest.Status])
	}

	// search RefNum in database for Double Spending
	err := s.db.GetTicketsByRefNum(receivedPaymentRequest.ResNum)
	if err != nil && err != gorm.ErrRecordNotFound {
		return models.Transaction{}, err
	}
	if err == nil {
		return models.Transaction{}, errors.New("double spending")
	}

	// call VerifyTransaction in bank whit RefNum and MID
	transaction, err := s.pgpr.VerifyTransaction(receivedPaymentRequest)
	if err != nil {
		return models.Transaction{}, err
	}

	// check if transaction is successful
	// if transaction is not successful return error
	if transaction.ResultCode != 0 {
		return models.Transaction{}, errors.New(transaction.ResultDescription)
	}

	// if the amount of the transaction is equal to the amount of the order
	// then we can update the status of the order to payed
	// and save the transaction to the database
	if transaction.PurchaseInfo.AffectiveAmount == transaction.PurchaseInfo.OrginalAmount {
		// TODO: update status to payed
		// TODO: save transaction to database
	} else {
		// otherwise if the amount of the transaction is not equal to the amount of the order
		// then we should reverse the transaction
		// TODO: reverse transaction
	}

	return transaction, nil
}
