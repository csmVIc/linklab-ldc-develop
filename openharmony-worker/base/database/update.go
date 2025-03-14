package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UpdateElem 更新表中的字段
func (dd *Driver) UpdateElem(cname string, filter interface{}, elem interface{}) (*mongo.UpdateResult, error) {
	if dd.ping() != nil && dd.init() != nil {
		return nil, errors.New("UpdateElem dd.ping()/dd.init() false")
	}

	collection := dd.client.Database(dd.info.Client.Db).Collection(cname)
	result, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": elem})
	return result, err
}
