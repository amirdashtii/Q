package service

import (
	"errors"
	"time"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/amirdashtii/Q/flight-ticket-service/ports"
	"github.com/amirdashtii/Q/flight-ticket-service/provider"
	"github.com/amirdashtii/Q/flight-ticket-service/repositories"
	"github.com/google/uuid"
)

type TicketService struct {
	db   ports.RepositoryContracts
	pr   ports.FlightProviderContract
	pgpr ports.PaymentGatewayProviderContract
}

func NewTicketService() *TicketService {
	db := repositories.NewPostgres()
	pr := provider.NewProviderClient()
	pgpr := provider.NewSamanGateway()

	return &TicketService{
		db:   db,
		pr:   pr,
		pgpr: pgpr,
	}
}

func (s *TicketService) CreateReservation(flightID string, passengerIDs []uuid.UUID, tickets *models.Tickets) error {

	var passengers []models.Passenger
	var flightProvider models.ProviderFlight

	if err := s.db.FindPassengersByIDs(&tickets.UserID, &passengerIDs, &passengers); err != nil {
		return err
	}

	if err := s.pr.RequestFlight(&flightID, &flightProvider); err != nil {
		return err
	}

	err := addFlightAndPassengersToTicket(tickets, flightProvider, passengers)
	if err != nil {
		return nil
	}
	tickets.FlightID = flightProvider.ID

	tickets.Status = "reserved"

	return s.db.Reserve(tickets)
}
func (s *TicketService) GetTicketsByID(tickets *models.Tickets) error {
	return s.db.GetReservationByID(tickets)
}

func (s *TicketService) CancelTicket(tickets *models.Tickets) error {
	if err := s.db.GetReservationByID(tickets); err != nil {
		return errors.New("ticket not found")
	}

	if tickets.Status == "cancelled" {
		return errors.New("ticket already cancelled")
	}

	if err := s.pr.CancelTicketWithProvider(len(tickets.TicketItems), tickets.FlightID.String()); err != nil {
		return err
	}

	updateItem := make(map[string]interface{})
	updateItem["status"] = "cancelled"
	return s.db.UpdateReservation(tickets.ID, updateItem)
}

func (s *TicketService) UpdateTickets(tickets *models.Tickets) error {
	updateItem := make(map[string]interface{})

	if tickets.Status != "" {
		updateItem["status"] = tickets.Status
	}
	if tickets.ReferenceNumber != "" {
		updateItem["reference_number"] = tickets.ReferenceNumber
	}

	return s.db.UpdateReservation(tickets.ID, updateItem)
}

func addFlightAndPassengersToTicket(tickets *models.Tickets, fp models.ProviderFlight, passengers []models.Passenger) error {
	for _, passenger := range passengers {
		var t models.TicketItem
		t.PassengerID = passenger.ID

		t.Flight.FlightNumber = fp.FlightNumber
		t.Flight.Source = fp.Source
		t.Flight.Destination = fp.Destination
		t.Flight.DepartureDate = fp.DepartureDate
		t.Flight.FlightDuration = fp.FlightDuration
		t.Flight.ArrivalDate = fp.ArrivalDate
		t.Flight.AirlineName = fp.AirlineName
		t.Flight.AircraftName = fp.AircraftName
		t.Flight.FlightClass = fp.FlightClass

		fareClass, price, err := addFareClass(*passenger.DateOfBirth, fp.FareClass)
		if err != nil {
			return err
		}

		t.Flight.FareClass = fareClass
		t.Price = price

		tickets.TicketItems = append(tickets.TicketItems, t)
		tickets.TotalPrice += price
	}

	return nil
}

func addFareClass(dateOfBirth time.Time, fareClass models.FareClass) (string, int64, error) {
	now := time.Now()

	if dateOfBirth.After(now.AddDate(0, -2, 0)) {
		return "", 0, errors.New("infant must be at least 2 months old")
	}
	if dateOfBirth.After(now.AddDate(-2, 0, 0)) {
		return "Infant", fareClass.InfantFare, nil
	}
	if dateOfBirth.After(now.AddDate(-12, 0, 0)) {
		return "Child", fareClass.ChildFare, nil
	}
	return "Adult", fareClass.AdultFare, nil

}
