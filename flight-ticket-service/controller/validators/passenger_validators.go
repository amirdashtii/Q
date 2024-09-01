package validators

import (
	"errors"
	"regexp"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
)

func ValidateNationalID(nationalID string) error {
	regex := regexp.MustCompile(`^[0-9]{10}$`)

	if !regex.MatchString(nationalID) {
		return errors.New("invalid national ID format")
	}

	return nil
}

func PassengerValidation(passengerReq models.PassengerReq, userIDStr string) error {
	if passengerReq.FirstName == "" {
		return errors.New("first name is required")
	}

	if passengerReq.LastName == "" {
		return errors.New("last name is required")
	}

	err := ValidateDate(passengerReq.DateOfBirth)
	if err != nil {
		return err
	}

	if passengerReq.NationalCode != "" {
		err := ValidateNationalID(passengerReq.NationalCode)
		if err != nil {
			return err
		}
	}

	err = IDValidation(map[string]string{"user_id": userIDStr})
	if err != nil {
		return errors.New("user ID is required")
	}

	return nil
}

func PassengerUpdateValidation(passengerIDStr, userIDStr, dateOfBirth string) error {
	err := IDValidation(map[string]string{"user id": userIDStr, "passenger id": passengerIDStr})
	if err != nil {
		return err
	}
	
	if dateOfBirth != "" {
		err := ValidateDate(dateOfBirth)
		if err != nil {
			return err
		}
	}

	return nil
}