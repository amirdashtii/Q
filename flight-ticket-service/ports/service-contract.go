package ports

import (
	"github.com/amirdashtii/Q/flight-ticket-service/models"
)

type FlightServiceContract interface {
	GetFlights(flightReq *models.FlightSearchRequest, flights *[]models.Flight) error
	GetFlightByID(id *string, flight *models.Flight) error
}
