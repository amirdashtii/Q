package ports

import (
	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/google/uuid"
)

type FlightServiceContract interface {
	GetFlights(flightReq *models.FlightSearchRequest, flights *[]models.FlightProvider) error
	GetFlightByID(id *string, flight *models.FlightProvider) error
}

type TicketServiceContract interface {
	CreateReservation(flightID string, passengerIDs []uuid.UUID, Reservation *models.Reservation) error
	GetReservationByID(Reservation *models.Reservation) error
}

type PassengerServiceContract interface {
	CreatePassenger(passengerReq *models.PassengerReq, passenger *models.Passenger) error
	GetPassengers(userID *uuid.UUID, passengers *[]models.Passenger) error
	GetPassengerByID(passenger *models.Passenger) error
	UpdatePassenger(passengerID, userID *uuid.UUID, passengerReq *models.PassengerReq) error
	DeletePassenger(passenger *models.Passenger) error
}
