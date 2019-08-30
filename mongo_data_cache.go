package main

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDataCache struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func newMongoDataCache() *mongoDataCache {
	mdc := new(mongoDataCache)
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		Fatal(err)
	}
	mdc.client = client
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		Fatal(err)
	}
	mdc.collection = mdc.client.Database("etsello").Collection("users")
	return mdc
}

func (mdc mongoDataCache) saveDetailsToCache(userID int, info userInfo) {
	filter := bson.D{{"_id", userID}}
	bsonDocument, _ := toDoc(info)
	upsertValue := true
	uOptions := options.MergeReplaceOptions()
	uOptions.Upsert = &upsertValue
	_, err := mdc.collection.ReplaceOne(context.TODO(), filter, bsonDocument, uOptions)
	if err != nil {
		Fatal(err)
	}
}

func (mdc mongoDataCache) getUserInfo(userID int) (*userInfo, error) {
	var result userInfo
	filter := bson.D{{"_id", userID}}
	err := mdc.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		Error(err)
		return nil, err
	}
	return &result, nil
}

func (mdc mongoDataCache) getUserMap() map[int]userInfo {
	results := make(map[int]userInfo)

	cur, err := mdc.collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem userInfo
		err := cur.Decode(&elem)
		if err != nil {
			Fatal(err)
		}
		results[elem.UserID] = elem
	}
	return results
}

func (mdc mongoDataCache) disconnectCache() {
	err := mdc.client.Disconnect(context.TODO())
	if err != nil {
		Fatal(err)
	}
	Info("Connection to mongodb closed.")
}

func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
