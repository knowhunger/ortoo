package model

type Operationer interface {
	ExecuteLocal(datatype OperationExecuter) (interface{}, error)
	ExecuteRemote(datatype OperationExecuter) (interface{}, error)
	GetBase() *BaseOperation
}

type OperationExecuter interface {
	ExecuteLocal(op interface{}) (interface{}, error)
	ExecuteRemote(op interface{}) (interface{}, error)
}

func NewOperation(opType TypeOperation) *BaseOperation {
	return &BaseOperation{
		Id:     NewOperationID(),
		OpType: opType,
	}
}

func (o *BaseOperation) SetOperationID(opID *OperationID) {
	o.Id = opID
}

//////////////////// IncreaseOperation ////////////////////

func NewIncreaseOperation(delta int32) *IncreaseOperation {
	return &IncreaseOperation{
		Base:  NewOperation(TypeOperation_INT_COUNTER_INCREASE),
		Delta: delta,
	}
}

func (i *IncreaseOperation) ExecuteLocal(datatype OperationExecuter) (interface{}, error) {
	return datatype.ExecuteLocal(i)
}

func (i *IncreaseOperation) ExecuteRemote(datatype OperationExecuter) (interface{}, error) {
	return datatype.ExecuteRemote(i)
}

func ToOperation(op Operationer) *Operation {
	switch o := op.(type) {
	case *IncreaseOperation:
		return &Operation{Body: &Operation_IncreaseOperation{o}}
	}
	return nil
}

func ToOperationer(op *Operation) Operationer {
	switch o := op.Body.(type) {
	case *Operation_IncreaseOperation:
		return o.IncreaseOperation
	}
	return nil
}