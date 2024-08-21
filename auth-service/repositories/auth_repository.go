package repositories

import (
	"context"
	"time"

	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/google/uuid"
)

func (p *Postgres) RegisterUser(user *models.User) error {
	result := p.db.Create(user)
	return result.Error
}

func (p *Postgres) LoginUser(user *models.User) error {
	result := p.db.Where("email = ? ", user.Email).First(user)
	return result.Error
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

func (r *RedisDB) AddToken(refreshToken, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.client.Set(ctx, id, refreshToken, 7*24*time.Hour).Err()
	if err != nil {
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

func (r *RedisDB) ReceiverToken(claims *models.Claims) (string, error) {

	ctx := context.Background()
	return r.client.Get(ctx, claims.ID).Result()
}
