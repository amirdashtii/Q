package service

import (
	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/amirdashtii/Q/auth-service/ports"
	"github.com/amirdashtii/Q/auth-service/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	db ports.UserRepositoryContracts
}

func NewAuthenticationService() *AuthenticationService {
	db := repositories.NewPostgres()

	return &AuthenticationService{
		db: db,
	}
}

func (u *AuthenticationService) AddUser(user *models.User) error {

	user.Password, _ = HashPassword(user.Password)
	user.Role = "user"
	user.Disabled = false

	err := u.db.AddUser(user)
	return err
}

func (u *AuthenticationService) GetUsers(users *[]models.User) error {

	err := u.db.GetUsers(users)
	return err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
