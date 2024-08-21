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
	authGroup := e.Group("/auth")
	authGroup.POST("/register", h.RegisterHandler)
	authGroup.POST("/login", h.LoginHandler)
	authGroup.POST("/logout", h.LogoutHandler)
	authGroup.POST("/refresh-token", h.RefreshTokenHandler)

	// User Routes
	userGroup := e.Group("/user")
	userGroup.Use(middleware.JwtMiddleware)
	userGroup.GET("/profile", h.GetUserProfileHandler)
	userGroup.PATCH("/profile", h.UpdateUserProfileHandler)
	userGroup.PATCH("/change-password", h.ChangePasswordHandler)

	// Admin Routes
	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.JwtMiddleware)
	adminGroup.GET("/users", h.GetUsersHandler)
	adminGroup.GET("/users/:user_id", h.GetUserByIdHandler)
	adminGroup.PATCH("/users/:user_id", h.UpdateUserByIdHandler)
	adminGroup.PATCH("/users/:user_id/promote", h.PromoteUserToAdminHandler)
	adminGroup.PATCH("/users/:user_id/deactivate", h.DeactivateUserHandler)
	adminGroup.PATCH("/users/:user_id/activate", h.ActivateUserHandler)
	adminGroup.DELETE("/users/:user_id", h.DeleteUserHandler)
}

func (h *AuthenticationHandler) RegisterHandler(c echo.Context) error {

	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	err := validators.RegisterValidation(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	err = h.svc.RegisterUser(&user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
}

func (h *AuthenticationHandler) LoginHandler(c echo.Context) error {

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	err := validators.LoginValidation(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error()})
	}

	accessToken, refreshToken, err := h.svc.LoginUser(&user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthenticationHandler) LogoutHandler(c echo.Context) error {

	token := c.Request().Header.Get("Authorization")[7:]

	err := validators.TokenValidation(token)
	if err != nil {
		return c.JSON(http.StatusForbidden, err.Error())
	}

	err = h.svc.Logout(token)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, "logout successful")
}

func (h *AuthenticationHandler) RefreshTokenHandler(c echo.Context) error {

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	newAccessToken, err := h.svc.RefreshToken(req.RefreshToken)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{"access_token": newAccessToken})
}

func (h *AuthenticationHandler) GetUserProfileHandler(c echo.Context) error {
	id := c.Get("id").(string)
	currentUserID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var user models.User
	user.ID = currentUserID

	err = h.svc.GetUserProfile(&user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
}

func (h *AuthenticationHandler) UpdateUserProfileHandler(c echo.Context) error {
	id := c.Get("id").(string)
	currentUserID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request body",
		})
	}
	user.ID = currentUserID

	err = validators.UpdateValidation(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	err = h.svc.UpdateUserProfile(&user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user successfully updated",
	})
}

func (h *AuthenticationHandler) ChangePasswordHandler(c echo.Context) error {
	id := c.Get("id").(string)
	currentUserID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request body",
		})
	}
	user.ID = currentUserID

	err = validators.PasswordValidation(user.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	err = h.svc.ChangePassword(&user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "password successfully changed",
	})
}

func (h *AuthenticationHandler) GetUsersHandler(c echo.Context) error {

	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Access denied"})
	}

	var users []models.User
	err := h.svc.GetUsers(&users)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"users": users,
	})
}

func (h *AuthenticationHandler) GetUserByIdHandler(c echo.Context) error {

	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Access denied"})
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var user models.User
	user.ID = userID
	err = h.svc.GetUserById(&user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
}

func (h *AuthenticationHandler) UpdateUserByIdHandler(c echo.Context) error {

	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Access denied"})
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request body",
		})
	}

	err = validators.UpdateValidation(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	user.ID = userID
	err = h.svc.UpdateUserById(&user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user successfully updated",
	})
}

func (h *AuthenticationHandler) PromoteUserToAdminHandler(c echo.Context) error {

	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Access denied"})
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Target user ID format is incorrect. Please check and try again.",
		})
	}

	err = h.svc.PromoteUserToAdmin(&userID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Failed to promote user to admin: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "User successfully promoted to admin.",
	})
}

func (h *AuthenticationHandler) DeactivateUserHandler(c echo.Context) error {

	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Access denied"})
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Target user ID format is incorrect. Please check and try again.",
		})
	}

	err = h.svc.DeactivateUser(&userID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Failed to deactivate user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "User successfully deactivated.",
	})
}

func (h *AuthenticationHandler) ActivateUserHandler(c echo.Context) error {

	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Access denied"})
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Target user ID format is incorrect. Please check and try again.",
		})
	}

	err = h.svc.ActivateUser(&userID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Failed to deactivate user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "User successfully activated.",
	})
}

func (h *AuthenticationHandler) DeleteUserHandler(c echo.Context) error {

	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Access denied"})
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Target user ID format is incorrect. Please check and try again.",
		})
	}

	var user models.User
	user.ID = userID
	err = h.svc.DeleteUser(&user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "User successfully deleted.",
	})
}
