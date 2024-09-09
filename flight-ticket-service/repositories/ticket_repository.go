package repositories

import (
	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/google/uuid"
)

func (p *Postgres) Reserve(reservation *models.Reservation) error {
	err := p.db.Create(reservation).Error
	if err != nil {
		return err
	}
	result := p.db.Preload("TicketItems.Passenger").Find(reservation)
	return result.Error
}

func (p *Postgres) GetReservationByID(reservation *models.Reservation) error {
	result := p.db.Preload("TicketItems.Passenger").Find(reservation)
	return result.Error
}

func (p *Postgres) FindPassengersByIDs(userID *uuid.UUID, passengerIDs *[]uuid.UUID, passengers *[]models.Passenger) error {
	result := p.db.Where("user_id = ?", userID).Find(passengers, passengerIDs)
	return result.Error
}
