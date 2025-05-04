package middlewares

import (
	"fmt"
	"net/http"

	"hintword.com/api/app/common/constants"
	cfg "hintword.com/api/app/configs"

	jwt "github.com/form3tech-oss/jwt-go"
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
				Message: "Missing or Malformed Authentication Token",
			},
		)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"errors": errorList})
	}

	var errorList []*fiber.Error
	errorList = append(
		errorList,
		&fiber.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid or Expired Authentication Token",
		},
	)
	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"errors": errorList})
}

// RequireEmployee RequireAdmin Ensures A route Can Only Be Accessed by an Admin user
// This function can be extended to handle different roles
func RequireEmployee(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	role := claims["userType"].(string)

	var errorList []*fiber.Error

	if role != constants.USER_TYPE_EMPLOYEE {
		errorList = append(
			errorList,
			&fiber.Error{
				Code:    fiber.StatusUnauthorized,
				Message: "You're Not Authorized",
			},
		)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"errors": errorList})
	}
	return c.Next()
}
