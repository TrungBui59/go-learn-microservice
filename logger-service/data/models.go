package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var client *mongo.Client

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_t" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

// insert value into database
func (l *LogEntry) Insert(entry LogEntry) error {
	// declare collection(same as table in SQL)
	collections := client.Database("logs").Collection("logs")

	// insert entry to collections
	_, err := collections.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Printf("Error inserting into log: ", err)
		return err
	}
	return nil
}

// get all entry
func (l *LogEntry) ALl() ([]*LogEntry, error) {
	// get context to prevent timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// get the collections to interact
	collection := client.Database("logs").Collection("logs")

	// specify option when manipulate daTa
	opts := options.Find()
	// specify sort the result by created_at
	opts.SetSort(bson.D{{"created_at", -1}})

	// get iterator to get data
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Error Finding all docs: ", err)
		return nil, err
	}

	// get all the rows from iterator
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		item := LogEntry{}
		err := cursor.Decode(item)
		if err != nil {
			log.Println("Error decodes entry: ", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

// get one entry by id
func (l *LogEntry) getOne(id string) (*LogEntry, error) {
	// get context to prevent timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// get the collections to interact
	collection := client.Database("logs").Collection("logs")

	// change id to correct format to use for searching in mongo
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// search for entry
	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// delete all the logs collections
func (l *LogEntry) DropCollection() error {
	// get context to prevent timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// get the collections to interact
	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}

// update entry
func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	// get context to prevent timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// get the collections to interact
	collection := client.Database("logs").Collection("logs")

	// change id to correct format to use for searching in mongo
	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docId},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		})
	if err != nil {
		return nil, err
	}

	return result, nil
}
