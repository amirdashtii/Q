package service

import (
	"sort"
	"strconv"
	"strings"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/amirdashtii/Q/flight-ticket-service/ports"
	"github.com/amirdashtii/Q/flight-ticket-service/provider"
	"github.com/amirdashtii/Q/flight-ticket-service/repositories"
)

type FlightService struct {
	db ports.FlightRepositoryContracts
	pr ports.FlightProviderContract
}

func NewFlightService() *FlightService {
	db := repositories.NewPostgres()
	pr := provider.NewFlightProviderClient()

	return &FlightService{
		db: db,
		pr: pr,
	}
}

func (s *FlightService) GetFlights(flightReq *models.FlightSearchRequest, flights *[]models.Flight) error {

	err := s.pr.RequestFlights(flightReq, flights)
	if err != nil {
		return err
	}

	filteredFlights := s.applyFilters(*flights, flightReq.Filter)

	sortedFlights := s.applySorting(filteredFlights, flightReq.SortBy, flightReq.SortOrder)

	*flights = sortedFlights

	return nil
}

func (s *FlightService) applyFilters(flights []models.Flight, filter string) []models.Flight {
	var filteredFlights []models.Flight

	if filter == "" {
		return flights
	}

	for _, flight := range flights {
		if s.applyFilter(flight, filter) {
			filteredFlights = append(filteredFlights, flight)
		}
	}

	return filteredFlights
}

func (s *FlightService) applyFilter(flight models.Flight, filter string) bool {
	parts := strings.Split(filter, "=")
	if len(parts) != 2 {
		return false
	}

	fieldAndOperator := parts[0]
	valueStr := parts[1]

	fieldParts := strings.Split(fieldAndOperator, "_")
	if len(fieldParts) != 2 {
		return false
	}

	field := fieldParts[0]
	operator := fieldParts[1]

	switch field {
	case "price":
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return false
		}
		switch operator {
		case "gte":
			return flight.FareClass.AdultFare >= value
		case "lte":
			return flight.FareClass.AdultFare <= value
		case "gt":
			return flight.FareClass.AdultFare > value
		case "lt":
			return flight.FareClass.AdultFare < value
		case "eq":
			return flight.FareClass.AdultFare == value
		}
	case "duration":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return false
		}
		switch operator {
		case "gte":
			return flight.FlightDuration >= value
		case "lte":
			return flight.FlightDuration <= value
		case "gt":
			return flight.FlightDuration > value
		case "lt":
			return flight.FlightDuration < value
		case "eq":
			return flight.FlightDuration == value
		}
	}

	return false
}

func (s *FlightService) applySorting(flights []models.Flight, sortBy, sortOrder string) []models.Flight {
	if sortBy == "" {
		return flights
	}

	sort.SliceStable(flights, func(i, j int) bool {
		var less bool

		switch sortBy {
		case "price":
			less = flights[i].FareClass.AdultFare < flights[j].FareClass.AdultFare
		case "departure_date":
			less = flights[i].DepartureDate.Before(flights[j].DepartureDate)
		case "duration":
			less = flights[i].FlightDuration < flights[j].FlightDuration
		}

		if sortOrder == "desc" {
			return !less
		}
		return less
	})

	return flights
}

func (s *FlightService) GetFlightByID(id *string, flight *models.Flight) error {
	return s.pr.RequestFlight(id, flight)
}
