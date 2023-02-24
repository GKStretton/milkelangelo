package mqtt

import (
	"flag"
	"fmt"
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

var (
	host = flag.String("host", "DEPTH", "the hostname of the mqtt broker")
)

const QOS = 1

var client paho.Client

type Callback func(topic string, payload []byte)
type subscriptions map[string][]Callback

var lock sync.Mutex

var subs = subscriptions{}

func Start() {
	_ = paho.CRITICAL
	client = paho.NewClient(paho.NewClientOptions().
		AddBroker(fmt.Sprintf("%s:1883", *host)).
		SetOnConnectHandler(on_connect).
		SetAutoReconnect(true).
		SetConnectionLostHandler(on_disconnect),
	)
	token := client.Connect()
	if !token.WaitTimeout(time.Second * time.Duration(5)) {
		panic("mqtt client connect timed out")
	}
	if token.Error() != nil {
		panic(token.Error())
	}

	client.AddRoute("#", on_message)
}

func on_message(client paho.Client, msg paho.Message) {
	// fmt.Printf("Received mqtt message on '%s': %v\n", msg.Topic(), string(msg.Payload()))
}

func on_disconnect(client paho.Client, err error) {
	fmt.Printf("Disconnected from broker: %v\n", err)
}

func on_connect(client paho.Client) {
	lock.Lock()
	defer lock.Unlock()
	for topic, cbs := range subs {
		for _, cb := range cbs {
			fmt.Printf("subscribing to %s\n", topic)
			t := client.Subscribe(topic, QOS, func(c paho.Client, m paho.Message) {
				cb(m.Topic(), m.Payload())
			})
			if !t.WaitTimeout(time.Second) {
				fmt.Printf("subscribe to %s timed out\n", topic)
			}
			if t.Error() != nil {
				fmt.Printf("error subscribing to %s: %v", topic, t.Error())
			}
		}
	}
	fmt.Println("Connected to broker")
}

func Publish(topic string, payload interface{}) error {
	if client == nil {
		return fmt.Errorf("mqtt client not initialised")
	}
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
	lock.Lock()
	defer lock.Unlock()
	if client == nil {
		fmt.Println("mqtt client not initialised")
	}
	subs[topic] = append(subs[topic], cb)
	if client != nil && client.IsConnectionOpen() {
		client.Subscribe(topic, QOS, func(c paho.Client, m paho.Message) {
			cb(m.Topic(), m.Payload())
		})
	}
}

func Unsubscribe(topic string) {
	lock.Lock()
	defer lock.Unlock()
	delete(subs, topic)
}
