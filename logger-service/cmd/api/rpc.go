package main

import (
	"context"
	"log"
	"logger-service/data"
	"time"
)

// type to RPC server
type RPCServer struct {
}

type RPCPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Error writing to mongo: ", err)
		return err
	}

	*resp = "Process Payload via RPC: " + payload.Name
	return nil
}
