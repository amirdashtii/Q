package repositories

import (
	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/google/uuid"
)

func (p *Postgres) Reserve(reservation *models.TicketReservation) error {
	result := p.db.Create(reservation)
	return result.Error
}

func (p *Postgres) FindPassengersByIDs(userID *uuid.UUID, passengerIDs *[]uuid.UUID, passengers *[]models.Passenger) error {
	result := p.db.Where("id IN ? AND user_id = ?", *passengerIDs, userID).Find(passengers)
	return result.Error
}