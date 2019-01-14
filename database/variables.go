package database

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

//var session *mgo.Session
var (
	client   *mongo.Client
	database *mongo.Database
)

// Initialize initializes the global variables
func Initialize(connectionURL string, databaseName string) {

	client, err := mongo.NewClient(connectionURL)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	database = client.Database(databaseName)
}
