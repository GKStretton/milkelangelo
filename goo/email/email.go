// Package email support sending emails, via mailjet. Supports sending via
// mqtt.
// Uses the machinepb.Email message. You can set the recipients in kv/x where
// x is one of EMAIL_RECIPIENT_MAINTENANCE, EMAIL_RECIPIENT_ROUTINE_MAINTENANCE,
// EMAIL_RECIPIENT_SOCIAL_NOTIFICATIONS; email addrs delimited by new line.
package email

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/mailjet/mailjet-apiv3-go/v4"
)

var l = log.New(os.Stdout, "[EMAIL] ", log.Flags())
var mj *mailjet.Client

func Start() {
	mqtt.Subscribe(topics_backend.TOPIC_EMAIL_SEND, handleMQTTEmailSend)

	API_KEY := keyvalue.Get("MAILJET_API_KEY")
	if len(API_KEY) < 1 {
		l.Println("failed to make email client, api key missing")
		return
	}
	API_SECRET := keyvalue.Get("MAILJET_API_SECRET")
	if len(API_SECRET) < 1 {
		l.Println("failed to make email client, api secret missing")
		return
	}

	mj = mailjet.NewMailjetClient(string(API_KEY), string(API_SECRET))
}

func GetAddressesOfRecipient(recipient machinepb.EmailRecipient) ([]string, error) {
	addr := string(keyvalue.Get(recipient.String()))
	if len(addr) < 3 {
		return nil, fmt.Errorf("recipient for %s not set in kv/", recipient.String())
	}
	addrs := strings.Split(addr, "\n")
	return addrs, nil
}

func SendEmail(email *machinepb.Email) error {
	if mj == nil {
		return errors.New("email client not initialised")
	}
	addrs, err := GetAddressesOfRecipient(email.Recipient)
	if err != nil {
		return err
	}

	recipients := []mailjet.RecipientV31{}
	for _, a := range addrs {
		recipients = append(recipients, mailjet.RecipientV31{
			Email: a,
		})
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "ops@study-of-light.com",
				Name:  "Milkelangelo Operations",
			},
			To:       (*mailjet.RecipientsV31)(&recipients),
			Subject:  email.Subject,
			TextPart: email.Body + footer(email),
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mj.SendMailV31(&messages)
	if err != nil {
		return err
	}
	l.Printf("Sent email '%s' to %s. Status: %s\n", email.Subject, recipients, res.ResultsV31[0].Status)
	return nil
}

func footer(email *machinepb.Email) string {
	if email.Recipient != machinepb.EmailRecipient_EMAIL_RECIPIENT_ROUTINE_OPERATIONS {
		return ""
	}

	b, err := os.ReadFile("../resources/social_text/splashtext.csv")
	if err != nil {
		l.Printf("error getting splashtext for footer: %s\n", err)
		return ""
	}
	lines := strings.Split(string(b), "\n")
	lines = lines[1:]
	if len(lines) < 1 {
		l.Printf("error getting splashtext for footer: no lines in file\n")
		return ""
	}
	i := rand.Intn(len(lines))
	choice := lines[i]

	return fmt.Sprintf("\n\n%s", choice)
}
