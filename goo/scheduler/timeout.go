package scheduler

import (
	"fmt"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/email"
)

func sessionTimeout(d time.Duration, doDrain bool) {
	fmt.Println("error: session timeout not implemented")
	err := email.SendEmail(&machinepb.Email{
		Subject:   fmt.Sprintf("session timeout not implemented (%s, %s)", d, doDrain),
		Body:      "session will not be ended automatically",
		Recipient: machinepb.EmailRecipient_EMAIL_RECIPIENT_MAINTENANCE,
	})
	if err != nil {
		fmt.Println(err)
	}
}
