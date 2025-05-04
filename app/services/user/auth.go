package user

import (
	"fmt"
	"hintword.com/api/app/configs"
	"hintword.com/api/app/models"

	hOuth "hintword.com/api/app/services/oauth"
)

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
