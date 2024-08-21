package main

import (
	"log"
	"os"

	"github.com/amirdashtii/Q/auth-service/controller"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	e := echo.New()
	e.Use(middleware.Logger())
	controller.AddAuthServiceRoutes(e)
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
