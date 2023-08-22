package thirdparty

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


func NewMongoClient(mongoConnectionString string) *mongo.Client {
	// create a context with timeout of 10s
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect with mongo
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnectionString))
	if err != nil {
		log.Fatal("Unable to connect mongodb")
	}
	// ping to establish the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Unable to connect mongodb")
	}
	log.Print("Connected to DB")

	return client
}