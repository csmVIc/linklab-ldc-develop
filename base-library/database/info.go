package database

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// AddressInfo mongo地址信息
type AddressInfo struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// ClientInfo 客户端信息
type ClientInfo struct {
	Address    []AddressInfo `json:"address"`
	User       string        `json:"user"`
	Password   string        `json:"password"`
	Db         string        `json:"db"`
	Replicaset string        `json:"replicaset"`
}

// TransactionInfo 事务信息
type TransactionInfo struct {
	TimeOut int `json:"timeout"`
}

// DInfo mongo信息
type DInfo struct {
	Client      ClientInfo      `json:"client"`
	Transaction TransactionInfo `json:"transaction"`
}

func (ci *ClientInfo) checkValue() bool {
	if len(ci.Address) < 1 {
		return false
	}
	return true
}

func (ci *ClientInfo) getURI() (string, error) {

	if ci.checkValue() == false {
		return "", errors.New("ci.checkValue() checkValue error")
	}

	addressstr := ""
	for i, address := range ci.Address {
		if i != 0 {
			addressstr = fmt.Sprintf("%s,%s:%s", addressstr, address.Host, address.Port)
		} else {
			addressstr = fmt.Sprintf("%s:%s", address.Host, address.Port)
		}
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=%s&replicaSet=%s", ci.User, ci.Password, addressstr, ci.Db, ci.Replicaset)
	log.Infof("mongo uri %s", uri)
	return uri, nil
}
