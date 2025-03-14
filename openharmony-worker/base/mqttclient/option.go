package mqttclient

import "errors"

// GetClientID 获取mqtt客户端clientid
func (md *Driver) GetClientID() (string, error) {
	if md.client == nil {
		return "", errors.New("mqtt client nil error")
	}

	r := (*md.client).OptionsReader()
	return r.ClientID(), nil
}

// GetUserName 获取mqtt客户端clientid
func (md *Driver) GetUserName() (string, error) {
	if md.client == nil {
		return "", errors.New("mqtt client nil error")
	}

	r := (*md.client).OptionsReader()
	return r.Username(), nil
}
