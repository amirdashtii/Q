package models

import "time"

type FlightReq struct {
	Source             string `gorm:"size:255" json:"source"`
	Destination        string `gorm:"size:255" json:"destination"`
	DepartureDate      string `json:"departure_date"`
	DepartureDateStart time.Time
	DepartureDateEnd   time.Time
}

type Flight struct {
	DBModel
	FlightNumber   string      `gorm:"size:255" json:"flight_number"`
	Source         string      `gorm:"size:255" json:"source"`
	Destination    string      `gorm:"size:255" json:"destination"`
	DepartureDate  time.Time   `json:"departure_date"`
	FlightDuration int         `json:"flight_duration"`
	ArrivalDate    time.Time   `json:"arrival_date"`
	AirlineName    string      `gorm:"size:255" json:"airline_name"`
	AircraftName   string      `gorm:"size:255" json:"aircraft_name"`
	FareClass      FareClass   `gorm:"embedded"`
	Tax            int64       `json:"tax"`
	FlightClass    FlightClass `gorm:"size:255" json:"flight_class"`
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
