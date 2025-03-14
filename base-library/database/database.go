package database

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Driver mongo客户端
type Driver struct {
	info   *DInfo
	client *mongo.Client
}

// Driver 全局实例
var (
	Mdriver *Driver
)

func (dd *Driver) ping() error {
	if dd.client == nil {
		err := errors.New("dd.client is nil error")
		log.Error(err)
		return err
	}
	if err := dd.client.Ping(context.TODO(), nil); err != nil {
		log.Errorf("Ping error {%v}", err)
		return err
	}
	return nil
}

func (dd *Driver) init() error {
	mongoURI, err := dd.info.Client.getURI()
	if err != nil {
		log.Errorf("getURI() error {%s}", err)
		return err
	}

	clientOpts := options.Client().ApplyURI(mongoURI)
	dd.client, err = mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Errorf("mongo.Connect error {%s}", err)
		return err
	}

	err = dd.ping()
	if err != nil {
		log.Errorf("dd.ping() error {%s}", err)
		return err
	}

	return nil
}

// New 创建实例
func New(info *DInfo) (*Driver, error) {
	dd := &Driver{info: info, client: nil}
	if err := dd.init(); err != nil {
		log.Errorf("dd.init() error {%v}\n", err)
		return nil, err
	}
	return dd, nil
}
