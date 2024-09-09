package validators

import (
	"errors"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
)

func ReservationValidation(reservationRequest *models.ReservationRequest, userIDStr string) error {

	if reservationRequest.FlightID == "" {
		return errors.New("flight ID is required")
	}

	if len(reservationRequest.PassengerIDs) == 0 {
		return errors.New("passenger IDs are required")
	}

	err := IDValidation(map[string]string{"user_id": userIDStr})
	if err != nil {
		return errors.New("user ID is required")
	}

	return nil
}
