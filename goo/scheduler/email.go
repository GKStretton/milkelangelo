package scheduler

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/email"
)

func requestSessionIntervention(e error) {
	err := email.SendEmail(&machinepb.Email{
		Subject:   fmt.Sprintf("Intervention required: %s", e),
		Body:      fmt.Sprintf("Please take over the session manually\n\n%s", e.Error()),
		Recipient: machinepb.EmailRecipient_EMAIL_RECIPIENT_MAINTENANCE,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func requestFridgeMilk() {
	err := email.SendEmail(&machinepb.Email{
		Subject: fmt.Sprintf(
			"Please provide %dml whole milk & align tripod by %s today",
			milkVolume,
			mainSessionStartTime.fmtLocal(),
		),
		Body: fmt.Sprintf(
			"As required for the next session at %s today.",
			mainSessionStartTime.fmtLocal(),
		),
		Recipient: machinepb.EmailRecipient_EMAIL_RECIPIENT_ROUTINE_OPERATIONS,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func sendReminder(skip bool) {
	var subject, body string

	if skip {
		subject = "No aSoL session today"
		body = "Session is being skipped"
	} else {
		subject = fmt.Sprintf(
			"Reminder: A Study of Light @ %s today",
			mainSessionStartTime.fmtLocal(),
		)
		body = fmt.Sprintf(
			"A Study of Light session will begin at %s today.",
			mainSessionStartTime.fmtLocal(),
		)

	}

	err := email.SendEmail(&machinepb.Email{
		Subject:   subject,
		Body:      body,
		Recipient: machinepb.EmailRecipient_EMAIL_RECIPIENT_ROUTINE_OPERATIONS,
	})
	if err != nil {
		fmt.Println(err)
	}
}
