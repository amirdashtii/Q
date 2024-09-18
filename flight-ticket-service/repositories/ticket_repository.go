package repositories

import (
	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/google/uuid"
)

func (p *Postgres) Reserve(tickets *models.Tickets) error {
	err := p.db.Create(tickets).Error
	if err != nil {
		return err
	}
	result := p.db.Preload("TicketItems.Passenger").Find(tickets)
	return result.Error
}

func (p *Postgres) GetTicketsByID(tickets *models.Tickets) error {
	result := p.db.Preload("TicketItems.Passenger").Where("id = ? AND user_id = ?", tickets.ID, tickets.UserID).First(tickets)
	return result.Error
}

func (p *Postgres) GetAllTickets(userID *uuid.UUID, tickets *[]models.Tickets) error {
	result := p.db.Preload("TicketItems.Passenger").Where("user_id = ?", userID).Find(tickets)
	return result.Error
}

func (p *Postgres) GetTicketsByRefNum(resNum string) error {
	var tickets models.Tickets
	result := p.db.Where("reference_number = ?", resNum).First(tickets)
	return result.Error
}

func (p *Postgres) UpdateReservation(id uuid.UUID, updateItem map[string]interface{}) error {
	var ticket models.Tickets
	ticket.ID = id
	result := p.db.Model(ticket).Updates(updateItem)
	return result.Error
}
