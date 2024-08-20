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

	err = u.redis.AddToken(token)

	if err != nil {
		return "", "", err
	}

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
