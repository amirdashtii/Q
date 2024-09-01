package controller

import (
	"net/http"

	"github.com/amirdashtii/Q/flight-ticket-service/controller/middleware"
	"github.com/amirdashtii/Q/flight-ticket-service/controller/validators"
	"github.com/amirdashtii/Q/flight-ticket-service/models"
	"github.com/amirdashtii/Q/flight-ticket-service/ports"
	"github.com/amirdashtii/Q/flight-ticket-service/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TicketHandler struct {
	svc ports.TicketServiceContract
}

func NewTicketHandler() *TicketHandler {
	svc := service.NewTicketService()
	return &TicketHandler{
		svc: svc,
	}
}
func AddTicketServiceRoutes(e *echo.Echo) {
	h := NewTicketHandler()

	// Ticket Routes
	ticketGroup := e.Group("/tickets")
	ticketGroup.Use(middleware.JwtMiddleware)
	ticketGroup.POST("/reserve", h.ReserveTicketHandler) // رزرو بلیت
	// ticketGroup.POST("/cancel", h.CancelTicketHandler)          // لغو رزرو بلیت
	// ticketGroup.POST("/pay", h.PayTicketHandler)                // پرداخت بلیت
	// ticketGroup.PATCH("/update", h.UpdateTicketHandler)         // تغییر یا به روزرسانی اطلاعات بلیت

}

func (h *TicketHandler) ReserveTicketHandler(c echo.Context) error {

	userIDStr := c.Get("id").(string)
	var reservation models.TicketReservation
	var reservationRequest models.TicketReservationRequest
	var passengerIDs []uuid.UUID

	if err := c.Bind(&reservationRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	if err := validators.ReservationValidation(&reservationRequest, userIDStr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	reservation.Status = "reserved"
	if reservationUserID, err := uuid.Parse(userIDStr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	} else {
		reservation.UserID = reservationUserID
	}

	if reservationFlightID, err := uuid.Parse(reservationRequest.FlightID); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	} else {
		reservation.FlightID = reservationFlightID
	}

	for _, passengerIDStr := range reservationRequest.PassengerIDs {
		passengerID, err := uuid.Parse(passengerIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": err.Error(),
			})
		}
		passengerIDs = append(passengerIDs, passengerID)
	}

	if err := h.svc.CreateReservation(&passengerIDs, &reservation); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to reserve ticket"})
	}

	return c.JSON(http.StatusOK, echo.Map{"reservation": reservation})
}
