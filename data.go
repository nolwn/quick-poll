package main

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "quickpoll"
const dbTimeout = 5 * time.Second

// Tables
const (
	TablePoll = "polls"
	TableVote = "votes"
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

func query(
	collection string, parameters interface{},
) (results []bson.M, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return
	}

	client.Database(dbName).Collection(collection)

	cursor, err := client.Database(dbName).Collection(collection).Find(ctx, parameters)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = cursor.All(ctx, &results)

	if results == nil {
		results = []bson.M{}
	}

	return renameIds(results), nil
}

func queryById(collection string, id string) (entity primitive.M, err error) {
	var oid primitive.ObjectID

	oid, err = primitive.ObjectIDFromHex(id)
	parameters := bson.D{{Key: "_id", Value: oid}}

	if err != nil {
		return nil, err
	}

	entities, err := query(collection, parameters)

	if len(entities) > 0 {
		entity = entities[0]
		return
	}

	return
}

func add(collection string, item interface{}) (id string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return
	}

	var result *mongo.InsertOneResult
	result, err = client.Database(dbName).Collection(collection).InsertOne(ctx, item)
	if err != nil {
		return
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		id = oid.Hex()
	} else {
		err = errors.New("MongoDB did not return a valid id for this document")
	}

	return
}

func update(collection string, id string, item idable) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return
	}

	var oid primitive.ObjectID
	oid, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	itemId := item.id()
	item.setId("")

	_, err = client.Database(dbName).Collection(collection).ReplaceOne(ctx, bson.D{{Key: "_id", Value: oid}}, item)

	item.setId(itemId)

	return
}

func renameIds(entities []primitive.M) []primitive.M {
	for _, entity := range entities {
		entity["id"] = entity["_id"]
		delete(entity, "_id")
	}

	return entities
}
