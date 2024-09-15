package ports

import "github.com/amirdashtii/Q/flight-ticket-service/models"

type FlightProviderContract interface {
	RequestFlights(flightReq *models.FlightSearchRequest, flights *[]models.ProviderFlight) error
	RequestFlight(id *string, flight *models.ProviderFlight) error
	ReserveTicketWithProvider(seats int, flightID string) error
	// CancelTicketWithProvider(reservation *models.Tickets) error
}
type PaymentGatewayProviderContract interface {
	CreatePayment(tickets *models.Tickets, phoneNumber string) (models.Response, error)
	VerifyTransaction(receivedPaymentRequest *models.ReceivedPaymentRequest) (models.Transaction, error)
}
