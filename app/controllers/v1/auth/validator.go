package auth_controller

type AuthResponse struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
	ID     string `json:"id"`
	Token  string `json:"token"`
}

var GoogleAuthRequest struct {
	Code string `json:"code" validate:"required"`
}
