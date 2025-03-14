package messenger

import (
	"fmt"
	"strings"
)

//GetIDFromTopic 从topic中解析出id
func GetIDFromTopic(topic string, utype string) (string, error) {
	array := strings.Split(topic, ".")
	if len(array) < 3 {
		return "", fmt.Errorf("topic {%v} split . len(array) {%v} < 3 error", topic, len(array))
	}
	if array[0] != utype {
		return "", fmt.Errorf("topic {%v} utype {%v} != input utype {%v}", topic, array[0], utype)
	}
	return array[len(array)-1], nil
}
