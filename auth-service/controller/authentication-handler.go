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

	e.POST("/users", h.CreateUserHandler)
	e.GET("/users", h.GetUsersHandler)
	e.POST("/login", h.LoginHandler)
	// e.POST("/logout", h.LogoutHandler)
	// e.POST("/forgot-password", h.ForgotPasswordHandler)
	// e.POST("/reset-password", h.ResetPasswordHandler)
	// e.GET("/verify-account", h.VerifyAccountHandler)
	// e.PUT("/update-profile", h.UpdateProfileHandler)
	// e.GET("/profile", h.GetUserProfileHandler)
	// e.POST("/refresh-token", h.RefreshTokenHandler)

	// e.POST("/login", h.loginHandler)
	// e.POST("/register", h.register) don
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

func (h *AuthenticationHandler) LoginHandler(c echo.Context) error {

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"Error": "Invalid request body",
		})
	}

	err := validators.LoginValidation(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	token, id, err := h.svc.LoginUser(&user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"Error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token, "id": id})
}
