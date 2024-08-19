package service

import (
	"errors"
	"time"

	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/amirdashtii/Q/auth-service/ports"
	"github.com/amirdashtii/Q/auth-service/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (u *AuthenticationService) LoginUser(user *models.User) (string, string, error) {

	foundedUser, err := u.db.LoginUser(user.Email)
	if err != nil {
		return "", "", err
	}
	if foundedUser.Disabled {
		return "", "", errors.New("your account is temporarily disabled by an admin")
	}

	decodedFoundedPassword := CheckPasswordHash(user.Password, foundedUser.Password)

	if !decodedFoundedPassword {
		err := errors.New("email or password mismatch")
		return "", "", err
	}

	token, err := GenerateToken(foundedUser.DBModel.ID)
	if err != nil {
		return "", "", err
	}

	// err = u.redis.AddToken(token)

	// if err != nil {
	// 	return "", "", err
	// }

	return token, foundedUser.ID.String(), nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Claims struct {
	ID uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}
type Credentials struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

var jwtKey = []byte("my_secret_key")

func GenerateToken(id uuid.UUID) (string, error) {
	var creds Credentials
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		ID: creds.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
