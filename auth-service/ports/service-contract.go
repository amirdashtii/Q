package ports

import (
	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/google/uuid"
)

type UserServiceContract interface {
	RegisterUser(user *models.User) error
	LoginUser(user *models.User) (string, error)
	GetUsers(currentUserID *uuid.UUID, users *[]models.User) error
	GetUserById(currentUserID *uuid.UUID, user *models.User) error
	UpdateUserById(currentUserID *uuid.UUID, user *models.User) error
	PromoteUserToAdmin(currentUserID, userID *uuid.UUID) error
	DeactivateUser(currentUserID, userID *uuid.UUID) error
	ActivateUser(currentUserID, userID *uuid.UUID) error
	DeleteUser(currentUserID *uuid.UUID, user *models.User) error
}
