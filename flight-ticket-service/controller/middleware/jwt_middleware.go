package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/amirdashtii/Q/auth-service/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

func JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Missing or malformed JWT"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token format"})
		}

		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid or expired token"})
		}

		claims := token.Claims.(*models.Claims)
		c.Set("id", claims.ID)
		c.Set("role", claims.Role)

		return next(c)
	}
}
