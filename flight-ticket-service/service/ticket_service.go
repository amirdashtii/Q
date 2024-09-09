package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/amirdashtii/Q/flight-ticket-service/ports"
	"github.com/amirdashtii/Q/flight-ticket-service/provider"
	"github.com/amirdashtii/Q/flight-ticket-service/repositories"
	"github.com/google/uuid"
)

type TicketService struct {
	db ports.TicketRepositoryContracts
	pr ports.ProviderContract
}

func NewTicketService() *TicketService {
	db := repositories.NewPostgres()
	pr := provider.NewProviderClient()

	return &TicketService{
		db: db,
		pr: pr,
	}
}

func (s *TicketService) CreateReservation(flightID string, passengerIDs []uuid.UUID, reservation *models.Reservation) error {

	var passengers []models.Passenger
	var flightProvider models.FlightProvider

	if err := s.db.FindPassengersByIDs(&reservation.UserID, &passengerIDs, &passengers); err != nil {
		return err
	}

	if err := s.pr.RequestFlight(&flightID, &flightProvider); err != nil {
		return err
	}

	err := addFlightAndPassengersToTicket(reservation, flightProvider, passengers)
	if err != nil {
		return nil
	}

	reservation.Status = "reserved"

	return s.db.Reserve(reservation)
}
func (s *TicketService) GetReservationByID(reservation *models.Reservation) error {
	return s.db.GetReservationByID(reservation)
}

func (s *TicketService) FindPassengersByIDs(userID *uuid.UUID, passengerIDs *[]uuid.UUID, passengers *[]models.Passenger) error {
	return s.db.FindPassengersByIDs(userID, passengerIDs, passengers)
}

func addFlightAndPassengersToTicket(reservation *models.Reservation, fp models.FlightProvider, passengers []models.Passenger) error {
	for _, passenger := range passengers {
		var t models.TicketItem
		t.PassengerID = passenger.ID

		t.Flight.FlightID = fp.ID
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

		fmt.Println(t)

		reservation.TicketItems = append(reservation.TicketItems, t)

		reservation.TotalPrice += price
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
