package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/google/uuid"
)

func (p *Postgres) RegisterUser(user *models.User) error {
	result := p.db.Create(user)
	return result.Error
}

func (p *Postgres) LoginUser(email string) (*models.User, error) {

	var fundedUser models.User
	if err := p.db.Where("email = ? ", email).First(&fundedUser).Error; err != nil {
		return nil, err
	}
	return &fundedUser, nil
}

func (p *Postgres) GetUsers(users *[]models.User) error {
	result := p.db.Find(users)
	return result.Error
}

func (p *Postgres) GetUserById(user *models.User) error {
	result := p.db.First(user)
	return result.Error
}

func (p *Postgres) UpdateUserById(id *uuid.UUID, updateItem map[string]interface{}) error {
	var user models.User
	user.ID = *id
	result := p.db.Model(user).Updates(updateItem)
	return result.Error
}

func (p *Postgres) DeleteUser(user *models.User) error {
	result := p.db.Delete(user)
	return result.Error
}

func (r *RedisDB) AddToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.client.Set(ctx, token, true, 0).Err()
	if err != nil {
		fmt.Printf("Error connecting to Redis: %v\n", err)
		return err
	}
	return nil
}

func (r *RedisDB) TokenReceiver(token string) (string, error) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, token).Result()

	return val, err
}
