package ortoo

import (
	"encoding/json"
	"fmt"
	"github.com/knowhunger/ortoo/ortoo/errors"
	"github.com/knowhunger/ortoo/ortoo/iface"
	"github.com/knowhunger/ortoo/ortoo/internal/datatypes"
	"github.com/knowhunger/ortoo/ortoo/log"
	"github.com/knowhunger/ortoo/ortoo/model"
	"github.com/knowhunger/ortoo/ortoo/operations"
	"github.com/knowhunger/ortoo/ortoo/types"
	"strings"
)

type List interface {
	Datatype
	ListInTxn
	DoTransaction(tag string, txnFunc func(listTxn ListInTxn) error) error
}

type ListInTxn interface {
	InsertMany(pos int, value ...interface{}) (interface{}, error)
	Get(pos int) (interface{}, error)
	GetMany(pos int, numOfNodes int) ([]interface{}, error)
	Delete(pos int) (interface{}, error)
	DeleteMany(pos int, numOfNodes int) ([]interface{}, error)
	Update(pos int, value ...interface{}) ([]interface{}, error)
	Size() int
}

func newList(key string, cuid types.CUID, wire iface.Wire, handlers *Handlers) List {
	list := &list{
		datatype: &datatype{
			ManageableDatatype: &datatypes.ManageableDatatype{},
			handlers:           handlers,
		},
		snapshot: newListSnapshot(),
	}
	list.Initialize(key, model.TypeOfDatatype_LIST, cuid, wire, list.snapshot, list)
	return list
}

type list struct {
	*datatype
	snapshot *listSnapshot
}

func (its *list) DoTransaction(tag string, txnFunc func(list ListInTxn) error) error {
	return its.ManageableDatatype.DoTransaction(tag, func(txnCtx *datatypes.TransactionContext) error {
		clone := &list{
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

func (its *list) GetAsJSON() interface{} {
	return struct {
		Value []interface{}
	}{
		Value: its.snapshot.GetAsJSON().([]interface{}),
	}
}

func (its *list) ExecuteLocal(op interface{}) (interface{}, error) {
	switch cast := op.(type) {
	case *operations.InsertOperation:
		target, ret, err := its.snapshot.insertLocal(cast.Pos, cast.ID.GetTimestamp(), cast.C.V...)
		if err != nil {
			return nil, err
		}
		cast.C.T = target
		return ret, nil
	case *operations.DeleteOperation:
		deletedTargets, deletedValues, err := its.snapshot.deleteLocal(cast.Pos, cast.NumOfNodes, cast.ID.GetTimestamp())
		if err != nil {
			return nil, err
		}
		cast.C.T = deletedTargets
		return deletedValues, nil
	case *operations.UpdateOperation:
		updatedTargets, updatedValues, err := its.snapshot.updateLocal(cast.Pos, cast.ID.GetTimestamp(), cast.C.V)
		if err != nil {
			return nil, err
		}
		cast.C.T = updatedTargets
		if len(cast.C.T) != len(cast.C.V) {
			return nil, errors.NewDatatypeError(errors.ErrDatatypeIllegalOperation, "not matched")
		}
		return updatedValues, nil
	}
	return nil, errors.NewDatatypeError(errors.ErrDatatypeIllegalOperation, op)
}

func (its *list) ExecuteRemote(op interface{}) (interface{}, error) {
	switch cast := op.(type) {
	case *operations.SnapshotOperation:
		var newSnap = newListSnapshot()
		if err := json.Unmarshal([]byte(cast.C.Snapshot), newSnap); err != nil {
			return nil, errors.NewDatatypeError(errors.ErrDatatypeSnapshot, err.Error())
		}
		its.snapshot = newSnap
		return nil, nil
	case *operations.InsertOperation:
		return its.snapshot.insertRemote(cast.C.T, cast.ID.GetTimestamp(), cast.C.V...)
	case *operations.DeleteOperation:
		return its.snapshot.deleteRemote(cast.C.T, cast.ID.GetTimestamp())
	case *operations.UpdateOperation:
		return its.snapshot.updateRemote(cast.C.T, cast.C.V, cast.ID.GetTimestamp())
	}
	return nil, errors.NewDatatypeError(errors.ErrDatatypeIllegalOperation, op)
}

func (its *list) Size() int {
	return its.Size()
}

func (its *list) GetSnapshot() iface.Snapshot {
	return its.snapshot
}

func (its *list) SetSnapshot(snapshot iface.Snapshot) {
	its.snapshot = snapshot.(*listSnapshot)
}

func (its *list) GetMetaAndSnapshot() ([]byte, iface.Snapshot, error) {
	meta, err := its.ManageableDatatype.GetMeta()
	if err != nil {
		return nil, nil, errors.NewDatatypeError(errors.ErrDatatypeSnapshot, err.Error())
	}
	return meta, its.snapshot, nil
}

func (its *list) SetMetaAndSnapshot(meta []byte, snapshot string) error {
	if err := its.ManageableDatatype.SetMeta(meta); err != nil {
		return errors.NewDatatypeError(errors.ErrDatatypeSnapshot, err.Error())
	}
	if err := json.Unmarshal([]byte(snapshot), its.snapshot); err != nil {
		return errors.NewDatatypeError(errors.ErrDatatypeSnapshot, err.Error())
	}
	return nil
}

func (its *list) Update(pos int, values ...interface{}) ([]interface{}, error) {
	if len(values) < 1 {
		return nil, errors.NewDatatypeError(errors.ErrDatatypeIllegalOperation, "at least one value should be inserted")
	}
	if err := its.snapshot.validateRange(pos, len(values)); err != nil {
		return nil, err
	}
	jsonValues, err := types.ConvertValueList(values)
	if err != nil {
		return nil, errors.NewDatatypeError(errors.ErrDatatypeIllegalOperation, err.Error())
	}
	op := operations.NewUpdateOperation(pos, jsonValues)
	ret, err := its.ExecuteOperationWithTransaction(its.TransactionCtx, op, true)
	if err != nil {
		return nil, err
	}
	return ret.([]interface{}), nil
}

func (its *list) InsertMany(pos int, values ...interface{}) (interface{}, error) {
	if len(values) < 1 {
		return nil, errors.NewDatatypeError(errors.ErrDatatypeIllegalOperation, "at least one value should be inserted")
	}
	jsonValues, err := types.ConvertValueList(values)
	if err != nil {
		return nil, errors.NewDatatypeError(errors.ErrDatatypeIllegalOperation, err.Error())
	}
	op := operations.NewInsertOperation(pos, jsonValues)
	ret, err := its.ExecuteOperationWithTransaction(its.TransactionCtx, op, true)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (its *list) Get(pos int) (interface{}, error) {
	return its.snapshot.get(pos)
}

func (its *list) GetMany(pos int, numOfNodes int) ([]interface{}, error) {
	return its.snapshot.getMany(pos, numOfNodes)
}

// Delete deletes one node at index pos.
func (its *list) Delete(pos int) (interface{}, error) {
	ret, err := its.DeleteMany(pos, 1)
	return ret[0], err
}

// DeleteMany deletes the nodes at index pos in sequence.
func (its *list) DeleteMany(pos int, numOfNode int) ([]interface{}, error) {
	if err := its.snapshot.validateRange(pos, numOfNode); err != nil {
		return nil, err
	}
	op := operations.NewDeleteOperation(pos, numOfNode)
	ret, err := its.ExecuteOperationWithTransaction(its.TransactionCtx, op, true)
	if err != nil {
		return nil, err
	}
	return ret.([]interface{}), nil
}

// ////////////////////////////////////////////////////////////////
//  listSnapshot
// ////////////////////////////////////////////////////////////////

type node struct {
	timedValue
	P    *model.Timestamp
	next *node
	prev *node
}

func (its *node) hash() string {
	return its.getTime().Hash()
}

func (its *node) String() string {
	var sb strings.Builder
	sb.WriteString(its.getTime().ToString())
	if its.P != nil {
		sb.WriteString(its.P.ToString())
	}
	if its.getValue() == nil {
		sb.WriteString(":DELETED")
	} else {
		_, _ = fmt.Fprintf(&sb, ":%v", its.getValue())
	}

	return sb.String()
}

func (its *node) getNextLiveNode() *node {
	ret := its.next
	for ret != nil {
		if ret.getValue() != nil {
			return ret
		}
		ret = ret.next
	}
	return nil
}

type listSnapshot struct {
	head *node
	size int
	Map  map[string]*node
}

func (its *listSnapshot) CloneSnapshot() iface.Snapshot {
	var cloneMap = make(map[string]*node)
	for k, v := range its.Map {
		cloneMap[k] = v
	}
	return &listSnapshot{
		head: its.head,
		size: its.size,
		Map:  cloneMap,
	}
}

func newListSnapshot() *listSnapshot {
	head := &node{
		timedValue: &timedValueImpl{
			V: nil,
			T: model.OldestTimestamp,
		},
		P:    nil,
		prev: nil,
		next: nil,
	}
	m := make(map[string]*node)
	m[head.hash()] = head
	return &listSnapshot{
		head: head,
		Map:  m,
		size: 0,
	}
}

func (its *listSnapshot) insertRemote(pos *model.Timestamp, ts *model.Timestamp, values ...interface{}) (interface{}, error) {
	var tvs []timedValue
	for _, v := range values {
		tvs = append(tvs, &timedValueImpl{
			V: v,
			T: ts.NextDeliminator(),
		})
	}
	return its.insertRemoteWithTimedValue(pos, ts, tvs...)
}

func (its *listSnapshot) insertRemoteWithTimedValue(pos *model.Timestamp, ts *model.Timestamp, tvs ...timedValue) (interface{}, error) {
	if target, ok := its.Map[pos.Hash()]; ok {
		for _, val := range tvs {
			nextTarget := target.next
			for nextTarget != nil && nextTarget.getTime().Compare(ts) > 0 {
				target = target.next
				nextTarget = nextTarget.next
			}
			newNode := &node{
				timedValue: val,
				P:          nil,
				next:       target.next,
				prev:       target,
			}
			target.next = newNode
			its.Map[newNode.hash()] = newNode
			its.size++
			target = newNode
		}
		return nil, nil
	}
	log.Logger.Warnf("no target exists for insertRemote")
	return nil, nil
}

func (its *listSnapshot) appendLocal(ts *model.Timestamp, values ...interface{}) (*model.Timestamp, []interface{}, error) {
	return its.insertLocal(its.size, ts, values...)
}

func (its *listSnapshot) insertLocal(pos int, ts *model.Timestamp, values ...interface{}) (*model.Timestamp, []interface{}, error) {
	var tvs []timedValue
	for _, v := range values {
		tvs = append(tvs, &timedValueImpl{
			V: v,
			T: ts.NextDeliminator(),
		})
	}
	return its.insertLocalWithTimedValue(pos, tvs...)
}

func (its *listSnapshot) insertLocalWithTimedValue(pos int, tvs ...timedValue) (*model.Timestamp, []interface{}, error) {
	if its.size < pos { // size:0 => possible indexes{0} , s:1 => p{0, 1}
		return nil, nil, errors.NewDatatypeError(errors.ErrDatatypeIllegalOperation, "out of bound index")
	}
	var inserted []interface{}
	target := its.findNthTarget(pos)
	targetTs := target.getTime()
	for _, v := range tvs {
		newNode := &node{
			timedValue: v,
			next:       target.next,
			prev:       target,
		}
		target.next = newNode
		its.Map[newNode.hash()] = newNode
		inserted = append(inserted, v.getValue())
		its.size++
		// ts = ts.NextDeliminator()
		target = newNode
	}
	return targetTs, inserted, nil
}

func (its *listSnapshot) isTombstone(n *node) bool {
	if n.getValue() == nil && n.P != nil {
		return true
	}
	return false
}

func (its *listSnapshot) updateLocal(pos int, ts *model.Timestamp, values []interface{}) ([]*model.Timestamp, []interface{}, error) {
	if err := its.validateRange(pos, len(values)); err != nil {
		return nil, nil, err
	}
	var updatedValues []interface{}
	var updatedTargets []*model.Timestamp
	target := its.findNthTarget(pos + 1)
	for _, v := range values {
		updatedValues = append(updatedValues, target.getValue())
		updatedTargets = append(updatedTargets, target.getTime())
		target.timedValue.setValue(v)
		target.P = ts
		target = target.getNextLiveNode()
	}
	return updatedTargets, updatedValues, nil
}

func (its *listSnapshot) updateRemote(targets []*model.Timestamp, values []interface{}, ts *model.Timestamp) (interface{}, error) {
	for i, t := range targets {
		if node, ok := its.Map[t.Hash()]; ok {
			if its.isTombstone(node) {
				continue
			}
			if node.P == nil || node.P.Compare(ts) < 0 {
				node.setValue(values[i])
				node.P = ts
			}
		}
	}
	return nil, nil
}

func (its *listSnapshot) deleteRemote(targets []*model.Timestamp, ts *model.Timestamp) (interface{}, error) {
	for _, t := range targets {
		if node, ok := its.Map[t.Hash()]; ok {
			if !its.isTombstone(node) {
				node.setValue(nil)
				its.size--
				node.P = ts
			} else { // concurrent deletes
				if node.P.Compare(ts) < 0 {
					node.P = ts
				}
			}
		} else {
			log.Logger.Warnf("fail to find delete target: %v", t.ToString())
		}
	}
	return nil, nil
}

func (its *listSnapshot) validateRange(pos int, numOfNodes int) error {
	// 1st condition: if size==4, pos==3 is ok, but 4 is not ok
	// 2nd condition: if size==4, (pos==3, numOfNodes==1) is ok, (pos==3, numOfNodes=2) is not ok.
	if numOfNodes < 1 {
		return errors.NewDatatypeError(errors.ErrDatatypeIllegalOperation, "numOfNodes should be more than 0")
	}
	if its.size-1 < pos || pos+numOfNodes > its.size {
		return errors.NewDatatypeError(errors.ErrDatatypeIllegalOperation, "out of bound index")
	}
	return nil
}

func (its *listSnapshot) deleteLocal(pos, numOfNodes int, ts *model.Timestamp) ([]*model.Timestamp, []interface{}, error) {
	if err := its.validateRange(pos, numOfNodes); err != nil {
		return nil, nil, err
	}
	var deletedTargets []*model.Timestamp
	var deletedValues []interface{}
	target := its.findNthTarget(pos + 1) // no head, but live node
	for i := 0; i < numOfNodes; i++ {
		deletedValues = append(deletedValues, target.getValue())
		deletedTargets = append(deletedTargets, target.getTime())
		target.setValue(nil)
		target.P = ts
		its.size--

		target = target.getNextLiveNode()
	}
	return deletedTargets, deletedValues, nil
}

// for example: h t1 n1 n2 t2 t3 n3 t4 (h:head, n:node, t: tombstone) size==3
// pos : 0 => h : when tombstones follows, the node before them is returned.
// pos : 1 => n1
// pos : 2 => n2
// pos : 3 => n3
func (its *listSnapshot) findNthTarget(pos int) *node {
	ret := its.head
	for i := 1; i <= pos; {
		ret = ret.next
		if !its.isTombstone(ret) { // not tombstone
			i++
		} else { // if tombstone
			for ret.next != nil && its.isTombstone(ret.next) { // while next is tombstone
				ret = ret.next
			}
		}
	}
	return ret
}

func (its *listSnapshot) get(pos int) (interface{}, error) {
	tv, err := its.getTimedValue(pos)
	if err != nil {
		return nil, err
	}
	return tv.getValue(), nil
}

func (its *listSnapshot) getTimedValue(pos int) (timedValue, error) {
	// size == 3, pos can be 0, 1, 2
	if its.size <= pos {
		return nil, errors.NewDatatypeError(errors.ErrDatatypeIllegalOperation, "out of bound index")
	}
	return its.findNthTarget(pos + 1).timedValue, nil
}

func (its *listSnapshot) getMany(pos int, numOfNodes int) ([]interface{}, error) {
	if err := its.validateRange(pos, numOfNodes); err != nil {
		return nil, err
	}
	var ret []interface{}
	for i := 1; i <= numOfNodes; i++ {
		target := its.findNthTarget(pos + i)
		ret = append(ret, target.getValue())
	}
	return ret, nil
}

func (its *listSnapshot) String() string {
	sb := strings.Builder{}
	_, _ = fmt.Fprintf(&sb, "(SIZE:%d) HEAD =>", its.size)
	n := its.head.next
	for n != nil {
		sb.WriteString(n.String())
		n = n.next
		if n != nil {
			sb.WriteString(" => ")
		}
	}
	return sb.String()
}

func (its *listSnapshot) GetAsJSON() interface{} {
	var l []interface{}
	n := its.head.getNextLiveNode()
	for n != nil {
		l = append(l, n.timedValue.getValue())
		n = n.getNextLiveNode()
	}
	return l
}

func (its *listSnapshot) Size() int {
	return its.size
}
