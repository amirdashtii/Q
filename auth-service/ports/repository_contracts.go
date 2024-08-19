package ports

import (
	"github.com/amirdashtii/Q/auth-service/models"
)

type UserRepositoryContracts interface {
	AddUser(user *models.User) error
	GetUsers(users *[]models.User) error
}
