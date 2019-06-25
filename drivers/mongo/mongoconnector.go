package mongo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload" // .env autoloader
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection details from env vars or a .env file
var (
	User     string = os.Getenv("MONGODB_USERNAME")
	Password string = os.Getenv("MONGODB_PASSWORD")
	Host     string = os.Getenv("MONGODB_HOST")
	Database string = os.Getenv("MONGODB_DATABASE")
)

// Basic structure to hold the connection
type MongoConnector struct {
	mongoClient *mongo.Client
}

// Connect to mongo and set up the structure with an active connection
func (d *MongoConnector) Connect() error {
	var err error

	dnsConnectionString := fmt.Sprintf("mongodb://%s:%s@%s:27017/%s", User, Password, Host, Database)

	d.mongoClient, err = mongo.NewClient(options.Client().ApplyURI(dnsConnectionString))

	if err != nil {
		return err
	}

	// Establish connection to database
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	if err = d.mongoClient.Connect(ctx); err != nil {
		return err
	}

	return nil
}

// Disconnect from mongo (should you ever want to)
func (d *MongoConnector) Disconnect() error {
	err := d.mongoClient.Disconnect(context.TODO())

	if err != nil {
		return err
	}

	return nil
}

// return the collection from the active connection
func (d *MongoConnector) GetCollection(name string) *mongo.Collection {
	return d.mongoClient.Database(Database).Collection(name)
}

// Just a basic panic recover that can be used in a defer somewhere
func PanicRecoverHandler(err error) {
	if r := recover(); r != nil {
		switch x := r.(type) {
		case string:
			err = errors.New(x)
		case error:
			err = x
		default:
			err = errors.New("Unknown panic")
		}
	}
}
