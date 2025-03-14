package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

// FindElem 查找数据
func (dd *Driver) FindElem(cname string, filter interface{}) (*mongo.Cursor, error) {
	if dd.ping() != nil && dd.init() != nil {
		return nil, errors.New("FindElem dd.ping()/dd.init() false")
	}
	collection := dd.client.Database(dd.info.Client.Db).Collection(cname)
	cusor, err := collection.Find(context.TODO(), filter)
	return cusor, err
}

// FindOneElem 查找一个数据
func (dd *Driver) FindOneElem(cname string, filter interface{}, result interface{}) error {
	if dd.ping() != nil && dd.init() != nil {
		return errors.New("FindOneElem dd.ping()/dd.init() false")
	}
	collection := dd.client.Database(dd.info.Client.Db).Collection(cname)
	err := collection.FindOne(context.TODO(), filter).Decode(result)
	return err
}

// DocExist 判断文档存在
func (dd *Driver) DocExist(cname string, filter interface{}) error {
	if dd.ping() != nil && dd.init() != nil {
		return errors.New("FindElem dd.ping()/dd.init() false")
	}
	collection := dd.client.Database(dd.info.Client.Db).Collection(cname)
	err := collection.FindOne(context.TODO(), filter).Err()
	return err
}
