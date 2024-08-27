package controller

import (
	"net/http"

	"github.com/amirdashtii/Q/flight-provider-service/controller/validators"
	"github.com/amirdashtii/Q/flight-provider-service/models"
	"github.com/amirdashtii/Q/flight-provider-service/ports"
	"github.com/amirdashtii/Q/flight-provider-service/service"
	"github.com/labstack/echo/v4"
)

type FlightHandler struct {
	svc ports.FlightServiceContract
}

func NewFlightHandler() *FlightHandler {
	svc := service.NewFlightService()
	return &FlightHandler{
		svc: svc,
	}
}
func AddFlightServiceRoutes(e *echo.Echo) {
	h := NewFlightHandler()

	e.GET("/flights", h.GetFlightsHandler)

}

func (h *FlightHandler) GetFlightsHandler(c echo.Context) error {
	err := h.svc.GenerateRandomFlightsForNext30Days()
	if err != nil {
	}

	var flightReq models.FlightReq
	err = c.Bind(&flightReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = validators.ValidateFlightParam(&flightReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var flights []models.Flight
	err = h.svc.GetFlights(&flightReq, &flights)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"flights": flights,
	})
}
