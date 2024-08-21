package ports

import (
	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/google/uuid"
)

type UserServiceContract interface {
	RegisterUser(user *models.User) error
	LoginUser(user *models.User) (string, string, error)
	Logout(token string) error
	RefreshToken(refreshToken string) (string, error)

	GetUserProfile(user *models.User) error
	UpdateUserProfile(user *models.User) error
	ChangePassword(user *models.User) error

	GetUsers(users *[]models.User) error
	GetUserById(user *models.User) error
	UpdateUserById(user *models.User) error
	PromoteUserToAdmin(userID *uuid.UUID) error
	DeactivateUser(userID *uuid.UUID) error
	ActivateUser(userID *uuid.UUID) error
	DeleteUser(user *models.User) error
}
