package validators

import (
	"errors"
	"time"

	"github.com/amirdashtii/Q/flight-provider-service/models"
)

func validateDepartingDate(departingStr string) error {
	if departingStr == "" {
		return errors.New("departing is required")
	}

	departing, err := time.Parse("2006-01-02", departingStr)
	if err != nil {
		return errors.New("invalid input date format")
	}

	now := time.Now().UTC()
	if departing.Before(now) {
		return errors.New("past date is not allowed")
	}

	return nil
}

func ValidateFlightParam(fReq *models.FlightSearchRequest) error {
	if fReq.Source == "" {
		return errors.New("source is required")
	}

	if fReq.Destination == "" {
		return errors.New("destination is required")
	}

	if fReq.Source == fReq.Destination {
		return errors.New("source and destination cannot be the same")
	}

	err := validateDepartingDate(fReq.DepartureDate)
	if err != nil {
		return err
	}

	return nil
}
