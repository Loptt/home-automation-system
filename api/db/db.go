package db

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	db          *mongo.Database
	initialized = false
)

// Connect creates a new MongoDB database and connects to it
func Connect(url string) error {
	// Database Config
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.NewClient(clientOptions)
	//Set up a context required by mongo.Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	//Cancel context to avoid memory leak
	defer cancel()

	// Ping our db connection
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}

	db = client.Database("home_automation")
	log.Printf("Connected to DB %s\n", url)
	initialized = true

	return nil
}

// Database gets the database object
func Database() (*mongo.Database, error) {
	if !initialized {
		return nil, errors.New("Calling database without initializing")
	}

	return db, nil
}
