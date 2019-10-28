package mongodb

import (
	"context"
	"github.com/knowhunger/ortoo/commons/log"
	"github.com/knowhunger/ortoo/server/constants"
	"github.com/knowhunger/ortoo/server/mongodb/schema"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionOperations struct {
	*baseCollection
}

func NewCollectionOperations(client *mongo.Client, operations *mongo.Collection) *CollectionOperations {
	return &CollectionOperations{newCollection(client, operations)}
}

func (c *CollectionOperations) InsertOperations(ctx context.Context, operations []interface{}) error {
	result, err := c.collection.InsertMany(ctx, operations)
	if err != nil {
		return log.OrtooError(err)
	}
	if len(result.InsertedIDs) != len(operations) {
		return log.OrtooErrorf(err, "fail to insert operation all")
	}
	return nil
}

func (c *CollectionOperations) DeleteOperation(ctx context.Context, duid string, sseq uint32) (int64, error) {
	f := schema.GetFilter().
		AddFilterEQ(schema.OperationDocFields.DUID, duid).
		AddFilterEQ(schema.OperationDocFields.Sseq, sseq)
	result, err := c.collection.DeleteOne(ctx, f)
	if err != nil {
		return 0, log.OrtooError(err)
	}
	return result.DeletedCount, nil
}

func (c *CollectionOperations) GetOperations(
	ctx context.Context,
	duid string,
	from, to uint64,
	operationDocHandler func(opDoc *schema.OperationDoc) error) error {
	f := schema.GetFilter().
		AddFilterEQ(schema.OperationDocFields.DUID, duid).
		AddFilterGTE(schema.OperationDocFields.Sseq, from)
	if to != constants.InfinitySseq {
		f.AddFilterLTE(schema.OperationDocFields.Sseq, to)
	}
	cursor, err := c.collection.Find(ctx, f)
	if err != nil {
		return log.OrtooError(err)
	}

	for cursor.Next(ctx) {
		var operationDoc schema.OperationDoc
		if err := cursor.Decode(&operationDoc); err != nil {
			return log.OrtooError(err)
		}
		if operationDocHandler != nil {
			if err := operationDocHandler(&operationDoc); err != nil {
				return log.OrtooError(err)
			}
		}
	}
	return nil
}
