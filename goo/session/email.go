package session

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/email"
)

const cleaningBody = `Please perform the cleaning routine:

http://192.168.0.37:5000/cleaning (Must be on home wifi)
`

func requestCleaning(session *Session) {
	err := email.SendEmail(&machinepb.Email{
		Subject:   fmt.Sprintf("%d (%d): Please clean/maintain system", session.ProductionId, session.Id),
		Body:      cleaningBody,
		Recipient: machinepb.EmailRecipient_EMAIL_RECIPIENT_ROUTINE_OPERATIONS,
	})
	if err != nil {
		fmt.Println(err)
	}
}

const choosePieceBody = `Please choose the image that will be used for the latest session.

- Go to http://192.168.0.37:5000/content (Must be on home wifi).
`

func requestPieceSelection(session *Session) {
	err := email.SendEmail(&machinepb.Email{
		Subject:   fmt.Sprintf("%d (%d): Time to choose a session image!", session.ProductionId, session.Id),
		Body:      choosePieceBody,
		Recipient: machinepb.EmailRecipient_EMAIL_RECIPIENT_ROUTINE_OPERATIONS,
	})
	if err != nil {
		fmt.Println(err)
	}
}
