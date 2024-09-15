package ports

import (
	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/google/uuid"
)

type FlightServiceContract interface {
	GetFlights(flightReq *models.FlightSearchRequest, flights *[]models.ProviderFlight) error
	GetFlightByID(id *string, flight *models.ProviderFlight) error
}

type TicketServiceContract interface {
	CreateReservation(flightID string, passengerIDs []uuid.UUID, ticket *models.Tickets) error
	GetReservationByID(ticket *models.Tickets) error
	UpdateReservation(ticket *models.Tickets) error
	CancelTicket(ticket *models.Tickets) error
}
type PaymentServiceContract interface {
	PayTicket(ticket *models.Tickets, paymentGateway string) (string, error)
	VerifyTransaction(receivedPaymentRequest *models.ReceivedPaymentRequest) (models.Transaction, error)
}

type PassengerServiceContract interface {
	CreatePassenger(passengerReq *models.PassengerReq, passenger *models.Passenger) error
	GetPassengers(userID *uuid.UUID, passengers *[]models.Passenger) error
	GetPassengerByID(passenger *models.Passenger) error
	UpdatePassenger(passengerID, userID *uuid.UUID, passengerReq *models.PassengerReq) error
	DeletePassenger(passenger *models.Passenger) error
}
