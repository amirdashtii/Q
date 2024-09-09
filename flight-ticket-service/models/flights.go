package models

import (
	"time"

	"github.com/google/uuid"
)

type FlightSearchRequest struct {
	Source        string `json:"source"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departure_date"`
	SortBy        string `json:"sort_by,omitempty"`
	SortOrder     string `json:"order,omitempty"`
	Filter        string `json:"filter_by,omitempty"`
}

type Flight struct {
	FlightID       uuid.UUID   `json:"flight_id" gorm:"type:uuid"`
	FlightNumber   string      `json:"flight_number" gorm:"size:255"`
	Source         string      `json:"source" gorm:"size:255"`
	Destination    string      `json:"destination" gorm:"size:255"`
	DepartureDate  time.Time   `json:"departure_date"`
	FlightDuration int         `json:"flight_duration"`
	ArrivalDate    time.Time   `json:"arrival_date"`
	AirlineName    string      `json:"airline_name" gorm:"size:255"`
	AircraftName   string      `json:"aircraft_name" gorm:"size:255"`
	FareClass      string      `json:"fare_class" gorm:"size:255"`
	FlightClass    FlightClass `json:"flight_class" gorm:"size:255"`
}

type FlightProvider struct {
	DBModel
	FlightNumber   string      `json:"flight_number" gorm:"size:255"`
	Source         string      `json:"source" gorm:"size:255"`
	Destination    string      `json:"destination" gorm:"size:255"`
	DepartureDate  time.Time   `json:"departure_date"`
	FlightDuration int         `json:"flight_duration"`
	ArrivalDate    time.Time   `json:"arrival_date"`
	AirlineName    string      `json:"airline_name" gorm:"size:255"`
	AircraftName   string      `json:"aircraft_name" gorm:"size:255"`
	FareClass      FareClass   `gorm:"embedded"`
	Tax            int64       `json:"tax"`
	FlightClass    FlightClass `json:"flight_class" gorm:"size:255"`
	RemainingSeat  int         `json:"remaining_seat"`
}

type FareClass struct {
	AdultFare  int64 `json:"adult_fare"`
	ChildFare  int64 `json:"child_fare"`
	InfantFare int64 `json:"infant_fare"`
}

type FlightClass string

const (
	Economy    FlightClass = "Economy"
	Business   FlightClass = "Business"
	FirstClass FlightClass = "First Class"
)
