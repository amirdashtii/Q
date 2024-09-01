package service

import (
	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/amirdashtii/Q/flight-ticket-service/ports"
	"github.com/amirdashtii/Q/flight-ticket-service/provider"
	"github.com/amirdashtii/Q/flight-ticket-service/repositories"
	"github.com/google/uuid"
)

type TicketService struct {
	db ports.TicketRepositoryContracts
	pr ports.TicketProviderContract
}

func NewTicketService() *TicketService {
	db := repositories.NewPostgres()
	pr := provider.NewTicketProviderClient()

	return &TicketService{
		db: db,
		pr: pr,
	}
}

func (s *TicketService) CreateReservation(passengerIDs *[]uuid.UUID, reservation *models.TicketReservation) error {

	if err := s.db.FindPassengersByIDs(&reservation.UserID, passengerIDs, &reservation.Passengers); err != nil {
		return err
	}

	if err := s.pr.ReserveTicketWithProvider(reservation); err != nil {
		return err
	}

	if err := s.db.Reserve(reservation); err != nil {
		if err := s.pr.CancelTicketWithProvider(reservation); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (s *TicketService) FindPassengersByIDs(userID *uuid.UUID, passengerIDs *[]uuid.UUID, passengers *[]models.Passenger) error {
	return s.db.FindPassengersByIDs(userID, passengerIDs, passengers)
}
