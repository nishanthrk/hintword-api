package company

type Application struct {
	CompanyId     string `json:"company_id"`
	Name          string `json:"name" validate:"required"`
	Mobile        string `json:"mobile" validate:"required"`
	Email         string `json:"email" validate:"required"`
	ContactPerson string `json:"contact_person" validate:"required"`
	NoOfUsers     int    `json:"no_of_users"`
	Dedupe        bool   `json:"dedupe,omitempty"`
	Otp           string `json:"otp,omitempty"`
	Resend        bool   `json:"resend,omitempty"`
	RetryType     string `json:"retry_type,omitempty"`
	// -1 -> Rejected, 1 -> Created, 2 -> Pending for approval, 3 -> Approval, 4 -> In Principle Approval
}

type AgentObject struct {
	AgentId      string `json:"agent_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	ImageUrl     string `json:"image_url"`
	LastSyncDate string `json:"last_sync_date"`
	Status       int    `json:"status"`
}

type SubscriptionObject struct {
	SubscriptionId string      `json:"subscription_id"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	LogoUrl        string      `json:"logo_url"`
	CaptureUrl     string      `json:"capture_url"`
	CallbackUrl    string      `json:"callback_url"`
	Credentials    interface{} `json:"credentials"`
}
