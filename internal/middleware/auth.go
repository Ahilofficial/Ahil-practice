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

		val:=c.Locals("user_id", userID)
		if val==nil{
			return helper.Error(c, 401, "user_id not found in context")
		}

		// email is optional
		if email, ok := claims["email"]; ok {
			c.Locals("email", email)
		}

		return c.Next()
	}
}

