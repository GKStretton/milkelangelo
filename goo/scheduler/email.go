package scheduler

import (
	"fmt"
	"os"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/email"
)

func requestFridgeMilk() {
	err := email.SendEmail(&machinepb.Email{
		Subject: fmt.Sprintf(
			"Please provide %d whole milk by %s today",
			bulkVolume,
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

const cleaningBody = `Please perform the cleaning routine:

- Fill altar with cleaning solution
- Run the cleaning cycle @ http://192.168.0.37:5000 "Cleaning" tab. (Must be on home wifi)
- Change the rinse glass water
- Ensure vials are topped up with correct fluids
- Wipe the bowl clean if necessary, once the cycle is done
- Issue a shutdown
- Empty bucket and clean it if necessary
- Replace bucket, ! with tube in it !
`

func requestCleaning() {
	err := email.SendEmail(&machinepb.Email{
		Subject:   "Please clean/maintain system",
		Body:      cleaningBody,
		Recipient: machinepb.EmailRecipient_EMAIL_RECIPIENT_ROUTINE_OPERATIONS,
	})
	if err != nil {
		fmt.Println(err)
	}
}

const choosePieceBody = `Please choose the image that will be used for the latest session.

- Go to http://192.168.0.37:5000 "Socials" tab. (Must be on home wifi).
`

func requestPieceSelection() {
	err := email.SendEmail(&machinepb.Email{
		Subject:   "Time to choose a session image!",
		Body:      choosePieceBody,
		Recipient: machinepb.EmailRecipient_EMAIL_RECIPIENT_ROUTINE_OPERATIONS,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	email.Start()
	requestCleaning()
	os.Exit(0)
}
