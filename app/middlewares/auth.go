package middlewares

import (
	"fmt"
	"net/http"

	cfg "hintword.com/api/app/configs"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// RequireLoggedIn ensures access only to login users by checking for token presence and validity
func RequireLoggedIn() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(cfg.GetConfig().JWTAccessSecret),
		ErrorHandler: jwtError,
	})
}

func OptionalAuth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(cfg.GetConfig().JWTAccessSecret), // Replace with your actual secret key
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Handle authentication errors (optional)
			fmt.Println("Authentication error:", err)
			return c.Next()
		},
	})
}

// TODO, Research how fiber checks for expired token?
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		var errorList []*fiber.Error
		errorList = append(
			errorList,
			&fiber.Error{
				Code:    fiber.StatusUnauthorized,
				Message: "Missing or Malformed Authentication AccessToken",
			},
		)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"errors": errorList})
	}

	var errorList []*fiber.Error
	errorList = append(
		errorList,
		&fiber.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid or Expired Authentication AccessToken",
		},
	)
	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"errors": errorList})
}
