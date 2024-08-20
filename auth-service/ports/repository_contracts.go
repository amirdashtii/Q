package ports

import (
	"github.com/amirdashtii/Q/auth-service/models"
)

type UserRepositoryContracts interface {
	AddUser(user *models.User) error
	GetUsers(users *[]models.User) error
	LoginUser(email string) (*models.User, error)
}

type InMemoryRespositoryContracts interface {
	AddToken(token string) error
	// RevokeToken(token string) error
	TokenReceiver(token string) (string, error)
}