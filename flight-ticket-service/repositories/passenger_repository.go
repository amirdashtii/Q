package repositories

import (
	"errors"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/google/uuid"
)

func (p *Postgres) CreatePassenger(passenger *models.Passenger) error {
	result := p.db.Create(passenger)
	return result.Error
}

func (p *Postgres) GetPassengers(userID *uuid.UUID, passengers *[]models.Passenger) error {

	result := p.db.Where("user_id = ?", userID).Find(passengers)
	return result.Error
}

func (p *Postgres) GetPassengerByID(passenger *models.Passenger) error {
	result := p.db.Where("user_id = ? AND id = ?", passenger.UserID, passenger.ID).First(passenger)
	return result.Error
}

func (p *Postgres) UpdatePassenger(id, userID *uuid.UUID, updateItem map[string]interface{}) error {
	var passenger models.Passenger
	result := p.db.Model(&passenger).Where("user_id = ? AND id = ?", userID, id).Updates(updateItem)
	if result.RowsAffected == 0 {
		return errors.New("passenger not found")
	}
	return result.Error
}

func (p *Postgres) DeletePassenger(passenger *models.Passenger) error {
	result := p.db.Where("user_id = ?", passenger.UserID).Delete(passenger)
	if result.RowsAffected == 0 {
		return errors.New("passenger not found")
	}
	return result.Error
}

