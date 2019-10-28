package datatypes

import (
	"github.com/knowhunger/ortoo/commons/errors"
	// "github.com/knowhunger/ortoo/commons"
	"github.com/knowhunger/ortoo/commons/log"
	"github.com/knowhunger/ortoo/commons/model"
	"sync"
)

// NotUserTransactionTag ...
const NotUserTransactionTag = "NotUserTransactionTag!@#$%ORTOO"

// TransactionDatatypeImpl is the datatype responsible for the transaction.
type TransactionDatatypeImpl struct {
	*WiredDatatypeImpl
	mutex            *sync.RWMutex
	isLocked         bool
	success          bool
	rollbackSnapshot model.Snapshot
	rollbackOps      []model.Operation
	rollbackOpID     *model.OperationID
	currentTrxCtx    *TransactionContext
}

// TransactionDatatype is an interface allowed for transactions.
type TransactionDatatype interface {
	ExecuteTransactionRemote(transaction []model.Operation) error
}

// TransactionContext is a context used for transactions
type TransactionContext struct {
	tag          string
	opBuffer     []model.Operation
	uuid         []byte
	rollbackOpID *model.OperationID
}

// func (t *TransactionContext) GetOpId() *model.OperationID {
//	if len(t.opBuffer) > 0 {
//		return t.opBuffer[0].GetBase().Id
//	}
//	return nil
// }

func (t *TransactionContext) appendOperation(op model.Operation) {
	t.opBuffer = append(t.opBuffer, op)
}

// newTransactionDatatype creates a new TransactionDatatype
func newTransactionDatatype(w *WiredDatatypeImpl, snapshot model.Snapshot) (*TransactionDatatypeImpl, error) {

	return &TransactionDatatypeImpl{
		WiredDatatypeImpl: w,
		mutex:             new(sync.RWMutex),
		isLocked:          false,
		success:           true,
		currentTrxCtx:     nil,
		rollbackSnapshot:  snapshot.CloneSnapshot(),
		rollbackOps:       nil,
		rollbackOpID:      w.opID.Clone(),
	}, nil
}

// ExecuteTransaction is a method to execute a transaction of operations
func (t *TransactionDatatypeImpl) ExecuteTransaction(ctx *TransactionContext, op model.Operation, isLocal bool) (interface{}, error) {
	transactionCtx, err := t.BeginTransaction(NotUserTransactionTag, ctx, false)
	if err != nil {
		return 0, t.Logger.OrtooErrorf(err, "fail to execute transaction")
	}
	defer func() {
		if err := t.EndTransaction(transactionCtx, false); err != nil {
			_ = log.OrtooError(err)
		}
	}()

	if isLocal {
		ret, err := t.executeLocalBase(op)
		if err != nil {
			return 0, t.Logger.OrtooErrorf(err, "fail to execute operation")
		}
		t.currentTrxCtx.appendOperation(op)
		return ret, nil
	} else {
		t.executeRemoteBase(op)
		return nil, nil
	}
}

// make a transaction and lock
func (t *TransactionDatatypeImpl) setTransactionContextAndLock(tag string) *TransactionContext {
	if tag != NotUserTransactionTag {
		t.Logger.Infof("Begin the transaction: `%s`", tag)
	}
	t.mutex.Lock()
	t.isLocked = true
	transactionCtx := &TransactionContext{
		tag:          tag,
		opBuffer:     nil,
		rollbackOpID: t.opID.Clone(),
	}
	return transactionCtx
}

// BeginTransaction is called before a transaction is executed.
// This sets TransactionDatatype.currentTrxCtx, lock, and generates a transaction operation
// This is called in either DoTransaction() or ExecuteTransaction().
// Note that TransactionDatatype.currentTrxCtx is currently working transaction context.
func (t *TransactionDatatypeImpl) BeginTransaction(tag string, trxCtxOfCommonDatatype *TransactionContext, withOp bool) (*TransactionContext, error) {
	if t.isLocked && t.currentTrxCtx == trxCtxOfCommonDatatype {
		return nil, nil // called after DoTransaction() succeeds.
	}
	t.currentTrxCtx = t.setTransactionContextAndLock(tag)
	if withOp {
		op, err := model.NewTransactionOperation(tag)
		if err != nil {
			return nil, t.Logger.OrtooErrorf(err, "fail to generate TransactionBeginOperation")
		}
		t.currentTrxCtx.uuid = op.Uuid
		t.SetNextOpID(op)
		t.currentTrxCtx.appendOperation(op)
	}
	return t.currentTrxCtx, nil
}

// Rollback is called to rollback a transaction
func (t *TransactionDatatypeImpl) Rollback() error {
	t.Logger.Infof("Begin the rollback: '%s'", t.currentTrxCtx.tag)
	snapshotDatatype, _ := t.finalDatatype.(SnapshotDatatype)
	redoOpID := t.opID
	redoSnapshot := snapshotDatatype.GetSnapshot().CloneSnapshot()
	t.SetOpID(t.currentTrxCtx.rollbackOpID)
	snapshotDatatype.SetSnapshot(t.rollbackSnapshot)
	for _, op := range t.rollbackOps {
		err := t.Replay(op)
		if err != nil {
			t.SetOpID(redoOpID)
			snapshotDatatype.SetSnapshot(redoSnapshot)
			return t.Logger.OrtooErrorf(err, "fail to replay operations")
		}
	}
	t.rollbackOpID = t.opID.Clone()
	t.rollbackSnapshot = snapshotDatatype.GetSnapshot().CloneSnapshot()
	t.rollbackOps = nil
	t.Logger.Infof("End the rollback: '%s'", t.currentTrxCtx.tag)
	return nil
}

// SetTransactionFail is called when a transaction fails
func (t *TransactionDatatypeImpl) SetTransactionFail() {
	t.success = false
}

// EndTransaction is called when a transaction ends
func (t *TransactionDatatypeImpl) EndTransaction(trxCtxOfCommonDatatype *TransactionContext, withOp bool) error {
	if trxCtxOfCommonDatatype == t.currentTrxCtx {
		defer t.unlock()
		if t.success {
			if withOp {
				beginOp, ok := t.currentTrxCtx.opBuffer[0].(*model.TransactionOperation)
				if !ok {
					return t.Logger.OrtooErrorf(&errors.ErrTransaction{}, "invalidate transaction: no begin operation")
				}
				beginOp.NumOfOps = uint32(len(t.currentTrxCtx.opBuffer))
			}
			t.rollbackOps = append(t.rollbackOps, t.currentTrxCtx.opBuffer...)
			t.deliverTransaction(t.currentTrxCtx.opBuffer)
			if t.currentTrxCtx.tag != NotUserTransactionTag {
				t.Logger.Infof("End the transaction: `%s`", t.currentTrxCtx.tag)
			}
		} else {
			t.Rollback()
		}
	}
	return nil
}

func (t *TransactionDatatypeImpl) unlock() {
	if t.isLocked {
		t.currentTrxCtx = nil
		t.success = true
		t.mutex.Unlock()
		t.isLocked = false
	}
}

func validateTransaction(transaction []model.Operation) error {
	beginOp, ok := transaction[0].(*model.TransactionOperation)
	if !ok {
		return log.OrtooErrorf(&errors.ErrTransaction{}, "invalidate transaction: no begin transaction")
	}
	if int(beginOp.NumOfOps) != len(transaction) {
		return log.OrtooErrorf(&errors.ErrTransaction{}, "invalidate transaction: incorrect number of operations")
	}
	return nil
}
