package operations

import (
	"github.com/knowhunger/ortoo/pkg/errors"
	"github.com/knowhunger/ortoo/pkg/iface"
	"github.com/knowhunger/ortoo/pkg/model"
)

// NewInsertOperation creates a new InsertOperation
func NewInsertOperation(pos int, values []interface{}) *InsertOperation {
	return &InsertOperation{
		baseOperation: newBaseOperation(nil),
		Pos:           pos,
		C: insertContent{
			V: values,
		},
	}
}

type insertContent struct {
	T *model.Timestamp
	V []interface{}
}

// InsertOperation is used to insert a value to a list
type InsertOperation struct {
	*baseOperation
	Pos int // for local
	C   insertContent
}

// ExecuteLocal enables the operation to perform something at the local client.
func (its *InsertOperation) ExecuteLocal(datatype iface.Datatype) (interface{}, errors.OrtooError) {
	return datatype.ExecuteLocal(its)
}

// ExecuteRemote enables the operation to perform something at the remote clients.
func (its *InsertOperation) ExecuteRemote(datatype iface.Datatype) (interface{}, errors.OrtooError) {
	return datatype.ExecuteRemote(its)
}

// ToModelOperation transforms this operation to the model.Operation.
func (its *InsertOperation) ToModelOperation() *model.Operation {
	return &model.Operation{
		ID:     its.ID,
		OpType: model.TypeOfOperation_LIST_INSERT,
		Body:   marshalContent(its.C),
	}
}

// GetType returns the type of operation.
func (its *InsertOperation) GetType() model.TypeOfOperation {
	return model.TypeOfOperation_LIST_INSERT
}

func (its *InsertOperation) String() string {
	return toString(its.ID, its.C)
}

// ////////////////// DeleteOperation ////////////////////

// NewDeleteOperation creates a new DeleteOperation.
func NewDeleteOperation(pos int, numOfNodes int) *DeleteOperation {
	return &DeleteOperation{
		baseOperation: newBaseOperation(nil),
		Pos:           pos,
		NumOfNodes:    numOfNodes,
		C:             deleteContent{},
	}
}

type deleteContent struct {
	T []*model.Timestamp
}

// DeleteOperation is used to delete a value from a list.
type DeleteOperation struct {
	*baseOperation
	Pos        int
	NumOfNodes int
	C          deleteContent
}

// ExecuteLocal enables the operation to perform something at the local client.
func (its *DeleteOperation) ExecuteLocal(datatype iface.Datatype) (interface{}, errors.OrtooError) {
	return datatype.ExecuteLocal(its)
}

// ExecuteRemote enables the operation to perform something at the remote clients.
func (its *DeleteOperation) ExecuteRemote(datatype iface.Datatype) (interface{}, errors.OrtooError) {
	return datatype.ExecuteRemote(its)
}

// ToModelOperation transforms this operation to the model.Operation.
func (its *DeleteOperation) ToModelOperation() *model.Operation {
	return &model.Operation{
		ID:     its.ID,
		OpType: model.TypeOfOperation_LIST_DELETE,
		Body:   marshalContent(its.C),
	}
}

// GetType returns the type of operation.
func (its *DeleteOperation) GetType() model.TypeOfOperation {
	return model.TypeOfOperation_LIST_DELETE
}

func (its *DeleteOperation) String() string {
	return toString(its.ID, its.C)
}

// ////////////////// UpdateOperation ////////////////////

// NewUpdateOperation creates a new UpdateOperation.
func NewUpdateOperation(pos int, values []interface{}) *UpdateOperation {
	return &UpdateOperation{
		baseOperation: newBaseOperation(nil),
		Pos:           pos,
		C: updateContent{
			V: values,
		},
	}
}

type updateContent struct {
	T []*model.Timestamp
	V []interface{}
}

// UpdateOperation is used to update a value in a list.
type UpdateOperation struct {
	*baseOperation
	Pos int
	C   updateContent
}

// ExecuteLocal enables the operation to perform something at the local client.
func (its *UpdateOperation) ExecuteLocal(datatype iface.Datatype) (interface{}, errors.OrtooError) {
	return datatype.ExecuteLocal(its)
}

// ExecuteRemote enables the operation to perform something at the remote clients.
func (its *UpdateOperation) ExecuteRemote(datatype iface.Datatype) (interface{}, errors.OrtooError) {
	return datatype.ExecuteRemote(its)
}

// ToModelOperation transforms this operation to the model.Operation.
func (its *UpdateOperation) ToModelOperation() *model.Operation {
	return &model.Operation{
		ID:     its.ID,
		OpType: model.TypeOfOperation_LIST_UPDATE,
		Body:   marshalContent(its.C),
	}
}

// GetType returns the type of operation.
func (its *UpdateOperation) GetType() model.TypeOfOperation {
	return model.TypeOfOperation_LIST_UPDATE
}

func (its *UpdateOperation) String() string {
	return toString(its.ID, its.C)
}
