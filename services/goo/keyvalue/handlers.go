package keyvalue

import (
	"fmt"
	"time"

	"github.com/gkstretton/dark/services/goo/mqtt"
)

var respDelay = time.Millisecond * time.Duration(200)

func setCallback(topic string, payload []byte) {
	key := getLastTopicValue(topic)

	err := setKeyValue(key, payload)
	if err != nil {
		fmt.Printf("failed to set key %s to %v: %v\n", key, payload, err)
	} else {
		go func() {
			time.Sleep(respDelay)
			mqtt.Publish(TOPIC_SET_RESP+key, []byte("ack"))
			sendToSubs(key)
		}()
	}
}

func reqCallback(topic string, payload []byte) {
	key := getLastTopicValue(topic)
	go func() {
		time.Sleep(respDelay)
		sendToSubs(key)
	}()
}

func sendToSubs(key string) {
	value, err := getKeyValue(key)
	if err != nil {
		fmt.Printf("error getting value for key %s: %v\n", key, err)
		value = []byte{}
	}
	err = mqtt.Publish(TOPIC_GET_RESP+key, value)
	if err != nil {
		fmt.Printf("error publishing value for key %s: %v\n", key, err)
	}
}
