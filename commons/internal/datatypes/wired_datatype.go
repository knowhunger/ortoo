package datatypes

import (
	"github.com/knowhunger/ortoo/commons/internal/constants"
	"github.com/knowhunger/ortoo/commons/log"
	"github.com/knowhunger/ortoo/commons/model"
)

type WiredDatatypeImpl struct {
	Wire
	*baseDatatype
	//*transactionalDatatype
	trans      *WiredDatatypeImpl
	checkPoint *model.CheckPoint
	buffer     []*model.OperationOnWire
	opExecuter model.OperationExecuter
}

type WiredDatatyper interface {
	GetWired() WiredDatatype
}

type PublicWiredDatatypeInterface interface {
	PublicTransactionalDatatypeInterface
}

type WiredDatatype interface {
	GetBase() *baseDatatype
	ExecuteRemote(op model.Operation)
	ExecuteTransactionRemote(transaction []model.Operation)
	CreatePushPullPack() *model.PushPullPack
	ApplyPushPullPack(*model.PushPullPack)
}

func NewWiredDataType(t model.TypeDatatype, w Wire) (*WiredDatatypeImpl, error) {
	baseDatatype, err := newBaseDatatype(t)
	if err != nil {
		return nil, log.OrtooError(err, "fail to create wiredDatatype due to baseDatatype")
	}
	return &WiredDatatypeImpl{
		baseDatatype: baseDatatype,
		checkPoint:   model.NewCheckPoint(),
		buffer:       make([]*model.OperationOnWire, 0, constants.OperationBufferSize),
		Wire:         w,
	}, nil
}

func (w *WiredDatatypeImpl) ExecuteWired(op model.Operation) (ret interface{}, err error) {
	ret, err = w.executeLocalBase(w.opExecuter, op)
	if err != nil {
		return ret, log.OrtooError(err, "fail to executeLocalTransactional")
	}

	w.buffer = append(w.buffer, model.ToOperationOnWire(op))
	w.DeliverOperation(w, op)
	return ret, nil
}

func (w *WiredDatatypeImpl) GetBase() *baseDatatype {
	return w.baseDatatype
}

func (w *WiredDatatypeImpl) String() string {
	return w.baseDatatype.String()
}

func (w *WiredDatatypeImpl) SetOperationExecuter(opExecuter model.OperationExecuter) {
	w.opExecuter = opExecuter
}

func (w *WiredDatatypeImpl) ExecuteRemote(op model.Operation) {
	w.opID.SyncLamport(op.GetBase().GetId().Lamport)
	w.GetBase().executeRemoteBase(w.opExecuter, op)
}

func (w *WiredDatatypeImpl) ExecuteTransactionRemote(transaction []model.Operation) {
	if err := validateTransaction(transaction); err != nil {
		return
	}
	for _, op := range transaction {
		w.opID.SyncLamport(op.GetBase().GetId().Lamport)
		w.GetBase().executeRemoteBase(w.opExecuter, op)
	}
}

func (w *WiredDatatypeImpl) CreatePushPullPack() *model.PushPullPack {
	seq := w.checkPoint.Cseq
	operations := w.getOperationOnWires(seq + 1)
	cp := &model.CheckPoint{
		Sseq: w.checkPoint.GetSseq(),
		Cseq: w.checkPoint.GetCseq() + uint64(len(operations)),
	}
	return &model.PushPullPack{
		CheckPoint: cp,
		Duid:       w.id,
		Era:        0,
		Type:       0,
		Operations: operations,
	}
}

func (w *WiredDatatypeImpl) ApplyPushPullPack(ppp *model.PushPullPack) {

	opList := ppp.GetOperations()
	for _, op := range opList {
		w.ExecuteRemote(model.ToOperation(op))
	}

}

func (w *WiredDatatypeImpl) getOperationOnWires(cseq uint64) []*model.OperationOnWire {
	op := model.ToOperation(w.buffer[0])
	startCseq := op.GetBase().Id.GetSeq()
	var start = int(cseq - startCseq)
	if len(w.buffer) > start {
		return w.buffer[start:]
	}
	return []*model.OperationOnWire{}

}

//
//func (w *WiredDatatypeImpl) BeginTransactionOnWired(tag string) error {
//	if err := w.BeginTransactionLocal(tag); err != nil {
//		return log.OrtooError(err, "fail to begin transaction on wired")
//	}
//	return nil
//}
//
//func (w *WiredDatatypeImpl) EndTransactionOnWired() {
//	buffer := w.EndTransactionLocal()
//	if buffer == nil {
//		return
//	}
//	w.DeliverTransaction(w, buffer)
//	w.EndTransaction()
//}
//
//func (w *WiredDatatypeImpl) GetTransactionalWiredDatatypeImpl() *WiredDatatypeImpl {
//	ww := &WiredDatatypeImpl{
//		trans:      w,
//		opExecuter: w.opExecuter,
//	}
//	return ww
//}
