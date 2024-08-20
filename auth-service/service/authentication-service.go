package service

import (
	"errors"
	"os"
	"time"

	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/amirdashtii/Q/auth-service/ports"
	"github.com/amirdashtii/Q/auth-service/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	db    ports.UserRepositoryContracts
	redis ports.InMemoryRespositoryContracts
}

func NewAuthenticationService() *AuthenticationService {
	db := repositories.NewPostgres()
	redis := repositories.RedisInit()

	return &AuthenticationService{
		db:    db,
		redis: redis,
	}
}

func (u *AuthenticationService) RegisterUser(user *models.User) error {

	user.Password, _ = HashPassword(user.Password)
	user.Role = "user"
	user.Disabled = false

	err := u.db.RegisterUser(user)
	return err
}

func (u *AuthenticationService) LoginUser(user *models.User) (string, error) {

	foundedUser, err := u.db.LoginUser(user.Email)
	if err != nil {
		return "", err
	}

	if foundedUser.Disabled {
		return "", errors.New("your account is temporarily disabled by an admin")
	}

	decodedFoundedPassword := CheckPasswordHash(user.Password, foundedUser.Password)

	if !decodedFoundedPassword {
		err := errors.New("email or password mismatch")
		return "", err
	}

	token, err := GenerateToken(foundedUser.DBModel.ID)
	if err != nil {
		return "", err
	}

	err = u.redis.AddToken(token)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *AuthenticationService) GetUsers(currentUserID *uuid.UUID, users *[]models.User) error {

	err := u.IsAdmin(*currentUserID)
	if err != nil {
		return err
	}

	return u.db.GetUsers(users)
}

func (u *AuthenticationService) GetUserById(currentUserID *uuid.UUID, user *models.User) error {

	err := u.IsAdmin(*currentUserID)
	if err != nil {
		return err
	}

	return u.db.GetUserById(user)
}

func (u *AuthenticationService) UpdateUserById(currentUserID *uuid.UUID, user *models.User) error {
	err := u.IsAdmin(*currentUserID)
	if err != nil {
		return err
	}

	updateItem := make(map[string]interface{})

	if user.FirstName != "" {
		updateItem["first_name"] = user.FirstName
	}

	if user.LastName != "" {
		updateItem["last_name"] = user.LastName
	}

	if user.Email != "" {
		updateItem["email"] = user.Email
	}

	if user.Password != "" {
		hashPassword, err := HashPassword(user.Password)
		if err != nil {
			return err
		}
		updateItem["password"] = hashPassword
	}

	if user.PhoneNumber != nil {
		updateItem["phone_number"] = user.Password
	}

	return u.db.UpdateUserById(&user.ID, updateItem)
}

func (u *AuthenticationService) PromoteUserToAdmin(currentUserID, userID *uuid.UUID) error {

	err := u.IsAdmin(*currentUserID)
	if err != nil {
		return err
	}
	var user models.User
	user.ID = *userID

	updateItem := make(map[string]interface{})

	updateItem["role"] = "admin"

	return u.db.UpdateUserById(&user.ID, updateItem)
}

func (u *AuthenticationService) DeactivateUser(currentUserID, userID *uuid.UUID) error {

	err := u.IsAdmin(*currentUserID)
	if err != nil {
		return err
	}
	var user models.User
	user.ID = *userID

	updateItem := make(map[string]interface{})

	updateItem["disabled"] = true

	return u.db.UpdateUserById(&user.ID, updateItem)
}

func (u *AuthenticationService) ActivateUser(currentUserID, userID *uuid.UUID) error {

	err := u.IsAdmin(*currentUserID)
	if err != nil {
		return err
	}
	var user models.User
	user.ID = *userID

	updateItem := make(map[string]interface{})

	updateItem["disabled"] = false

	return u.db.UpdateUserById(&user.ID, updateItem)
}

func (u *AuthenticationService) DeleteUser(currentUserID *uuid.UUID, user *models.User) error {
	err := u.IsAdmin(*currentUserID)
	if err != nil {
		return err
	}

	return u.db.DeleteUser(user)
}

func (u *AuthenticationService) IsAdmin(id uuid.UUID) error {
	var user models.User
	user.ID = id
	err := u.db.GetUserById(&user)
	if err != nil {
		return err
	}

	if user.Role != "admin" {
		return errors.New("access denied. admin privileges are required to perform this action")
	}

	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(id uuid.UUID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
