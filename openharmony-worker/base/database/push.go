package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// PushElemToArray 向数组后添加
func (dd *Driver) PushElemToArray(cname string, filter interface{}, elem interface{}) (*mongo.UpdateResult, error) {
	if dd.ping() != nil && dd.init() != nil {
		return nil, errors.New("PushElemToArray dd.ping()/dd.init() false")
	}

	collection := dd.client.Database(dd.info.Client.Db).Collection(cname)
	result, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$push": elem})
	return result, err
}
