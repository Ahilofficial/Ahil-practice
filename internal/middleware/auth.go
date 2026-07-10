package middleware

import (
	"backend_institutions/internal/helper"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
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

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "supersecretkey"
		}

		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return helper.Error(c, 401, "Invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return helper.Error(c, 401, "Invalid token claims")
		}

		// Check user_id exists
		userID, ok := claims["user_id"]
		if !ok || userID == nil {
			return helper.Error(c, 401, "user_id not found in token")
		}

		c.Locals("user_id", userID)

		return c.Next()
	}
}

func OptionalAuth() fiber.Handler {
	return func(c fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			c.Locals("user_id", nil)
			return c.Next()
		}

		tokensplit := strings.Split(authHeader, " ")
		if len(tokensplit) != 2 || strings.ToLower(tokensplit[0]) != "bearer" {
			c.Locals("user_id", nil)
			return c.Next()
		}
		JwtSecret := os.Getenv("JWT_SECRET")
		if JwtSecret == "" {
			JwtSecret = "supersecretkey"
		}
		token, err := jwt.Parse(tokensplit[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(JwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.Locals("user_id", nil)
			return c.Next()
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Locals("user_id", nil)
			return c.Next()
		}
		userID, ok := claims["user_id"]
		if !ok || userID == nil {
			c.Locals("user_id", nil)
			return c.Next()
		}
		c.Locals("user_id", userID)
		return c.Next()
	}
}
