package service

import (
	"errors"
	"time"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/amirdashtii/Q/flight-ticket-service/ports"
	"github.com/amirdashtii/Q/flight-ticket-service/repositories"
	"github.com/google/uuid"
)

type PassengerService struct {
	db ports.PassengerRepositoryContracts
}

func NewPassengerService() *PassengerService {
	db := repositories.NewPostgres()

	return &PassengerService{
		db: db,
	}
}

func (u *PassengerService) CreatePassenger(passengerReq *models.PassengerReq, passenger *models.Passenger) error {

	passenger.FirstName = passengerReq.FirstName
	passenger.LastName = passengerReq.LastName
	dateOfBirth, err := time.Parse("2006-01-02", passengerReq.DateOfBirth)
	if err != nil {
		return errors.New("invalid input date format")
	}
	passenger.DateOfBirth = &dateOfBirth
	passenger.Nationality = passengerReq.Nationality
	passenger.NationalCode = passengerReq.NationalCode
	passenger.PassportNumber = passengerReq.PassportNumber
	passenger.Gender = passengerReq.Gender
	return u.db.CreatePassenger(passenger)
}

func (u *PassengerService) GetPassengers(userID *uuid.UUID, passenger *[]models.Passenger) error {
	return u.db.GetPassengers(userID, passenger)
}

func (u *PassengerService) GetPassengerByID(passenger *models.Passenger) error {
	return u.db.GetPassengerByID(passenger)
}

func (u *PassengerService) UpdatePassenger(passengerID, userID *uuid.UUID, passengerReq *models.PassengerReq) error {
	updateItem := make(map[string]interface{})

	if passengerReq.FirstName != "" {
		updateItem["first_name"] = passengerReq.FirstName
	}

	if passengerReq.LastName != "" {
		updateItem["last_name"] = passengerReq.LastName
	}

	if passengerReq.DateOfBirth != "" {
		dateOfBirth, err := time.Parse("2006-01-02", passengerReq.DateOfBirth)
		if err != nil {
			return errors.New("invalid input date format")
		}
		updateItem["date_of_birth"] = dateOfBirth
	}

	if passengerReq.Nationality != "" {
		updateItem["nationality"] = passengerReq.Nationality
	}

	if passengerReq.NationalCode != "" {
		updateItem["national_code"] = passengerReq.NationalCode
	}

	if passengerReq.PassportNumber != "" {
		updateItem["passport_number"] = passengerReq.PassportNumber
	}

	if passengerReq.Gender != "" {
		updateItem["gender"] = passengerReq.Gender
	}

	return u.db.UpdatePassenger(passengerID, userID, updateItem)
}

func (u *PassengerService) DeletePassenger(passenger *models.Passenger) error {
	return u.db.DeletePassenger(passenger)
}
