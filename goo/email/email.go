package email

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Start() {
	mqtt.Subscribe(topics_backend.TOPIC_EMAIL_SEND, func(topic string, payload []byte) {
		email := &machinepb.Email{}
		err := handleMessage(payload, email)
		if err != nil {
			fmt.Printf("failed to parse email: %v\n", err)
			return
		}
		err = SendEmail(email)
		if err != nil {
			fmt.Println(err)
			return
		}
	})
}

func SendEmail(email *machinepb.Email) error {
	recipient := os.Getenv("EMAIL_RECIPIENT")
	if recipient == "" {
		return fmt.Errorf("cannot send email %+v, recipient email not set", email)
	}
	cmd := exec.Command("sh", "-c", fmt.Sprintf(`echo "%s" | mail -s "%s" %s`, email.Body, email.Subject, recipient))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to send email %+v: %v", email, err)
	}
	return nil
}

func handleMessage(data []byte, msg protoreflect.ProtoMessage) error {
	// try unmarshaling as protobuf
	if errProto := proto.Unmarshal(data, msg); errProto != nil {
		// if protobuf fails, try unmarshaling as JSON
		if errJson := protojson.Unmarshal(data, msg); errJson != nil {
			// if both fail, return an error
			return fmt.Errorf("unable to unmarshal message. errProto: %v, errJson: %v", errProto, errJson)
		}
	}

	return nil
}
