package keyvalue

import (
	"fmt"
	"strings"
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

func GetBool(key string) bool {
	b, err := getKeyValue(key)
	if err != nil {
		return false
	}
	if strings.TrimSpace(strings.ToLower(string(b))) == "true" {
		return true
	}
	return false
}

func SetBool(key string, value bool) {
	val := "false"
	if value {
		val = "true"
	}
	err := setKeyValue(key, []byte(val))
	if err != nil {
		fmt.Println(err)
	}
}

func Get(key string) []byte {
	b, err := getKeyValue(key)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

func GetString(key string) string {
	return string(Get(key))
}

func Set(key string, value []byte) {
	err := setKeyValue(key, value)
	if err != nil {
		fmt.Println(err)
	}
}
