package sms_provider

import (
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/datatypes"
	"hintword.com/api/app/common/constants"
)

type SmsProvider interface {
	SendOtp(templateId string, mobile string, payload interface{}, spoof bool) error

	ResendOtp(retryType string, mobile string, payload interface{}, spoof bool) error

	VerifyOtp(otp string, mobile string, spoof bool) error
}

type Credentials struct {
	AuthKey string `json:"auth_key"`
}

func GetSmsProvider(provider string, data datatypes.JSON) (p SmsProvider) {
	credentials := Credentials{}
	fmt.Println("data :: ", data)
	err := json.Unmarshal(data, &credentials)
	if err != nil {
		log.Println("Unable to process credentials")
		return
	}
	switch provider {
	case "MSG91":
		msg91 := &Msg91{
			AuthKey: credentials.AuthKey,
			BaseUrl: constants.Msg91BaseUrl,
		}
		return msg91.GetInstance()
	default:
		panic("Unsupported cloud storage provider")
	}
}
