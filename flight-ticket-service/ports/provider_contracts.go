package ports

import "github.com/amirdashtii/Q/flight-ticket-service/models"

type FlightProviderContract interface {
	RequestFlights(flightReq *models.FlightSearchRequest, flights *[]models.Flight) error
	RequestFlight(id *string, flight *models.Flight) error
}
