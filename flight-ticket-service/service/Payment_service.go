package service

import (
	"errors"
	"os"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/amirdashtii/Q/flight-ticket-service/ports"
	"github.com/amirdashtii/Q/flight-ticket-service/provider"
	"github.com/amirdashtii/Q/flight-ticket-service/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentService struct {
	db             ports.RepositoryContracts
	ticketService  ports.TicketServiceContract
	flightprovider ports.FlightProviderContract
	paymentGetway  ports.PaymentGatewayProviderContract
}

func NewPaymentService() *PaymentService {
	db := repositories.NewPostgres()
	ticketService := NewTicketService()
	flightprovider := provider.NewProviderClient()
	paymentGetway := provider.NewSamanGateway()

	return &PaymentService{
		db:             db,
		ticketService:  ticketService,
		flightprovider: flightprovider,
		paymentGetway:  paymentGetway,
	}
}

func (s *PaymentService) PayTicketBySaman(tickets *models.Tickets) (string, error) {
	var pFlight models.ProviderFlight

	if err := s.db.GetTicketsByID(tickets); err != nil {
		return "", err
	}
	flightID := tickets.FlightID.String()
	seats := len(tickets.TicketItems)

	if err := s.flightprovider.RequestFlight(&flightID, &pFlight); err != nil {
		return "", err
	}

	if pFlight.RemainingSeat < seats {
		return "", errors.New("not enough seats")
	}

	if err := s.flightprovider.ReserveTicketWithProvider(seats, flightID); err != nil {
		return "", err
	}

	phoneNumber := ""

	response, err := s.paymentGetway.CreatePayment(tickets, phoneNumber)
	if err != nil {
		return "", err
	}

	if response.Status != 1 {
		return "", errors.New(models.BankErrorMessage[response.ErrorCode])
	}
	paymentLink := os.Getenv("BANK_SEND_TOKEN_URL") + "?token=" + response.Token

	return paymentLink, nil

}
func (s *PaymentService) PayTicket(tickets *models.Tickets, paymentGateway string) (string, error) {
	if paymentGateway == "Samsn Geteway" {
		return s.PayTicketBySaman(tickets)
	}
	return "", errors.New("payment gateway not found")
}

func (s *PaymentService) VerifyTransaction(receivedPaymentRequest *models.PaymentReceipt) (models.Transaction, error) {

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
	transaction, err := s.paymentGetway.VerifyTransaction(receivedPaymentRequest)
	if err != nil {
		return models.Transaction{}, err
	}

	// check if transaction is successful
	// if transaction is not successful return error
	if transaction.ResultCode != 0 {
		return models.Transaction{}, errors.New(transaction.ResultDescription)
	}

	if err := s.db.CreatePaymentReceipt(receivedPaymentRequest); err != nil {
		return models.Transaction{}, err
	}
	// if the amount of the transaction is equal to the amount of the order
	// then we can update the status of the order to payed
	// and save the transaction to the database
	if transaction.TransactionDetail.AffectiveAmount == transaction.TransactionDetail.OrginalAmount {
		var tickets models.Tickets
		reservationID, err := uuid.Parse(receivedPaymentRequest.ResNum)

		if err != nil {
			return models.Transaction{}, err
		}

		tickets.ID = reservationID

		tickets.Status = "payed"
		tickets.ReferenceNumber = transaction.TransactionDetail.RefNum
		if err := s.ticketService.UpdateTickets(&tickets); err != nil {
			return models.Transaction{}, err
		}

	} else {
		// otherwise if the amount of the transaction is not equal to the amount of the order
		// then we should reverse the transaction
		transaction, err := s.paymentGetway.ReverseTransaction(receivedPaymentRequest)
		if err != nil {
			return models.Transaction{}, err
		}

		if transaction.ResultCode != 0 {
			return models.Transaction{}, errors.New(transaction.ResultDescription)
		}

	}

	return transaction, nil
}
