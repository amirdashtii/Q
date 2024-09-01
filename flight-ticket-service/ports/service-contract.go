package ports

import (
	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/google/uuid"
)

type FlightServiceContract interface {
	GetFlights(flightReq *models.FlightSearchRequest, flights *[]models.Flight) error
	GetFlightByID(id *string, flight *models.Flight) error
}

type TicketServiceContract interface {
	CreateReservation(passengerIDs *[]uuid.UUID, reservation *models.TicketReservation) error
}

type PassengerServiceContract interface {
	CreatePassenger(passengerReq *models.PassengerReq, passenger *models.Passenger) error
	GetPassengers(userID *uuid.UUID, passengers *[]models.Passenger) error
	GetPassengerByID(passenger *models.Passenger) error
	UpdatePassenger(passengerID, userID *uuid.UUID, passengerReq *models.PassengerReq) error
	DeletePassenger(passenger *models.Passenger) error
}
