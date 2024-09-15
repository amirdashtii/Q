package ports

import (
	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/google/uuid"
)

type RepositoryContracts interface {
	// passenger
	CreatePassenger(passenger *models.Passenger) error
	GetPassengers(userID *uuid.UUID, passengers *[]models.Passenger) error
	GetPassengerByID(passenger *models.Passenger) error
	UpdatePassenger(id, userID *uuid.UUID, updateItem map[string]interface{}) error
	DeletePassenger(passenger *models.Passenger) error
	FindPassengersByIDs(userID *uuid.UUID, passengerIDs *[]uuid.UUID, passenger *[]models.Passenger) error

	// ticket
	Reserve(ticket *models.Tickets) error
	GetReservationByID(ticket *models.Tickets) error
	GetTicketsByRefNum(resNum string) error
	UpdateReservation(ticket *models.Tickets) error
}
