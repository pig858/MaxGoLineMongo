package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	Name    string `bson:name, omitempty`
	Content string `bson:content, omitempty`
	Time    int64  `bson:sendTime, omitempty`
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

func GetByName(collection *mongo.Collection, name string) []Message {
	var messages []Message

	filter := bson.D{{"name", name}}
	cursor, err := collection.Find(Ctx, filter)

	if err != nil {
		panic(err)
	}

	if err = cursor.All(Ctx, &messages); err != nil {
		panic(err)
	}
	return messages
}

func Insert(collection *mongo.Collection, message Message) {
	_, err := collection.InsertOne(Ctx, message)
	if err != nil {
		panic(err)
	}
}
