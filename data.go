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

// TODO: this should come from the environment
const dbUri = "mongodb://127.0.0.1:27017" //ultimately this should come from the env

var opts = options.Client().ApplyURI(dbUri)

// query looks up an item in the database. It takes a collection, and a parameters which
// will be a list of properties to look up.
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

// queryById takes a collection and an id and builds a query to look up the item with
// that id.
func queryById(collection string, id string) (entity primitive.M, err error) {
	var oid primitive.ObjectID

	// mongo db ids aren't strings, so the id string needs to be turned into a mongo id
	// object.
	oid, err = primitive.ObjectIDFromHex(id)
	parameters := bson.D{{Key: "_id", Value: oid}}

	if err != nil { // Something must be wrong with the id. So... not found!
		return nil, nil
	}

	entities, err := query(collection, parameters)

	// need to make sure that something was actually found before you try to
	if len(entities) > 0 {
		entity = entities[0]
		return
	}

	return
}

// add an item to a given collection.
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

	// Mongodb returns an empty interface type as an ID. Why? Beats the Hell out of me.
	// It needs to be casted into an ObjectID type.
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		id = oid.Hex() // turns the ObjectID into a hex string.
	} else {
		err = errors.New("MongoDB did not return a valid id for this document")
	}

	return
}

// update takes a collection, and ID and an item and updates whatever items has that ID
// with the fields of the item passed.
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

	// the item id needs to be zeroed out so it is ignored by MongoDB
	itemId := item.id()
	item.setId("")

	_, err = client.Database(dbName).Collection(collection).ReplaceOne(
		ctx,
		bson.D{{Key: "_id", Value: oid}},
		item,
	)

	item.setId(itemId) // reset it incase the item is used after this function is called

	return
}

// renameIds takes a list of BSON entities that have been returned by MongoDB and renames
// the ID field from _id to id.
func renameIds(entities []primitive.M) []primitive.M {
	for _, entity := range entities {
		entity["id"] = entity["_id"]
		delete(entity, "_id")
	}

	return entities
}
