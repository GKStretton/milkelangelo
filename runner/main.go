package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/email"
	"github.com/gkstretton/dark/services/goo/mqtt"
)

func main() {
	fmt.Println("Launching runner")

	mqtt.Start()
	mqtt.Subscribe(topics_backend.TOPIC_GENERATE_CONTENT, func(topic string, payload []byte) {
		n, err := strconv.Atoi(string(payload))
		if err != nil {
			fmt.Printf("failed to parse session number from content generation request: %v\n", err)
			return
		}
		generateContent(n)
	})

	email.Start()

	for {
		time.Sleep(time.Second)
	}
}

func generateContent(sessionNumber int) {
	fmt.Printf("generate %d\n", sessionNumber)

	planCmd := exec.Command("./scripts/generate-content-plan", strconv.Itoa(sessionNumber))
	planCmd.Stdout = os.Stdout
	planCmd.Stderr = os.Stderr
	err := planCmd.Run()
	if err != nil {
		handleError(fmt.Errorf("error generating content plan: %v", err))
		return
	}
	fmt.Printf("generated content plan for %d\n", sessionNumber)

	contentCmd := exec.Command("./scripts/generate-content", strconv.Itoa(sessionNumber))
	contentCmd.Stdout = os.Stdout
	contentCmd.Stderr = os.Stderr
	err = contentCmd.Run()
	if err != nil {
		handleError(fmt.Errorf("error generating content: %v %v", err, planCmd.Err))
		return
	}

	fmt.Printf("generated content for %d\n", sessionNumber)
	handleSuccess(sessionNumber)
}

func handleSuccess(sessionNumber int) {
	e := &machinepb.Email{
		Recipient: machinepb.EmailRecipient_EMAIL_RECIPIENT_MAINTENANCE,
		Subject:   fmt.Sprintf("Generated content for %d", sessionNumber),
		Body:      fmt.Sprintf("successfully generated content for session id %d\n", sessionNumber),
	}
	err := email.SendEmail(e)
	if err != nil {
		fmt.Printf("failed to email confirmation of generation: %v", err)
		return
	}
	fmt.Printf("sent email confirming generation of content for %d\n", sessionNumber)
}

func handleError(err error) {
	fmt.Println(err)
	e := &machinepb.Email{
		Recipient: machinepb.EmailRecipient_EMAIL_RECIPIENT_MAINTENANCE,
		Subject:   "error in asol runner",
		Body:      err.Error(),
	}
	emailErr := email.SendEmail(e)
	if emailErr != nil {
		fmt.Printf("failed to email error: %v (%v)", emailErr, err)
		return
	}
}
