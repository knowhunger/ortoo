package ortoo

import (
	"fmt"
	"github.com/knowhunger/ortoo/ortoo/model"
	"github.com/knowhunger/ortoo/ortoo/types"
)

type timedValue interface {
	getValue() types.JSONValue
	setValue(v types.JSONValue)
	getTime() *model.Timestamp
	setTime(ts *model.Timestamp)
	makeTomb(ts *model.Timestamp) bool
	isTomb() bool
	String() string
}

type timedValueImpl struct {
	V types.JSONValue
	T *model.Timestamp
	P *model.Timestamp
}

func (its *timedValueImpl) getValue() types.JSONValue {
	return its.V
}

func (its *timedValueImpl) setValue(v types.JSONValue) {
	its.V = v
}

func (its *timedValueImpl) getTime() *model.Timestamp {
	return its.T
}

func (its *timedValueImpl) setTime(ts *model.Timestamp) {
	its.T = ts
}

// this is for hash_map
func (its *timedValueImpl) makeTomb(ts *model.Timestamp) bool {
	its.V = nil
	its.T = ts
	return true
}

func (its *timedValueImpl) isTomb() bool {
	if its.V == nil {
		return true
	}
	return false
}

func (its *timedValueImpl) String() string {
	if its.V == nil {
		return fmt.Sprintf("Φ|%s", its.T.ToString())
	}
	return fmt.Sprintf("TV[%v|T%s]", its.V, its.T.ToString())
}