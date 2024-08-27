package service

import (
	"fmt"
	"time"

	"math/rand/v2"

	"github.com/amirdashtii/Q/flight-provider-service/models"
	"github.com/amirdashtii/Q/flight-provider-service/ports"
	"github.com/amirdashtii/Q/flight-provider-service/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlightService struct {
	db ports.FlightRepositoryContracts
}

func NewFlightService() *FlightService {
	db := repositories.NewPostgres()

	return &FlightService{db: db}
}

func (s *FlightService) GetLastFlightDate() (time.Time, error) {
	var lastFlight models.Flight
	err := s.db.GetLastFlightDate(&lastFlight)

	if err != nil && err != gorm.ErrRecordNotFound {
		return time.Time{}, err
	}

	if lastFlight.ID == uuid.Nil {
		return time.Now(), nil
	}

	return lastFlight.DepartureDate, nil
}

func (s *FlightService) GenerateRandomFlightsForNext30Days() error {
	lastDate, err := s.GetLastFlightDate()
	if err != nil {
		return err
	}

	today := time.Now()
	daysToGenerate := 30 - int(lastDate.Sub(today).Hours()/24)

	if daysToGenerate <= 0 {
		return nil
	}
	if daysToGenerate > 30 {
		daysToGenerate = 30
	}

	for i := 1; i <= daysToGenerate; i++ {
		date := lastDate.AddDate(0, 0, i)
		flights := s.generateRandomFlightsForDay(date)
		err := s.db.CreateFlights(&flights)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *FlightService) generateRandomFlightsForDay(day time.Time) []models.Flight {
	numFlights := 250
	flights := make([]models.Flight, numFlights)

	for i := 0; i < numFlights; i++ {
		source, destination := generateRandomCityPair()

		flightDuration := rand.IntN(90) + 30

		departureTime := day.Add(time.Duration(rand.IntN(24)) * time.Hour)

		arrivalTime := departureTime.Add(time.Duration(flightDuration) * time.Minute)

		flight := models.Flight{
			DBModel:        models.DBModel{ID: uuid.New()},
			FlightNumber:   generateRandomFlightNumber(),
			Source:         source,
			Destination:    destination,
			DepartureDate:  departureTime,
			FlightDuration: flightDuration,
			ArrivalDate:    arrivalTime,
			AirlineName:    generateRandomAirline(),
			AircraftName:   generateRandomAircraft(),
			FareClass: models.FareClass{
				AdultFare:  (rand.Int64N(100) + 100) * 10000,
				ChildFare:  (rand.Int64N(100) + 70) * 10000,
				InfantFare: (rand.Int64N(100) + 40) * 10000,
			},
			Tax:           (rand.Int64N(20) + 10) * 10000,
			FlightClass:   generateRandomFlightClass(),
			RemainingSeat: rand.IntN(100),
		}
		flights[i] = flight
	}

	return flights
}

func generateRandomCityPair() (string, string) {
	cities := []string{"Tehran", "Mashhad", "Isfahan", "Shiraz", "Tabriz", "Ahvaz", "Kerman", "Rasht", "Yazd", "Zahedan", "Kermanshah", "Urmia", "Qazvin", "Hamadan", "Bandar Abbas"}
	source := cities[rand.IntN(len(cities))]
	destination := cities[rand.IntN(len(cities))]

	for source == destination {
		destination = cities[rand.IntN(len(cities))]
	}

	return source, destination
}

func generateRandomFlightNumber() string {
	return fmt.Sprintf("FL%d", rand.IntN(9000)+1000)
}

func generateRandomAirline() string {
	airlines := []string{"Iran Air", "Mahan Air", "Aseman Airlines", "Taban Air", "Kish Air", "Qeshm Air", "Caspian Airlines", "Zagros Airlines", "ATA Airlines", "Saha Airlines"}
	return airlines[rand.IntN(len(airlines))]
}

func generateRandomAircraft() string {
	aircrafts := []string{"Airbus A320", "Airbus A310", "Airbus A300", "Airbus A321", "Airbus A340", "Boeing 737", "Boeing 747", "Fokker 100", "Fokker 50", "Fokker 70"}
	return aircrafts[rand.IntN(len(aircrafts))]
}

func generateRandomFlightClass() models.FlightClass {
	classes := []models.FlightClass{"Economy", "Business", "First Class"}
	return classes[rand.IntN(len(classes))]
}

func (s *FlightService) GetFlights(flightReq *models.FlightReq, flights *[]models.Flight) error {

	dayStart, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", flightReq.DepartureDate))
	if err != nil {
		return err
	}

	flightReq.DepartureDateStart = dayStart
	flightReq.DepartureDateEnd = dayStart.Add(24 * time.Hour)

	return s.db.GetFlights(flightReq, flights)
}
