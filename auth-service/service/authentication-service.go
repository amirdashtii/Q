package service

import (
	"errors"

	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/amirdashtii/Q/auth-service/ports"
	"github.com/amirdashtii/Q/auth-service/repositories"
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

func (u *AuthenticationService) LoginUser(user *models.User) (string, string, error) {

	inpPassword := user.Password
	err := u.db.LoginUser(user)
	if err != nil {
		return "", "", err
	}

	if user.Disabled {
		return "", "", errors.New("your account is temporarily disabled by an admin")
	}

	decodedFoundedPassword := CheckPasswordHash(inpPassword, user.Password)

	if !decodedFoundedPassword {
		err := errors.New("email or password mismatch")
		return "", "", err
	}

	accessToken, err := GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(user.ID.String())
	if err != nil {
		return "", "", err
	}

	err = u.redis.AddToken(refreshToken, user.ID.String())
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *AuthenticationService) Logout(token string) error {
	err := u.redis.RevokeToken(token)

	return err
}

func (u *AuthenticationService) RefreshToken(refreshToken string) (string, error) {
	claims, err := ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	storedToken, err := u.redis.ReceiverToken(claims)
	if err != nil || storedToken != refreshToken {
		return "", err
	}

	newAccessToken, err := GenerateAccessToken(claims.ID, claims.Role)
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}

func (u *AuthenticationService) GetUserProfile(user *models.User) error {
	return u.db.GetUserById(user)
}

func (u *AuthenticationService) UpdateUserProfile(user *models.User) error {

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

	if user.PhoneNumber != nil {
		updateItem["phone_number"] = user.PhoneNumber
	}

	return u.db.UpdateUserById(&user.ID, updateItem)
}

func (u *AuthenticationService) ChangePassword(user *models.User) error {

	updateItem := make(map[string]interface{})

	hashPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	updateItem["password"] = hashPassword

	return u.db.UpdateUserById(&user.ID, updateItem)
}

func (u *AuthenticationService) GetUsers(users *[]models.User) error {
	return u.db.GetUsers(users)
}

func (u *AuthenticationService) GetUserById(user *models.User) error {
	return u.db.GetUserById(user)
}

func (u *AuthenticationService) UpdateUserById(user *models.User) error {

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

func (u *AuthenticationService) PromoteUserToAdmin(userID *uuid.UUID) error {

	var user models.User
	user.ID = *userID

	updateItem := make(map[string]interface{})

	updateItem["role"] = "admin"

	return u.db.UpdateUserById(&user.ID, updateItem)
}

func (u *AuthenticationService) DeactivateUser(userID *uuid.UUID) error {

	var user models.User
	user.ID = *userID

	updateItem := make(map[string]interface{})

	updateItem["disabled"] = true

	return u.db.UpdateUserById(&user.ID, updateItem)
}

func (u *AuthenticationService) ActivateUser(userID *uuid.UUID) error {

	var user models.User
	user.ID = *userID

	updateItem := make(map[string]interface{})

	updateItem["disabled"] = false

	return u.db.UpdateUserById(&user.ID, updateItem)
}

func (u *AuthenticationService) DeleteUser(user *models.User) error {
	return u.db.DeleteUser(user)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
