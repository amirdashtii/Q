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
	ticketGroup.GET("/:id", h.GetTicketsByIDHandler)
	ticketGroup.GET("/", h.GetAllTicketsHandler)
	ticketGroup.POST("/", h.ReserveTicketHandler)
	ticketGroup.POST("/cancel/:id", h.CancelTicketHandler)
}

func (h *TicketHandler) ReserveTicketHandler(c echo.Context) error {

	userIDStr := c.Get("id").(string)
	var tickets models.Tickets
	var reservationRequest models.ReservationRequest
	var passengerIDs []uuid.UUID

	if err := c.Bind(&reservationRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	if err := validators.ReservationValidation(&reservationRequest, userIDStr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	tickets.UserID = userID

	for _, passengerIDStr := range reservationRequest.PassengerIDs {
		passengerID, err := uuid.Parse(passengerIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": err.Error(),
			})
		}
		passengerIDs = append(passengerIDs, passengerID)
	}

	if err := h.svc.CreateReservation(reservationRequest.FlightID, passengerIDs, &tickets); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"tickets": tickets})
}

func (h *TicketHandler) GetTicketsByIDHandler(c echo.Context) error {

	var tickets models.Tickets

	userIDStr := c.Get("id").(string)
	ticketsIDStr := c.Param("id")

	if err := validators.IDValidation(map[string]string{"user id": userIDStr, "tickets id": ticketsIDStr}); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	tickets.UserID = userID

	ticketID, err := uuid.Parse(ticketsIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	tickets.ID = ticketID

	if err := h.svc.GetTicketsByID(&tickets); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"tickets": tickets})
}

func (h *TicketHandler) GetAllTicketsHandler(c echo.Context) error {

	var allTickets []models.Tickets

	userIDStr := c.Get("id").(string)

	if err := validators.IDValidation(map[string]string{"user id": userIDStr}); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if err := h.svc.GetAllTickets(&userID, &allTickets); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"all tickets": allTickets})
}

func (h *TicketHandler) CancelTicketHandler(c echo.Context) error {
	var tickets models.Tickets

	userIDStr := c.Get("id").(string)
	ticketsIDStr := c.Param("id")

	if err := validators.IDValidation(map[string]string{"user id": userIDStr, "tickets id": ticketsIDStr}); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	tickets.UserID = userID

	ticketID, err := uuid.Parse(ticketsIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	tickets.ID = ticketID

	if err := h.svc.CancelTicket(&tickets); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Ticket canceled successfully"})
}
