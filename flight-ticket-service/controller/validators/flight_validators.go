package validators

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
)

// Validate the departure date
func validateDepartingDate(departingStr string) error {

	if err := ValidateDate(departingStr); err != nil {
		return err
	}

	now := time.Now().UTC()
	departing, _ := time.Parse("2006-01-02", departingStr)
	if departing.Before(now) {
		return errors.New("past date is not allowed")
	}

	return nil
}

func validateSortOrder(sortOrder string) error {
	if sortOrder == "" {
		return nil // Sort order is optional
	}

	sortOrder = strings.ToLower(sortOrder)
	if sortOrder != "asc" && sortOrder != "desc" {
		return errors.New("invalid sort order, must be 'asc' or 'desc'")
	}

	return nil
}

func validateSortBy(sortBy string) error {
	if sortBy == "" {
		return nil
	}

	validSortFields := []string{"price", "departure_date", "duration"}
	for _, field := range validSortFields {
		if sortBy == field {
			return nil
		}
	}

	return errors.New("invalid sort field, must be one of 'price', 'departure_date', or 'duration'")
}

func validateFilter(filter string) error {
	if filter == "" {
		return nil
	}

	validFields := []string{"price", "duration"}
	validOperators := []string{"lt", "gt", "lte", "gte", "eq"}

	parts := strings.Split(filter, "=")
	if len(parts) != 2 {
		return errors.New("invalid filter format")
	}

	fieldAndOperator := parts[0]
	value := parts[1]

	parts = strings.Split(fieldAndOperator, "_")
	if len(parts) != 2 {
		return errors.New("invalid filter format")
	}

	field := parts[0]
	operator := parts[1]

	isValidField := false
	for _, validField := range validFields {
		if strings.TrimSpace(field) == validField {
			isValidField = true
			break
		}
	}

	if !isValidField {
		return errors.New("invalid filter field, must be one of 'price', 'duration'")
	}

	isValidOperator := false
	for _, validOperator := range validOperators {
		if strings.TrimSpace(operator) == validOperator {
			isValidOperator = true
			break
		}
	}

	if !isValidOperator {
		return errors.New("invalid filter operator, must be one of 'lt', 'gt', 'lte', 'gte', 'eq'")
	}

	if _, err := strconv.Atoi(strings.TrimSpace(value)); err != nil {
		return errors.New("invalid filter value, must be a number")
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

	err = validateSortOrder(fReq.SortOrder)
	if err != nil {
		return err
	}

	err = validateSortBy(fReq.SortBy)
	if err != nil {
		return err
	}

	err = validateFilter(fReq.Filter)
	if err != nil {
		return err
	}

	return nil
}
