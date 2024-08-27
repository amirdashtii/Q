package controller

import (
	"net/http"

	"github.com/amirdashtii/Q/flight-provider-service/controller/validators"
	"github.com/amirdashtii/Q/flight-provider-service/models"
	"github.com/amirdashtii/Q/flight-provider-service/ports"
	"github.com/amirdashtii/Q/flight-provider-service/service"
	"github.com/google/uuid"
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
	e.GET("/flights/:id", h.GetFlightByIDHandler)
	e.PATCH("/flights/:id/reserve", h.DecreaseFlightCapacityHandler)
	e.PATCH("/flights/:id/cancel", h.IncreaseFlightCapacityHandler)
}

func (h *FlightHandler) GetFlightsHandler(c echo.Context) error {
	err := h.svc.GenerateRandomFlightsForNext30Days()
	if err != nil {
	}

	var flightReq models.FlightReq
	err = c.Bind(&flightReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	err = validators.ValidateFlightParam(&flightReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
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

func (h *FlightHandler) GetFlightByIDHandler(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
	}

	var flight models.Flight
	flight.ID = id
	err = h.svc.GetFlightByID(&flight)

	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Flight not found"})
	}

	return c.JSON(http.StatusOK, flight)
}

func (h *FlightHandler) DecreaseFlightCapacityHandler(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
	}

	var req struct {
		Seats int `json:"seats"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	err = h.svc.DecreaseFlightCapacity(id, req.Seats)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Seats successfully reserved"})
}

func (h *FlightHandler) IncreaseFlightCapacityHandler(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
	}

	var req struct {
		Seats int `json:"seats"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	err = h.svc.IncreaseFlightCapacity(id, req.Seats)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Seats successfully released"})
}
