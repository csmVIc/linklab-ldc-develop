package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertElem 数据插入
func (dd *Driver) InsertElem(cname string, elem interface{}) (*mongo.InsertOneResult, error) {
	if dd.ping() != nil && dd.init() != nil {
		return nil, errors.New("InsertElem dd.ping()/dd.init() false")
	}
	collection := dd.client.Database(dd.info.Client.Db).Collection(cname)
	result, err := collection.InsertOne(context.TODO(), elem)
	return result, err
}

// InsertElemIfNotExist 如果不存在则插入
func (dd *Driver) InsertElemIfNotExist(cname string, filter interface{}, elem interface{}) (*mongo.UpdateResult, error) {
	if dd.ping() != nil && dd.init() != nil {
		return nil, errors.New("InsertElemIfNotExist dd.ping()/dd.init() false")
	}

	collection := dd.client.Database(dd.info.Client.Db).Collection(cname)
	result, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$setOnInsert": elem}, options.Update().SetUpsert(true))
	return result, err
}

// ReplaceElem 不管是否存在,都会替换数据
func (dd *Driver) ReplaceElem(cname string, filter interface{}, elem interface{}) (*mongo.UpdateResult, error) {
	if dd.ping() != nil && dd.init() != nil {
		return nil, errors.New("ReplaceElem dd.ping()/dd.init() false")
	}

	collection := dd.client.Database(dd.info.Client.Db).Collection(cname)
	result, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": elem}, options.Update().SetUpsert(true))
	return result, err
}
