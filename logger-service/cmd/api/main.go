package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	webPort  = "80"
	rpcPort  = "50001" // rpc port
	mongoURL = "mongodb://mongo:27107"
	gRpc     = "50001" // for gRPC
)

var client *mongo.Client

type Config struct {
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connection option
	clientOptions := options.Client().ApplyURI(mongoURL)

	// set the authentication(should specify in dockerfile)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect to mongo
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error Connecting: ", err)
		return nil, err
	}

	return c, nil
}
