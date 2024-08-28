package controller

import (
	"net/http"

	"github.com/amirdashtii/Q/flight-ticket-service/controller/validators"
	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/amirdashtii/Q/flight-ticket-service/ports"
	"github.com/amirdashtii/Q/flight-ticket-service/service"
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
}

func (h *FlightHandler) GetFlightsHandler(c echo.Context) error {
	var flightReq models.FlightSearchRequest

	flightReq.Source = c.QueryParam("source")
	flightReq.Destination = c.QueryParam("destination")
	flightReq.DepartureDate = c.QueryParam("departure_date")

	flightReq.SortBy = c.QueryParam("sort_by")
	flightReq.SortOrder = c.QueryParam("order")
	flightReq.Filter = c.QueryParam("filter_by")

	err := validators.ValidateFlightParam(&flightReq)
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
	id := c.Param("id")

	err := validators.ValidateID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	var flight models.Flight
	err = h.svc.GetFlightByID(&id, &flight)

	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Flight not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"flight": flight,
	})

}
