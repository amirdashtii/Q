package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/amirdashtii/Q/auth-service/models"
)

func (p *Postgres) AddUser(user *models.User) error {
	result := p.db.Create(user)
	return result.Error
}

func (p *Postgres) GetUsers(users *[]models.User) error {
	result := p.db.Find(users)
	return result.Error
}

func (p *Postgres) LoginUser(email string) (*models.User, error) {

	var fundedUser models.User
	if err := p.db.Where("email = ? ", email).First(&fundedUser).Error; err != nil {
		return nil, err
	}
	return &fundedUser, nil
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

func (r *RedisDB) RevokeToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.client.Set(ctx, token, false, 0).Err()
	return err
}

func (r *RedisDB) TokenReceiver(token string) (string, error) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, token).Result()

	return val, err
}
