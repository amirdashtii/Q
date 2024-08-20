package ports

import (
	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/google/uuid"
)

type UserRepositoryContracts interface {
	RegisterUser(user *models.User) error
	LoginUser(email string) (*models.User, error)
	GetUsers(users *[]models.User) error
	GetUserById(user *models.User) error
	UpdateUserById(id *uuid.UUID, updateItem map[string]interface{}) error
	DeleteUser(user *models.User) error
}

type InMemoryRespositoryContracts interface {
	AddToken(token string) error
	TokenReceiver(token string) (string, error)
}
