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

	// Flight Routes
	flightGroup := e.Group("/flights")
	flightGroup.GET("", h.GetFlightsHandler)
	flightGroup.GET("/:id", h.GetFlightByIDHandler)

	// // User Reservation Routes
	// userReservationGroup := e.Group("/user")
	// userReservationGroup.Use(middleware.JwtMiddleware)
	// userReservationGroup.GET("/reservations", h.ListUserReservationsHandler) // مشاهده رزروهای یک کاربر

	// // Reservation Routes
	// reservationGroup := e.Group("/reservations")
	// reservationGroup.GET("/:id", h.GetTicketsByIDHandler)   // دریافت جزئیات یک رزرو

	// // Payment Status Routes
	// paymentGroup := e.Group("/payment-status")
	// paymentGroup.GET("", h.GetPaymentStatusHandler)             // بررسی وضعیت پرداخت

	// // Airline Routes
	// airlineGroup := e.Group("/airlines")
	// airlineGroup.GET("", h.ListAirlinesHandler)                 // لیست شرکت های هواپیمایی

}

func (h *FlightHandler) GetFlightsHandler(c echo.Context) error {

	var flightReq models.FlightSearchRequest
	var flights []models.ProviderFlight

	flightReq.Source = c.QueryParam("source")
	flightReq.Destination = c.QueryParam("destination")
	flightReq.DepartureDate = c.QueryParam("departure_date")

	flightReq.SortBy = c.QueryParam("sort_by")
	flightReq.SortOrder = c.QueryParam("order")
	flightReq.Filter = c.QueryParam("filter_by")

	if err := validators.ValidateFlightParam(&flightReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.svc.GetFlights(&flightReq, &flights); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"flights": flights,
	})
}

func (h *FlightHandler) GetFlightByIDHandler(c echo.Context) error {

	var flight models.ProviderFlight

	userIDStr := c.Param("id")

	if err := validators.IDValidation(map[string]string{"id": userIDStr}); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.svc.GetFlightByID(&userIDStr, &flight); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"flight": flight,
	})

}
