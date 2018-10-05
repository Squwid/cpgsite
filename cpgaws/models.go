package cpgaws

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Region is the region that the server is hosted in
var Region = "us-east-1"

// AdminEmails holds the admin emails that have op privilige
var AdminEmails = []string{"bwhitelaw@sbcglobal.net"}

// Inquiry is an inquiry about a real estate property
type Inquiry struct {
	UUID        string `json:"uuid"`
	TimeCreated string `json:"timeCreated"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateInquiry creates a new inquiry to use to send on SES
func CreateInquiry(email, phone, name, desc string) *Inquiry {
	uid, _ := uuid.NewV4()
	return &Inquiry{
		UUID:        uid.String(),
		TimeCreated: time.Now().String(),
		Email:       email,
		Phone:       phone,
		Name:        name,
		Description: desc,
	}
}
