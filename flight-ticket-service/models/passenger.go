package models

import (
	"time"

	"github.com/google/uuid"
)

type Passenger struct {
	DBModel
	UserID         uuid.UUID  `json:"user_id" gorm:"not null"`
	FirstName      string     `json:"first_name" gorm:"size:255 not null"`
	LastName       string     `json:"last_name" gorm:"size:255 not null"`
	DateOfBirth    *time.Time `json:"date_of_birth" gorm:"not null"`
	Nationality    string     `json:"nationality"`
	NationalCode   string     `json:"national_code"`
	PassportNumber string     `json:"passport_number"`
	Gender         string     `json:"gender"`
}

type PassengerReq struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DateOfBirth    string `json:"date_of_birth"`
	Nationality    string `json:"nationality"`
	NationalCode   string `json:"national_code"`
	PassportNumber string `json:"passport_number"`
	Gender         string `json:"gender"`
}
