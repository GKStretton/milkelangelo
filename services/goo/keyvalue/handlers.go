package keyvalue

import (
	"fmt"

	"github.com/gkstretton/dark/services/goo/mqtt"
)

func setCallback(topic string, payload []byte) {
	key := getLastTopicValue(topic)

	err := setKeyValue(key, payload)
	if err != nil {
		fmt.Printf("failed to set key %s to %v: %v\n", key, payload, err)
	} else {
		mqtt.Publish(TOPIC_ROOT+"set-ack", []byte(key))
		sendToSubs(key)
	}
}

func reqCallback(topic string, payload []byte) {
	key := string(payload)
	sendToSubs(key)
}

func sendToSubs(key string) {
	value, err := getKeyValue(key)
	if err != nil {
		fmt.Printf("error getting value for key %s: %v\n", key, err)
		value = []byte{}
	}
	err = mqtt.Publish(TOPIC_ROOT+"get-resp/"+key, value)
	if err != nil {
		fmt.Printf("error publishing value for key %s: %v\n", key, err)
	}
}
