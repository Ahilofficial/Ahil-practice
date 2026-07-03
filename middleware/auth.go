package middleware

import (
	"backend_institutions/helper"
	"fmt"
	"os"
	"strings"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return helper.Error(c, 401, "Authorization header is required")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return helper.Error(c, 401, "Authorization header format must be Bearer {token}")
		}

		tokenString := parts[1]
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "supersecretkey"
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			return helper.Error(c, 401, "Invalid or expired token: "+err.Error())
		}

		if !token.Valid {
			return helper.Error(c, 401, "Invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return helper.Error(c, 401, "Invalid token claims")
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("email", claims["email"])

		return c.Next()
	}
}
