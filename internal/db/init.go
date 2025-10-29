package db

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	dbName string
)

func Init() error {
	var err error

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		fmt.Println("MONGODB_URI environment variable not found!")
		return errors.New("missing MONGODB_URI environment variable")
	}

	dbName = os.Getenv("MONGODB_DBNAME")
	if dbName == "" {
		fmt.Println("MONGODB_DBNAME environment variable not found!")
		return errors.New("missing MONGODB_DBNAME environment variable")
	}

	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	fmt.Println("Connected to MongoDB database:", dbName)

	return nil
}

func GetClient() *mongo.Client {
	return Client
}

func GetDBName() string {
	return dbName
}
