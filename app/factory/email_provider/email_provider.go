package email_provider

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type EmailProvider interface {
	SendEmail(toMail string, toName string, subject string, htmlContent string) error
}

type Credentials struct {
	AuthKey      string `json:"auth_key"`
	SenderName   string `json:"sender_name"`
	NoReplyEmail string `json:"no_reply_email"`
}

func GetEmailProvider(provider string, data datatypes.JSON) EmailProvider {
	credentials := Credentials{}
	err := json.Unmarshal(data, &credentials)
	if err != nil {
		panic("Unable to process credentials")
	}
	switch provider {
	case "SENDGRID":
		sendgrid := &Sendgrid{
			AuthKey:     credentials.AuthKey,
			SenderName:  credentials.SenderName,
			NoReplyMail: credentials.NoReplyEmail,
		}
		return sendgrid.GetInstance()
	default:
		panic("Unsupported cloud storage provider")
	}
}
