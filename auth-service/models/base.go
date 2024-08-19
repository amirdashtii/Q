package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DBModel struct {
	ID         uuid.UUID       `gorm:"primaryKey;type:uuid" json:"id"`
	CreatedAt  time.Time       `json:"created_at"`
	ModifiedAt *time.Time      `json:"modified_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at"`
}

func (d *DBModel) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.New()
	return nil
}
