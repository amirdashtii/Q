package ports

import "github.com/amirdashtii/Q/flight-ticket-service/models"

type ProviderContract interface {
	RequestFlights(flightReq *models.FlightSearchRequest, flights *[]models.FlightProvider) error
	RequestFlight(id *string, flight *models.FlightProvider) error
	// ReserveTicketWithProvider(reservation *models.Tickets) error
	// CancelTicketWithProvider(reservation *models.Tickets) error
}
