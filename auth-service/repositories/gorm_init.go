package repositories

import (
	"fmt"
	"log"
	"os"

	"github.com/amirdashtii/Q/auth-service/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func GormInit() (*gorm.DB, error) {

	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB_NAME")
	port := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, user, password, dbName, port)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = database.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	Initialize(database)

	return database, nil
}
func NewPostgres() *Postgres {
	db, _ := GormInit()
	return &Postgres{
		db: db,
	}
}

func Initialize(db *gorm.DB) {
	var user models.User
	if err := db.Where("role = ?", "admin").First(&user).Error; err != nil {

		if err == gorm.ErrRecordNotFound {

			admin := models.User{
				Email:    "admin@domain.com",
				Password: "Passw0rd",
				Role:     "admin",
				Disabled: false,
			}
			admin.Password, _ = HashPassword(admin.Password)

			result := db.Create(&admin)

			if result.Error != nil {
				log.Fatalf("Failed to create admin user: %v", result.Error)
			}
		}
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
