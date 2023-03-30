package keyvalue

import (
	"fmt"
	"sync"

	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/mqtt"
)

const TOPIC_SET = topics_backend.TOPIC_KV_SET
const TOPIC_SET_RESP = topics_backend.TOPIC_KV_SET_RESP
const TOPIC_GET = topics_backend.TOPIC_KV_GET
const TOPIC_GET_RESP = topics_backend.TOPIC_KV_GET_RESP

var (
	lock = &sync.Mutex{}
)

func Start() {
	// set
	mqtt.Subscribe(TOPIC_SET+"#", setCallback)
	// req
	mqtt.Subscribe(TOPIC_GET+"#", reqCallback)
}

func Get(key string) []byte {
	b, err := getKeyValue(key)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

func Set(key string, value []byte) {
	err := setKeyValue(key, value)
	if err != nil {
		fmt.Println(err)
	}
}
