package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	UserID  string `bson:userId, omitempty`
	Name    string `bson:name, omitempty`
	Content string `bson:content, omitempty`
	Time    string `bson:sendTime, omitempty`
}
type Follower struct {
	Name   string `bson:name, omitempty`
	UserID string `bson:userId, omitempty`
	Time   string `bson:sendTime, omitempty`
}

var (
	Ctx = context.TODO()
)

func Connect(host string, port string) *mongo.Client {

	dbConnectURI := "mongodb://" + host + ":" + port
	clientOptions := options.Client().ApplyURI(dbConnectURI)
	client, err := mongo.Connect(Ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	return client
}

func GetByFilter[T interface{}](collection *mongo.Collection, filter interface{}) []T {
	var res []T

	cursor, err := collection.Find(Ctx, filter)

	if err != nil {
		panic(err)
	}

	if err = cursor.All(Ctx, &res); err != nil {
		panic(err)
	}
	return res
}

func Insert(collection *mongo.Collection, data any) {
	_, err := collection.InsertOne(Ctx, data)
	if err != nil {
		panic(err)
	}
}

func Delete(collection *mongo.Collection, filter interface{}) {

	_, err := collection.DeleteOne(Ctx, filter)
	if err != nil {
		panic(err)
	}

}
