package cpgaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// Send sends an inquiry to the customer
func (inq *Inquiry) Send() error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(Region)},
	}))

	svc := ses.New(sess)
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: aws.StringSlice(AdminEmails),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("utf-8"),
					Data:    aws.String(inq.Description),
				},
				Text: &ses.Content{
					Charset: aws.String("utf-8"),
					Data:    aws.String(inq.Description),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    aws.String(fmt.Sprintf("New Inquiry from %s", inq.Name)),
			},
		},
		Source: aws.String("bwhitelaw@sbcglobal.net"),
	}
	_, err := svc.SendEmail(input)
	if err != nil {
		return err
	}
	return nil
}
