package user_service

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v4"
	"hintword.com/api/app/common/utility"
	"hintword.com/api/app/models"

	cfg "hintword.com/api/app/configs"
)

type UserClaims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type UserDetails struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type GoogleUserInfo struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Mobile   string `json:"mobile"`
	Avatar   string `json:"avatar" validate:"required"`
	GoogleID string `json:"google_id" validate:"required"`
}

func HandleGoogleUserInfo(params GoogleUserInfo) (string, error) {
	var user models.Users
	user, _ = user.FindByEmail(params.Email)
	if user.UserID == "" {
		user = models.Users{
			UserID:   utility.GenerateUUID(),
			Email:    params.Email,
			Name:     params.Name,
			GoogleID: params.GoogleID,
			Status:   models.StatusActive,
		}
	}
	user.AvatarURL = params.Avatar

	_, err := user.Save()
	if err != nil {
		return "", fmt.Errorf("failed to save user: %v", err)
	}

	// Generate JWT token
	jwtToken, err := GenerateToken(user.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return jwtToken, nil
}

func GetUserObject(c interface{}) (details UserDetails) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered from panic %v", r)
		}
		return
	}()

	var user *jwt.Token
	if ctx, ok := c.(*fiber.Ctx); ok {
		user = ctx.Locals("user").(*jwt.Token)
	} else if wCtx, _ok := c.(*websocket.Conn); _ok {
		user = wCtx.Locals("user").(*jwt.Token)
	} else {
		fmt.Println("Content type not found unable to fetch user")
		return
	}

	claims := user.Claims.(jwt.MapClaims)
	jsonData, err := utility.MapToJSON(claims)
	if err != nil {
		fmt.Println("Error converting map to JSON:", err)
		return
	}
	err = utility.JSONToStruct(jsonData, &details)
	if err != nil {
		fmt.Println("Error converting JSON to struct:", err)
		return
	}
	return
}

func GenerateToken(userID string) (string, error) {
	// Get user details from database
	var user models.Users
	user, _ = user.FindById(userID)
	if user.UserID == "" {
		return "", errors.New("user not found")
	}

	expireTime := time.Now().Add(time.Hour * 24)

	claims := UserClaims{
		user.UserID,
		user.Email,
		user.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    cfg.GetConfig().JWTIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.GetConfig().JWTAccessSecret))
}
