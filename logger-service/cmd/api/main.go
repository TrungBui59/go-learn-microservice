package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"logger-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

const (
	webPort  = "80"
	rpcPort  = "5001" // rpc port
	mongoURL = "mongodb://mongo:27017"
	gRpc     = "50001" // for gRPC
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient
	// create a context to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// setup app for webserver
	app := Config{
		Models: data.New(client),
	}

	// register and start RPC server
	err = rpc.Register(new(RPCServer))
	go app.rpcListen()
	//// start server
	app.serve()

}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// starting RPC
func (app *Config) rpcListen() error {
	log.Printf("Listening on port %s", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", rpcPort))
	if err != nil {
		log.Printf("Error listening RPC on port %d: %v", rpcPort, err)
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
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

	log.Println("Connected to MongoDB")

	return c, nil
}
