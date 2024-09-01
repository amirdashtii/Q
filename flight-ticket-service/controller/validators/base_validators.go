package validators

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

func IDValidation(idsStr map[string]string) error {

	for key, id := range idsStr {
		if id == "" {
			return errors.New(key + " is required")
		}

		if _, err := uuid.Parse(id); err != nil {
			return errors.New(key + " is invalid")
		}
	}

	return nil
}

func ValidateDate(date string) error {
	if date == "" {
		return errors.New("departing date is required")
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		return errors.New("invalid input date format")
	}
	
	return nil
}
