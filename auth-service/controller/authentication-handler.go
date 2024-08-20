package controller

import (
	"net/http"

	"github.com/amirdashtii/Q/auth-service/controller/middleware"
	"github.com/amirdashtii/Q/auth-service/controller/validators"
	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/amirdashtii/Q/auth-service/ports"
	"github.com/amirdashtii/Q/auth-service/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthenticationHandler struct {
	svc ports.UserServiceContract
}

func NewAuthenticationHandler() *AuthenticationHandler {
	svc := service.NewAuthenticationService()
	return &AuthenticationHandler{
		svc: svc,
	}
}
func AddAuthServiceRoutes(e *echo.Echo) {
	h := NewAuthenticationHandler()
	// Authentication Routes
	e.POST("/auth/register", h.RegisterHandler)
	e.POST("/auth/login", h.LoginHandler)
	// e.POST("/auth/logout", h.LogoutHandler)
	// e.POST("/auth/refresh-token", h.RefreshTokenHandler)

	// User Routes
	// e.GET("/user/profile", h.GetUserProfileHandler)
	// e.PATCH("/user/profile", h.UpdateUserProfileHandler)
	// e.PATCH("/user/change-password", h.ChangePasswordHandler)

	// Admin Routes
	e.GET("/admin/users", h.GetUsersHandler, middleware.AuthMiddleware)
	e.GET("/admin/users/:user_id", h.GetUserByIdHandler, middleware.AuthMiddleware)
	e.PATCH("/admin/users/:user_id", h.UpdateUserByIdHandler, middleware.AuthMiddleware)
	e.PATCH("/admin/users/:user_id/promote", h.PromoteUserToAdminHandler, middleware.AuthMiddleware)
	e.PATCH("/admin/users/:user_id/deactivate", h.DeactivateUserHandler, middleware.AuthMiddleware)
	e.PATCH("/admin/users/:user_id/activate", h.ActivateUserHandler, middleware.AuthMiddleware)
	e.DELETE("/admin/users/:user_id", h.DeleteUserHandler, middleware.AuthMiddleware)

	// e.POST("/logout", h.logout)
	// e.POST("/create-admin", h.CreateAdmin)
	// e.GET("/is-admin/:id/:role", h.IsAdmin)
	// e.GET("/verify/:number/:id", h.Verify)
	// e.POST("/disable-user", h.DisableUser)
	// e.GET("/test", h.Test, middleware.AuthMiddleware)
}

func (h *AuthenticationHandler) RegisterHandler(c echo.Context) error {

	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	err := validators.RegisterValidation(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	err = h.svc.RegisterUser(&user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]models.User{
		"user": user,
	})
}

func (h *AuthenticationHandler) LoginHandler(c echo.Context) error {

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	err := validators.LoginValidation(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error()})
	}

	token, err := h.svc.LoginUser(&user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token})
}

func (h *AuthenticationHandler) GetUsersHandler(c echo.Context) error {
	id := c.Get("id").(string)
	currentUserID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	var users []models.User

	err = h.svc.GetUsers(&currentUserID, &users)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string][]models.User{
		"users": users,
	})
}

func (h *AuthenticationHandler) GetUserByIdHandler(c echo.Context) error {
	id := c.Get("id").(string)
	currentUserID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	var user models.User
	user.ID = userID
	err = h.svc.GetUserById(&currentUserID, &user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]models.User{
		"user": user,
	})
}

func (h *AuthenticationHandler) UpdateUserByIdHandler(c echo.Context) error {
	id := c.Get("id").(string)
	currentUserID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	err = validators.UpdateValidation(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	user.ID = userID
	err = h.svc.UpdateUserById(&currentUserID, &user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "user successfully updated",
	})
}

func (h *AuthenticationHandler) PromoteUserToAdminHandler(c echo.Context) error {
	id := c.Get("id").(string)
	currentUserID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Target user ID format is incorrect. Please check and try again.",
		})
	}

	err = h.svc.PromoteUserToAdmin(&currentUserID, &userID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to promote user to admin: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User successfully promoted to admin.",
	})
}

func (h *AuthenticationHandler) DeactivateUserHandler(c echo.Context) error {
	id := c.Get("id").(string)
	currentUserID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Target user ID format is incorrect. Please check and try again.",
		})
	}

	err = h.svc.DeactivateUser(&currentUserID, &userID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to deactivate user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User successfully deactivated.",
	})
}

func (h *AuthenticationHandler) ActivateUserHandler(c echo.Context) error {
	id := c.Get("id").(string)
	currentUserID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Target user ID format is incorrect. Please check and try again.",
		})
	}

	err = h.svc.ActivateUser(&currentUserID, &userID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to deactivate user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User successfully activated.",
	})
}

func (h *AuthenticationHandler) DeleteUserHandler(c echo.Context) error {
	id := c.Get("id").(string)
	currentUserID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Target user ID format is incorrect. Please check and try again.",
		})
	}

	var user models.User
	user.ID = userID
	err = h.svc.DeleteUser(&currentUserID, &user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User successfully deleted.",
	})
}
