package ports

import (
	"github.com/amirdashtii/Q/flight-provider-service/models"
)

type FlightRepositoryContracts interface {
	GetLastFlightDate(lastFlight *models.Flight) error
	CreateFlights(flights *[]models.Flight) error
	GetFlights(flightReq *models.FlightSearchRequest, flights *[]models.Flight) error
	GetFlightByID(flight *models.Flight) error
	UpdateFlight(flight *models.Flight) error
}
