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

type PassengerHandler struct {
	svc ports.PassengerServiceContract
}

func NewPassengerHandler() *PassengerHandler {
	svc := service.NewPassengerService()
	return &PassengerHandler{
		svc: svc,
	}
}
func AddPassengerServiceRoutes(e *echo.Echo) {
	h := NewPassengerHandler()

	// Passenger Routes
	passengerGroup := e.Group("/user")
	passengerGroup.Use(middleware.JwtMiddleware)
	passengerGroup.POST("/passengers", h.CreatePassengerHandler)
	passengerGroup.GET("/passengers", h.GetPassengersHandler)
	passengerGroup.GET("/passengers/:id", h.GetPassengerByIDHandler)
	passengerGroup.PATCH("/passengers/:id", h.UpdatePassengerHandler)
	passengerGroup.DELETE("/passengers/:id", h.DeletePassengerHandler)
}

func (h *PassengerHandler) CreatePassengerHandler(c echo.Context) error {

	var passenger models.Passenger
	var passengerReq models.PassengerReq

	userIDStr := c.Get("id").(string)

	if err := c.Bind(&passengerReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if err := validators.PassengerValidation(passengerReq, userIDStr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if passengerUserID, err := uuid.Parse(userIDStr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	} else {
		passenger.UserID = passengerUserID
	}

	if err := h.svc.CreatePassenger(&passengerReq, &passenger); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"passenger": passenger,
	})
}

func (h *PassengerHandler) GetPassengersHandler(c echo.Context) error {

	var passengers []models.Passenger

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

	if err = h.svc.GetPassengers(&userID, &passengers); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"passengers": passengers,
	})
}

func (h *PassengerHandler) GetPassengerByIDHandler(c echo.Context) error {
	var passenger models.Passenger

	userIDStr := c.Get("id").(string)
	passengerIDStr := c.Param("id")

	if err := validators.IDValidation(map[string]string{"user id": userIDStr, "passenger id": passengerIDStr}); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if passengerUserID, err := uuid.Parse(userIDStr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	} else {
		passenger.UserID = passengerUserID
	}

	if passengerID, err := uuid.Parse(passengerIDStr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	} else {
		passenger.ID = passengerID
	}

	if err := h.svc.GetPassengerByID(&passenger); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Passenger not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"passenger": passenger,
	})

}

func (h *PassengerHandler) UpdatePassengerHandler(c echo.Context) error {

	var passengerReq models.PassengerReq

	userIDStr := c.Get("id").(string)
	passengerIDStr := c.Param("id")

	if err := c.Bind(&passengerReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if err := validators.PassengerUpdateValidation(passengerIDStr, userIDStr, passengerReq.DateOfBirth); err != nil {
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

	passengerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if err = h.svc.UpdatePassenger(&passengerID, &userID, &passengerReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "passenger successfully updated",
	})

}

func (h *PassengerHandler) DeletePassengerHandler(c echo.Context) error {

	var passenger models.Passenger

	userIDStr := c.Get("id").(string)
	passengerIDStr := c.Param("id")

	if err := validators.IDValidation(map[string]string{"user id": userIDStr, "passenger id": passengerIDStr}); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if passengerUserID, err := uuid.Parse(userIDStr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	} else {
		passenger.UserID = passengerUserID
	}

	if passengerID, err := uuid.Parse(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	} else {
		passenger.ID = passengerID
	}

	if err := h.svc.DeletePassenger(&passenger); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "passenger successfully Deleted",
	})

}
