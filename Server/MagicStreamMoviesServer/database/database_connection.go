package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("can't load .env file")
	}

	MongoDB_URI := os.Getenv("MONGODB_URI")
	if MongoDB_URI == "" {
		log.Fatal("MONGODB_URI not set!")
	}

	fmt.Println("MongoDB URI: ", MongoDB_URI)

	clientOptions := options.Client().ApplyURI(MongoDB_URI)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil
	}
	return client
}

var Client *mongo.Client = Connect()	
// i changed this to internal_fn_var inside (OpenCollection), but cancelled it later; cz it will executed with every connection operation!

func OpenCollection(collectionName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to open .env file")
	}

	databaseName := os.Getenv("DATABASE_NAME")
	fmt.Println("DATABASE_NAME: ", databaseName)

	collection := Client.Database(databaseName).Collection(collectionName)
	if collection == nil {
		return nil 
	}
	return collection
}