package repositories

import (
	"github.com/amirdashtii/Q/flight-ticket-service/models"
)

func (p *Postgres) Reserve(reservation *models.Tickets) error {
	err := p.db.Create(reservation).Error
	if err != nil {
		return err
	}
	result := p.db.Preload("TicketItems.Passenger").Find(reservation)
	return result.Error
}

func (p *Postgres) GetReservationByID(reservation *models.Tickets) error {
	result := p.db.Preload("TicketItems.Passenger").Find(reservation)
	return result.Error
}

func (p *Postgres) GetTicketsByRefNum(resNum string) error {
	var tickets models.Tickets
	result := p.db.Where("reference_number = ?", resNum).First(tickets)
	return result.Error
}
