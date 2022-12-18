package mqtt

import (
	"fmt"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

const QOS = 1

var client paho.Client

type Callback func(topic string, payload []byte)
type subscriptions map[string][]Callback

var subs = subscriptions{}

func Start() {
	_ = paho.CRITICAL
	client = paho.NewClient(paho.NewClientOptions().
		AddBroker("DEPTH:1883").
		SetOnConnectHandler(on_connect).
		SetAutoReconnect(true),
	)
	token := client.Connect()
	if !token.WaitTimeout(time.Second * time.Duration(5)) {
		panic("mqtt client connect timed out")
	}
	if token.Error() != nil {
		panic(token.Error())
	}
}

func on_connect(client paho.Client) {
	for topic, cbs := range subs {
		for _, cb := range cbs {
			client.Subscribe(topic, QOS, func(c paho.Client, m paho.Message) {
				cb(m.Topic(), m.Payload())
			})
		}
	}
	fmt.Println("Connected to broker")
}

func Publish(topic string, payload interface{}) error {
	token := client.Publish(topic, QOS, false, payload)
	if !token.WaitTimeout(time.Second) {
		return fmt.Errorf("publish to %s timed out", topic)
	}
	if token.Error() != nil {
		return fmt.Errorf("error publishing to %s: %v", topic, token.Error())
	}
	return nil
}

func Subscribe(topic string, cb Callback) {
	subs[topic] = append(subs[topic], cb)
	if client != nil && client.IsConnectionOpen() {
		client.Subscribe(topic, QOS, func(c paho.Client, m paho.Message) {
			cb(m.Topic(), m.Payload())
		})
	}
}

func Unsubscribe(topic string) {
	delete(subs, topic)
}
