package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client   *mongo.Client
	Database *mongo.Database
}

const databaseName = "sharedacardb"
const mongoURI = "mongodb://localhost/" + databaseName

var Instance MongoInstance

func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	database := client.Database(databaseName)

	if err != nil {
		return err
	}

	Instance = MongoInstance{
		Client:   client,
		Database: database,
	}

	return nil
}
