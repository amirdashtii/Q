package controller

import (
	"net/http"

	"github.com/amirdashtii/Q/auth-service/controller/validators"
	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/amirdashtii/Q/auth-service/ports"
	"github.com/amirdashtii/Q/auth-service/service"
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

	e.POST("/api/auth/users", h.CreateUserHandler)
	e.GET("/api/auth/users", h.GetUsersHandler)
	// e.POST("/api/auth/login", h.LoginHandler)
	// e.POST("/api/auth/logout", h.LogoutHandler)
	// e.POST("/api/auth/forgot-password", h.ForgotPasswordHandler)
	// e.POST("/api/auth/reset-password", h.ResetPasswordHandler)
	// e.GET("/api/auth/verify-account", h.VerifyAccountHandler)
	// e.PUT("/api/auth/update-profile", h.UpdateProfileHandler)
	// e.GET("/api/auth/profile", h.GetUserProfileHandler)
	// e.POST("/api/auth/refresh-token", h.RefreshTokenHandler)

	// e.POST("/login", h.loginHandler)
	// e.POST("/register", h.register)
	// e.POST("/logout", h.logout)
	// e.POST("/create-admin", h.CreateAdmin)
	// e.GET("/is-admin/:id/:role", h.IsAdmin)
	// e.GET("/verify/:number/:id", h.Verify)
	// e.POST("/disable-user", h.DisableUser)
	// e.GET("/test", h.Test, middleware.AuthMiddleware)
}

func (h *AuthenticationHandler) CreateUserHandler(c echo.Context) error {

	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	err := validators.RegisterValidation(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = h.svc.AddUser(&user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]models.User{
		"user": user,
	})
}

func (h *AuthenticationHandler) GetUsersHandler(c echo.Context) error {

	var users []models.User

	err := h.svc.GetUsers(&users)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string][]models.User{
		"users": users,
	})
}
