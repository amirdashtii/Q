package repositories

import (
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