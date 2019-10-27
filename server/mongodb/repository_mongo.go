package mongodb

import (
	"context"
	"github.com/knowhunger/ortoo/commons/log"
	"github.com/knowhunger/ortoo/server/mongodb/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RepositoryMongo is a tool struct for MongoDB
type RepositoryMongo struct {
	*CollectionCounters
	*CollectionClients
	*CollectionCollections
	*CollectionDatatypes
	*CollectionOperations
	config *Config
	ctx    context.Context
	client *mongo.Client
	db     *mongo.Database
}

// New creates a new RepositoryMongo
func New(ctx context.Context, conf *Config) (*RepositoryMongo, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Host)) // "mongodb://root:ortoo-test@localhost:27017"))
	if err != nil {
		return nil, log.OrtooErrorf(err, "fail to connect MongoDB")
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, log.OrtooErrorf(err, "fail to ping MongoDB")
	}
	db := client.Database(conf.OrtooDB)
	repo := &RepositoryMongo{
		db:     db,
		ctx:    ctx,
		client: client,
	}
	if err := repo.InitializeCollections(ctx); err != nil {
		return nil, log.OrtooError(err)
	}
	return repo, nil
}

// InitializeCollections initialize collections
func (r *RepositoryMongo) InitializeCollections(ctx context.Context) error {

	r.CollectionCounters = NewCollectionCounters(r.client, r.db.Collection(schema.CollectionNameCounters))
	r.CollectionClients = NewCollectionClients(r.client, r.db.Collection(schema.CollectionNameClients))
	r.CollectionCollections = NewCollectionCollections(r.client, r.CollectionCounters, r.db.Collection(schema.CollectionNameCollections))
	r.CollectionDatatypes = NewCollectionDatatypes(r.client, r.db.Collection(schema.CollectionNameDatatypes))
	r.CollectionOperations = NewCollectionOperations(r.client, r.db.Collection(schema.CollectionNameOperations))
	names, err := r.db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return log.OrtooErrorf(err, "fail to list collection names")
	}
	var realCollections = make(map[string]bool)
	for _, v := range names {
		realCollections[v] = true
	}

	if _, ok := realCollections[schema.CollectionNameClients]; !ok {
		if err := r.CollectionClients.create(ctx, &schema.ClientDoc{}); err != nil {
			return log.OrtooErrorf(err, "fail to create the client collection")
		}
	}
	if _, ok := realCollections[schema.CollectionNameCollections]; !ok {
		if err := r.CollectionCollections.create(ctx, &schema.CollectionDoc{}); err != nil {
			return log.OrtooErrorf(err, "fail to create the collections collection")
		}
	}
	if _, ok := realCollections[schema.CollectionNameDatatypes]; !ok {
		if err := r.CollectionDatatypes.create(ctx, &schema.DatatypeDoc{}); err != nil {
			return log.OrtooErrorf(err, "fail to create the collections collection")
		}
	}
	if _, ok := realCollections[schema.CollectionNameOperations]; !ok {
		if err := r.CollectionOperations.create(ctx, &schema.OperationDoc{}); err != nil {
			return log.OrtooErrorf(err, "fail to create the operations collection")
		}
	}
	return nil
}

// GetOrCreateCollectionSnapshot is a method that gets or creates a collection of snapshot
func (r *RepositoryMongo) GetOrCreateCollectionSnapshot(ctx context.Context, name string) (*CollectionSnapshots, error) {

	names, err := r.db.ListCollectionNames(ctx, schema.FilterByName(name))
	if err != nil {
		return nil, log.OrtooErrorf(err, "fail to list collections")
	}
	collection := newCollectionSnapshot(r.client, r.db.Collection(name), name)
	if len(names) == 0 {
		if err := collection.create(ctx, nil); err != nil {
			return nil, log.OrtooErrorf(err, "fail to create collection")
		}
	}
	return collection, nil
}
