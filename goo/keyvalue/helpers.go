package keyvalue

import "strings"

func getLastTopicValue(topic string) string {
	sections := strings.Split(topic, "/")
	if len(sections) >= 1 {
		return sections[len(sections)-1]
	}
	return ""
}
