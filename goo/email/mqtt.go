package email

import (
	"fmt"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func handleMQTTEmailSend(topic string, payload []byte) {
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
