package snapshot

import (
	"github.com/knowhunger/ortoo/pkg/context"
	"github.com/knowhunger/ortoo/pkg/errors"
	"github.com/knowhunger/ortoo/pkg/iface"
	"github.com/knowhunger/ortoo/pkg/model"
	"github.com/knowhunger/ortoo/pkg/ortoo"
	"github.com/knowhunger/ortoo/server/constants"
	"github.com/knowhunger/ortoo/server/mongodb"
	"github.com/knowhunger/ortoo/server/mongodb/schema"
)

// Manager is a struct that updates snapshot of a datatype in Ortoo server
type Manager struct {
	ctx           context.OrtooContext
	mongo         *mongodb.RepositoryMongo
	datatypeDoc   *schema.DatatypeDoc
	collectionDoc *schema.CollectionDoc
}

// NewManager returns an instance of Snapshot Manager
func NewManager(
	ctx context.OrtooContext,
	mongo *mongodb.RepositoryMongo,
	datatypeDoc *schema.DatatypeDoc,
	collectionDoc *schema.CollectionDoc,
) *Manager {
	return &Manager{
		ctx:           ctx,
		mongo:         mongo,
		datatypeDoc:   datatypeDoc,
		collectionDoc: collectionDoc,
	}
}

// UpdateSnapshot updates snapshot for specified datatype
func (its *Manager) UpdateSnapshot() errors.OrtooError {
	var lastSseq uint64 = 0
	client := ortoo.NewClient(ortoo.NewLocalClientConfig(its.collectionDoc.Name), "server")
	datatype := client.CreateDatatype(its.datatypeDoc.Key, its.datatypeDoc.GetType(), nil).(iface.Datatype)
	datatype.SetLogger(its.ctx.L())
	snapshotDoc, err := its.mongo.GetLatestSnapshot(its.ctx, its.collectionDoc.Num, its.datatypeDoc.DUID)
	if err != nil {
		return err
	}
	if snapshotDoc != nil {
		lastSseq = snapshotDoc.Sseq
		if err := datatype.SetMetaAndSnapshot(snapshotDoc.Meta, snapshotDoc.Snapshot); err != nil {
			return err
		}
	}
	var opList []*model.Operation
	var sseqList []uint64
	opList, sseqList, err = its.mongo.GetOperations(its.ctx, its.datatypeDoc.DUID, lastSseq+1, constants.InfinitySseq)

	if len(sseqList) <= 0 {
		return nil
	}
	if _, err = datatype.ReceiveRemoteModelOperations(opList, false); err != nil {
		// TODO: should fix corruption
		return err
	}
	lastSseq = sseqList[len(sseqList)-1]

	meta, snap, err := datatype.GetMetaAndSnapshot()
	if err != nil {
		return err
	}

	if err := its.mongo.InsertSnapshot(its.ctx, its.collectionDoc.Num, its.datatypeDoc.DUID, lastSseq, meta, snap); err != nil {
		return err
	}

	data := datatype.GetSnapshot().GetAsJSONCompatible()
	if err := its.mongo.InsertRealSnapshot(its.ctx, its.collectionDoc.Name, its.datatypeDoc.Key, data, lastSseq); err != nil {
		return err
	}
	its.ctx.L().Infof("update snapshot and real snapshot")
	return nil
}
