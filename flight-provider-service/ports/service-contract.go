package ports

import (
	"github.com/amirdashtii/Q/flight-provider-service/models"
	"github.com/google/uuid"
)

type FlightServiceContract interface {
	GenerateRandomFlightsForNext30Days() error
	GetFlights(flightReq *models.FlightReq, flights *[]models.Flight) error
	GetFlightByID(flight *models.Flight) error
	DecreaseFlightCapacity(id uuid.UUID, seats int) error
	IncreaseFlightCapacity(id uuid.UUID, seats int) error
}
