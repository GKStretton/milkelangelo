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
			"Please provide %d whole milk by %s today",
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
