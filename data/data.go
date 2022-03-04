package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/nolwn/quick-poll/resources"
)

const dbName = "quickpoll"
const dbTimeout = 5 * time.Second

// Tables
const (
	TablePoll = "polls"
)

// const dbUriEnvVar = "MONGO_DB_URI"
const dbUri = "mongodb://127.0.0.1:27017" //ultimately this should come from the env

var opts = options.Client().ApplyURI(dbUri)

type data struct {
	tables map[string][]interface{}
}

var d = data{
	tables: make(map[string][]interface{}),
}

func Query(
	collection string, parameters resources.AddPoll,
) (results []bson.M, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return
	}

	client.Database(dbName).Collection(collection)

	cursor, err := client.Database(dbName).Collection(collection).Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return
	}

	err = cursor.All(ctx, &results)

	return
}

func Add(collection string, item interface{}) (id string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	fmt.Printf("%v\n", item)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return
	}

	var result *mongo.InsertOneResult
	result, err = client.Database(dbName).Collection(collection).InsertOne(ctx, item)
	if err != nil {
		return
	}

	id = fmt.Sprintf("%v", result.InsertedID)

	return
}
