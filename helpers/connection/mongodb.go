package connection

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// NewMongoDB connects to MongoDB database using the official MongoDB Go driver and returns the database instance.
// It takes a timeout duration, a MongoDB URI, and a database name as parameters.
// The MongoDB URI should be in the format: mongodb://username:password@host:port/database?query_params
//
// Example: mongodb://user:password@localhost:27017/go-template?ssl=false
func NewMongoDB(timeout time.Duration, URI, dbName string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	clientOptions := options.Client()
	clientOptions.ApplyURI(URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if dbName == "" {
		connString, _ := connstring.Parse(URI)
		dbName = connString.Database
	}

	return client.Database(dbName)
}
