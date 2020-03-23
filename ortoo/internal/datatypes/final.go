package datatypes

import (
	"github.com/gogo/protobuf/proto"
	"github.com/knowhunger/ortoo/ortoo/errors"
	"github.com/knowhunger/ortoo/ortoo/log"
	"github.com/knowhunger/ortoo/ortoo/model"
	"github.com/knowhunger/ortoo/ortoo/operations"
)

// FinalDatatype implements the datatype features finally used.
type FinalDatatype struct {
	*TransactionDatatype
	TransactionCtx *TransactionContext
}

// FinalDatatypeInterface defines the interface to obtain FinalDatatype which provide final interfaces.
type FinalDatatypeInterface interface {
	GetFinal() *FinalDatatype
}

// Initialize is a method for initialization
func (c *FinalDatatype) Initialize(
	key string,
	typeOf model.TypeOfDatatype,
	cuid model.CUID,
	w Wire,
	snapshot model.Snapshot,
	datatype model.Datatype) {

	baseDatatype := newBaseDatatype(key, typeOf, cuid)
	wiredDatatype := newWiredDatatype(baseDatatype, w)
	transactionDatatype := newTransactionDatatype(wiredDatatype, snapshot)

	c.TransactionDatatype = transactionDatatype
	c.TransactionCtx = nil
	c.SetDatatype(datatype)
}

// GetMeta returns the binary of metadata of the datatype.
func (c *FinalDatatype) GetMeta() ([]byte, error) {
	meta := model.DatatypeMeta{
		Key:    c.Key,
		DUID:   c.id,
		OpID:   c.opID,
		TypeOf: c.TypeOf,
		State:  c.state,
	}
	metab, err := proto.Marshal(&meta)
	if err != nil {
		return nil, log.OrtooError(err)
	}
	return metab, nil
}

// SetMeta sets the metadata with binary metadata.
func (c *FinalDatatype) SetMeta(meta []byte) error {
	m := model.DatatypeMeta{}
	if err := proto.Unmarshal(meta, &m); err != nil {
		return log.OrtooError(err)
	}
	c.Key = m.Key
	c.id = m.DUID
	c.opID = m.OpID
	c.TypeOf = m.TypeOf
	c.state = m.State
	return nil
}

func (c *FinalDatatype) DoTransaction(tag string, fn func(txnCtx *TransactionContext) error) error {
	txnCtx, err := c.BeginTransaction(tag, c.TransactionCtx, true)
	if err != nil {
		return err
	}
	defer func() {
		if err := c.EndTransaction(txnCtx, true, true); err != nil {
			_ = log.OrtooError(err)
		}
	}()
	if err := fn(txnCtx); err != nil {
		c.SetTransactionFail()
		return errors.NewDatatypeError(errors.ErrDatatypeTransaction, err.Error())
	}
	return nil
}

// SubscribeOrCreate enables a datatype to subscribe and create itself.
func (c *FinalDatatype) SubscribeOrCreate(state model.StateOfDatatype) error {
	if state == model.StateOfDatatype_DUE_TO_SUBSCRIBE {
		c.state = state
		return nil
	}
	subscribeOp, err := operations.NewSnapshotOperation(c.TypeOf, state, c.datatype.GetSnapshot())
	if err != nil {
		return log.OrtooErrorf(err, "fail to subscribe")
	}
	_, err = c.ExecuteOperationWithTransaction(c.TransactionCtx, subscribeOp, true)
	if err != nil {
		return log.OrtooErrorf(err, "fail to execute SubscribeOperation")
	}
	return nil
}

// ExecuteTransactionRemote is a method to execute a transaction of remote operations
func (c FinalDatatype) ExecuteTransactionRemote(transaction []*model.Operation, opList []interface{}) error {
	var trxCtx *TransactionContext
	var err error
	if len(transaction) > 1 {
		trxOp, ok := operations.ModelToOperation(transaction[0]).(*operations.TransactionOperation)
		if !ok {
			return errors.NewDatatypeError(errors.ErrDatatypeTransaction, "no transaction operation")
		}
		if int(trxOp.GetNumOfOps()) != len(transaction) {
			return errors.NewDatatypeError(errors.ErrDatatypeTransaction, "not matched number of operations")
		}
		trxCtx, err = c.BeginTransaction(trxOp.C.Tag, c.TransactionCtx, false)
		if err != nil {
			return log.OrtooError(err)
		}
		defer func() {
			if err := c.EndTransaction(trxCtx, false, false); err != nil {
				_ = log.OrtooError(err)
			}
		}()
		transaction = transaction[1:]
	}
	for _, modelOp := range transaction {
		op := operations.ModelToOperation(modelOp)
		if opList != nil {
			opList = append(opList, op.GetAsJSON())
		}
		c.ExecuteOperationWithTransaction(trxCtx, op, false)
	}
	return nil
}

func validateTransaction(transaction []operations.Operation) error {
	beginOp, ok := transaction[0].(*operations.TransactionOperation)
	if !ok {
		return errors.NewDatatypeError(errors.ErrDatatypeTransaction, "no transaction operation")
	}
	if int(beginOp.GetNumOfOps()) != len(transaction) {
		return errors.NewDatatypeError(errors.ErrDatatypeTransaction, "not matched number of operations")
	}
	return nil
}
