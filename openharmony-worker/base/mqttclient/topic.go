package mqttclient

import (
	"fmt"
	"strings"
)

// GetClientInfoFromTopic 从topic中解析出client的username和clientid
func GetClientInfoFromTopic(topic string) (string, string, error) {
	array := strings.Split(topic, "/")
	if len(array) < 3 {
		return "", "", fmt.Errorf("topic {%v} split / array length {%v} < 3", topic, len(array))
	}
	return array[1], array[2], nil
}
