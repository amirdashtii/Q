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

type PaymentHandler struct {
	svc ports.PaymentServiceContract
}

func NewPaymentHandler() *PaymentHandler {
	svc := service.NewPaymentService()
	return &PaymentHandler{
		svc: svc,
	}
}

func AddPaymentServiceRoutes(e *echo.Echo) {
	h := NewPaymentHandler()

	// Ticket Routes
	paymentGroup := e.Group("/payment")
	paymentGroup.Use(middleware.JwtMiddleware)
	paymentGroup.POST("/pay", h.PayTicketHandler)
	e.POST("/payment/success", h.PaymentSuccessTicketHandler)
}

func (h *PaymentHandler) PayTicketHandler(c echo.Context) error {

	type Re struct {
		TicketID       string `json:"tickets_id"`
		PaymentGateway string `json:"payment_gateway"`
	}

	var re Re
	var ticket models.Tickets

	userIDStr := c.Get("id").(string)

	if err := c.Bind(&re); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}
	ticketIDStr := re.TicketID

	if err := validators.IDValidation(map[string]string{"user id": userIDStr, "ticket id": ticketIDStr}); err != nil {
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
	ticket.UserID = userID

	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	ticket.ID = ticketID

	paymentLink, err := h.svc.PayTicket(&ticket, re.PaymentGateway)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"payment link": paymentLink})
}

func (h *PaymentHandler) PaymentSuccessTicketHandler(c echo.Context) error {

	var receivedPaymentRequest models.ReceivedPaymentRequest

	if err := c.Bind(&receivedPaymentRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	transaction, err := h.svc.VerifyTransaction(&receivedPaymentRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"transaction": transaction})
}
