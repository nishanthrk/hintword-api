package user

import (
	"fmt"

	"hintword.com/api/app/common/utility"
	"hintword.com/api/app/configs"
	"hintword.com/api/app/models"

	hOuth "hintword.com/api/app/services/oauth"
)

type GoogleUserInfo struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Mobile   string `json:"mobile"`
	Avatar   string `json:"avatar" validate:"required"`
	GoogleID string `json:"google_id" validate:"required"`
}

func HandleGoogleAuth(code string) (string, error) {
	// Exchange code for token
	token, err := hOuth.ExchangeGoogleCode(code, configs.GetConfig().GoogleOauthClientId, configs.GetConfig().GoogleOauthClientSecret, configs.GetConfig().GoogleOauthRedirectionUrl)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %v", err)
	}

	// Get user info from Google
	userInfo, err := hOuth.GetGoogleUserInfo(token.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %v", err)
	}

	// Find or create user
	var user models.Users
	user, _ = user.FindByEmail(userInfo.Email)
	if user.ID != "" {
		user = models.Users{

			Email:     userInfo.Email,
			Name:      userInfo.Name,
			GoogleID:  userInfo.ID,
			AvatarURL: userInfo.Picture,
		}

		_, err = user.Save()
	}

	// Generate JWT token
	jwtToken, err := hOuth.GenerateJWTToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return jwtToken, nil
}

func HandleGoogleUserInfo(params GoogleUserInfo) (string, error) {
	var user models.Users
	user, _ = user.FindByEmail(params.Email)
	if user.ID == "" {
		user = models.Users{
			ID:       utility.GenerateUUID(),
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
	jwtToken, err := hOuth.GenerateJWTToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return jwtToken, nil
}
