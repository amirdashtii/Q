package ports

import "github.com/amirdashtii/Q/flight-provider-service/models"

type FlightServiceContract interface {
	GenerateRandomFlightsForNext30Days() error
	GetFlights(flightReq *models.FlightReq, flights *[]models.Flight) error
}
