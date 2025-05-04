package email_provider

import (
	"errors"
	"sync"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Sendgrid struct {
	singletonInstance *Sendgrid

	AuthKey string

	SenderName string

	NoReplyMail string
	// Mutex for thread safety
	mu sync.Mutex
}

func (m *Sendgrid) GetInstance() *Sendgrid {

	// Check if the instance already exists
	if m.singletonInstance == nil {
		// Use a mutex to ensure thread safety during instance creation
		m.mu.Lock()
		defer m.mu.Unlock()

		// Check again inside the critical section to prevent race conditions
		if m.singletonInstance == nil {
			// Create the singleton instance
			m.singletonInstance = m
		}
	}

	return m.singletonInstance
}

func (m *Sendgrid) SendEmail(toMail string, toName string, subject string, htmlContent string) (err error) {
	from := mail.NewEmail(m.SenderName, m.NoReplyMail)
	to := mail.NewEmail(toName, toMail)
	plainTextContent := ""
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(m.AuthKey)
	response, err := client.Send(message)
	if err != nil {
		return
	} else {
		if response.StatusCode != 202 {
			return errors.New(response.Body)
		}
	}

	return
}
