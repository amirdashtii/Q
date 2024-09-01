package models

import (
	"github.com/google/uuid"
)

type TicketReservation struct {
	DBModel
	UserID     uuid.UUID   `json:"user_id" gorm:"not null"`
	FlightID   uuid.UUID   `json:"flight_id" gorm:"not null"`
	Status     string      `json:"status" gorm:"not null"`
	Passengers []Passenger `json:"passengers" gorm:"many2many:TicketReservation_passengers;ForeignKey:id"`
}

type TicketReservationRequest struct {
	PassengerIDs []string `json:"passenger_ids"`
	FlightID     string   `json:"flight_id"`
}
