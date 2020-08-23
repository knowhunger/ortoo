package ortoo

import (
	"encoding/json"
	"github.com/knowhunger/ortoo/ortoo/iface"
	"github.com/knowhunger/ortoo/ortoo/operations"
	"github.com/knowhunger/ortoo/ortoo/types"

	"fmt"
	"github.com/knowhunger/ortoo/ortoo/errors"
	"github.com/knowhunger/ortoo/ortoo/internal/datatypes"
	"github.com/knowhunger/ortoo/ortoo/log"
	"github.com/knowhunger/ortoo/ortoo/model"
)

// Counter is an Ortoo datatype which provides int counter interfaces.
type Counter interface {
	Datatype
	CounterInTxn
	DoTransaction(tag string, txnFunc func(intCounter CounterInTxn) error) error
}

// CounterInTxn is an Ortoo datatype which provides int counter interfaces in a transaction.
type CounterInTxn interface {
	Get() int32
	Increase() (int32, error)
	IncreaseBy(delta int32) (int32, error)
}

type counter struct {
	*datatype
	snapshot *counterSnapshot
}

// newCounter creates a new int counter
func newCounter(key string, cuid types.CUID, wire iface.Wire, handler *Handlers) Counter {
	base := datatypes.NewBaseDatatype(key, model.TypeOfDatatype_COUNTER, cuid)
	counter := &counter{
		datatype: &datatype{
			ManageableDatatype: &datatypes.ManageableDatatype{},
			handlers:           handler,
		},
		snapshot: &counterSnapshot{
			base:  base,
			Value: 0,
		},
	}
	counter.Initialize(base, wire, counter.snapshot, counter)
	return counter
}

func (its *counter) DoTransaction(tag string, txnFunc func(intCounter CounterInTxn) error) error {
	return its.ManageableDatatype.DoTransaction(tag, func(txnCtx *datatypes.TransactionContext) error {
		clone := &counter{
			datatype: &datatype{
				ManageableDatatype: &datatypes.ManageableDatatype{
					TransactionDatatype: its.ManageableDatatype.TransactionDatatype,
					TransactionCtx:      txnCtx,
				},
				handlers: its.handlers,
			},
			snapshot: its.snapshot,
		}
		return txnFunc(clone)
	})
}

func (its *counter) GetFinal() *datatypes.ManageableDatatype {
	return its.ManageableDatatype
}

// ExecuteLocal enables the operation to perform something at the local client.
func (its *counter) ExecuteLocal(op interface{}) (interface{}, errors.OrtooError) {
	iop := op.(*operations.IncreaseOperation)
	return its.snapshot.increaseCommon(iop.C.Delta), nil
}

// ExecuteRemote is called by operation.ExecuteRemote()
func (its *counter) ExecuteRemote(op interface{}) (interface{}, errors.OrtooError) {
	switch cast := op.(type) {
	case *operations.SnapshotOperation:
		newSnap := counterSnapshot{}
		if err := json.Unmarshal([]byte(cast.C.Snapshot), &newSnap); err != nil {
			return nil, errors.ErrDatatypeSnapshot.New(its.Logger, err.Error())
		}
		its.snapshot = &newSnap
		return nil, nil
	case *operations.IncreaseOperation:
		return its.snapshot.increaseCommon(cast.C.Delta), nil
	}

	return nil, errors.ErrDatatypeIllegalOperation.New(its.Logger, op)
}

func (its *counter) Get() int32 {
	return its.snapshot.Value
}

func (its *counter) Increase() (int32, error) {
	return its.IncreaseBy(1)
}

func (its *counter) IncreaseBy(delta int32) (int32, error) {
	op := operations.NewIncreaseOperation(delta)
	ret, err := its.ExecuteOperationWithTransaction(its.TransactionCtx, op, true)
	if err != nil {
		return 0, log.OrtooErrorf(err, "fail to execute operation")
	}
	return ret.(int32), nil
}

func (its *counter) GetSnapshot() iface.Snapshot {
	return its.snapshot
}

func (its *counter) SetSnapshot(snapshot iface.Snapshot) {
	its.snapshot = snapshot.(*counterSnapshot)
}

func (its *counter) GetAsJSON() interface{} {
	return its.snapshot.GetAsJSONCompatible()
}

func (its *counter) GetMetaAndSnapshot() ([]byte, iface.Snapshot, errors.OrtooError) {
	meta, err := its.ManageableDatatype.GetMeta()
	if err != nil {
		return nil, nil, errors.ErrDatatypeSnapshot.New(its.Logger, err.Error())
	}
	return meta, its.snapshot, nil
}

func (its *counter) SetMetaAndSnapshot(meta []byte, snapshot string) errors.OrtooError {
	if err := its.ManageableDatatype.SetMeta(meta); err != nil {
		return errors.ErrDatatypeSnapshot.New(its.Logger, err.Error())
	}

	if err := json.Unmarshal([]byte(snapshot), its.snapshot); err != nil {
		return errors.ErrDatatypeSnapshot.New(its.Logger, err.Error())
	}
	return nil
}

// ////////////////////////////////////////////////////////////////
//  counterSnapshot
// ////////////////////////////////////////////////////////////////

type counterSnapshot struct {
	base  *datatypes.BaseDatatype
	Value int32 `json:"value"`
}

func (its *counterSnapshot) CloneSnapshot() iface.Snapshot {
	return &counterSnapshot{
		Value: its.Value,
	}
}

func (its *counterSnapshot) GetAsJSONCompatible() interface{} {
	return its
}

func (its *counterSnapshot) increaseCommon(delta int32) int32 {
	its.Value = its.Value + delta
	return its.Value
}

func (its *counterSnapshot) String() string {
	return fmt.Sprintf("Map: %d", its.Value)
}
