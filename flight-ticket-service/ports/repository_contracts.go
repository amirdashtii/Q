package ports

import (
	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/google/uuid"
)

type FlightRepositoryContracts interface {
}

type TicketRepositoryContracts interface {
	Reserve(reservation *models.TicketReservation) error
	FindPassengersByIDs(userID *uuid.UUID, passengerIDs *[]uuid.UUID, passengers *[]models.Passenger) error
}

type PassengerRepositoryContracts interface {
	CreatePassenger(passenger *models.Passenger) error
	GetPassengers(userID *uuid.UUID, passengers *[]models.Passenger) error
	GetPassengerByID(passenger *models.Passenger) error
	UpdatePassenger(id, userID *uuid.UUID, updateItem map[string]interface{}) error
	DeletePassenger(passenger *models.Passenger) error
}
