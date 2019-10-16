package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	//CollectionNameCounters is the name of the collection for Counters
	CollectionNameCounters = "-_-Counters"
	//CollectionNameClients is the name of the collection for Clients
	CollectionNameClients = "-_-Clients"
	//CollectionNameCollections is the name of the collection for Collections
	CollectionNameCollections = "-_-Collections"
	//CollectionNameDatatypes is the name of the collection for Datatypes
	CollectionNameDatatypes = "-_-Datatypes"
)

const (
	//ID is an identifier of MongoDB
	ID = "_id"
)

func filterByID(id interface{}) bson.D {
	return bson.D{bson.E{Key: ID, Value: id}}
}

func filterByName(name string) bson.D {
	return bson.D{bson.E{Key: "name", Value: name}}
}

// options
var (
	upsert       = true
	upsertOption = &options.UpdateOptions{
		Upsert: &upsert,
	}
)
