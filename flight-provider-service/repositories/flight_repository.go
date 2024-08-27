package repositories

import (
	"github.com/amirdashtii/Q/flight-provider-service/models"
)

func (p *Postgres) GetLastFlightDate(lastFlight *models.Flight) error {

	result := p.db.Order("departure_date DESC").First(&lastFlight)
	return result.Error
}

func (p *Postgres) CreateFlights(flights *[]models.Flight) error {

	result := p.db.Create(flights)
	return result.Error
}

func (p *Postgres) GetFlights(flightReq *models.FlightReq, flights *[]models.Flight) error {
	
	result := p.db.Where("source = ? AND destination = ? AND departure_date >= ? AND departure_date < ?", flightReq.Source, flightReq.Destination, flightReq.DepartureDateStart, flightReq.DepartureDateEnd).Find(flights)
	return result.Error
}
