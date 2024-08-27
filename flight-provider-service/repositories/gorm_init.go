package repositories

import (
	"fmt"
	"log"
	"os"

	"github.com/amirdashtii/Q/flight-provider-service/models"
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

	err = database.AutoMigrate(&models.Flight{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	return database, nil
}
func NewPostgres() *Postgres {
	db, _ := GormInit()
	return &Postgres{
		db: db,
	}
}
