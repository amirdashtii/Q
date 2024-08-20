package ports

import (
	"github.com/amirdashtii/Q/auth-service/models"
)

type UserServiceContract interface {
	AddUser(user *models.User) error
	GetUsers(users *[]models.User) error
	LoginUser(user *models.User) (string, string, error)
	// PromoteUserToAdmin(id *string) error
}

type InMemoryServiceContracts interface {
	AddToken(token string) error
	// RevokeToken(token string) *redis.StatusCmd
	TokenReceiver() (string, error)
}