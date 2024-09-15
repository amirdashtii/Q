package models

import (
	"github.com/google/uuid"
)

type Tickets struct {
	DBModel
	UserID          uuid.UUID    `json:"user_id" gorm:"not null"`
	FlightID        uuid.UUID    `json:"flight_id"`
	TicketItems     []TicketItem `json:"ticket_items" gorm:"foreignKey:ReservationID;constraint:OnDelete:CASCADE;"`
	TotalPrice      int64        `json:"total_price"`
	Status          string       `json:"status" gorm:"not null"`
	ReferenceNumber string       `json:"reference_number"`
}

type TicketItem struct {
	DBModel
	ReservationID uuid.UUID `json:"reservation_id" gorm:"not null"`
	PassengerID   uuid.UUID `json:"passenger_id" gorm:"not null"`
	Passenger     Passenger `json:"passenger" gorm:"foreignKey:PassengerID;constraint:OnDelete:CASCADE;"`
	Flight        Flight    `json:"flight" gorm:"embedded"`
	Price         int64     `json:"price"`
}

type ReservationRequest struct {
	PassengerIDs []string `json:"passenger_ids"`
	FlightID     string   `json:"flight_id"`
}

