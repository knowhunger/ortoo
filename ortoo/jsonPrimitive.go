package ortoo

import (
	"fmt"
	"github.com/knowhunger/ortoo/ortoo/log"
	"github.com/knowhunger/ortoo/ortoo/model"
	"github.com/knowhunger/ortoo/ortoo/types"
	"reflect"
)

type TypeOfJSON int

const (
	typeJSONPrimitive TypeOfJSON = iota
	TypeJSONElement
	TypeJSONObject
	TypeJSONArray
)

// jsonType extends precededType
// jsonElement extends jsonType
// jsonObject extends jsonType
// jsonArray extends jsonType

// ////////////////////////////////////
//  jsonType
// ////////////////////////////////////

type jsonType interface {
	precededType
	getType() TypeOfJSON
	getRoot() *jsonRoot
	setRoot(r *jsonObject)
	getParent() jsonType
	getParentAsJSONObject() *jsonObject
	findJSONArray(ts *model.Timestamp) (j *jsonArray, ok bool)
	findJSONObject(ts *model.Timestamp) (j *jsonObject, ok bool)
	findJSONElement(ts *model.Timestamp) (j *jsonElement, ok bool)
	findJSONPrimitive(ts *model.Timestamp) (j jsonType, ok bool)
	addToNodeMap(primitive jsonType)
	addToCemetery(primitive jsonType)
	createJSONObject(parent jsonType, value interface{}, ts *model.Timestamp) *jsonObject
	createJSONArray(parent jsonType, value interface{}, ts *model.Timestamp) *jsonArray
	toJSONPrimitiveForMarshal() *jsonPrimitiveForMarshal
}

type jsonRoot struct {
	root     *jsonObject
	nodeMap  map[string]jsonType
	cemetery map[string]jsonType
}

type jsonPrimitive struct {
	root    *jsonRoot
	parent  jsonType
	K       *model.Timestamp // used for key that is immutable and used in the root
	P       *model.Timestamp // used for precedence; for example makeTomb
	deleted bool
}

func (its *jsonPrimitive) getType() TypeOfJSON {
	return typeJSONPrimitive
}

func (its *jsonPrimitive) isTomb() bool {
	return its.deleted
}

func (its *jsonPrimitive) makeTomb(ts *model.Timestamp) bool {
	if its.deleted {
		if its.P.Compare(ts) > 0 { // This condition makes newer timestamps remain in nodes.
			log.Logger.Infof("fail to makeTomb() of jsonPrimitive:%v", its.K.ToString())
			return false
		}
	}
	its.P = ts
	its.deleted = true
	log.Logger.Infof("makeTomb() of jsonPrimitive:%v", its.K.ToString())
	return true
}

func (its *jsonPrimitive) getTime() *model.Timestamp {
	return its.K
}

func (its *jsonPrimitive) setTime(ts *model.Timestamp) {
	its.K = ts
}

func (its *jsonPrimitive) getPrecedence() *model.Timestamp {
	return its.P
}

func (its *jsonPrimitive) setPrecedence(ts *model.Timestamp) {
	its.P = ts
}

func (its *jsonPrimitive) findJSONPrimitive(ts *model.Timestamp) (j jsonType, ok bool) {
	node, ok := its.getRoot().nodeMap[ts.Hash()]
	return node, ok
}

func (its *jsonPrimitive) findJSONElement(ts *model.Timestamp) (j *jsonElement, ok bool) {
	if node, ok := its.getRoot().nodeMap[ts.Hash()]; ok {
		if j, ok2 := node.(*jsonElement); ok2 {
			return j, ok2
		}
	}
	return nil, false
}

func (its *jsonPrimitive) findJSONObject(ts *model.Timestamp) (json *jsonObject, ok bool) {
	if node, ok := its.getRoot().nodeMap[ts.Hash()]; ok {
		if j, ok2 := node.(*jsonObject); ok2 {
			return j, ok2
		}
	}
	return nil, false
}

func (its *jsonPrimitive) findJSONArray(ts *model.Timestamp) (json *jsonArray, ok bool) {
	if node, ok := its.getRoot().nodeMap[ts.Hash()]; ok {
		if j, ok2 := node.(*jsonArray); ok2 {
			return j, ok2
		}
	}
	return nil, false
}

func (its *jsonPrimitive) addToNodeMap(primitive jsonType) {
	its.getRoot().nodeMap[primitive.getTime().Hash()] = primitive
}

func (its *jsonPrimitive) addToCemetery(primitive jsonType) {
	its.getRoot().cemetery[primitive.getTime().Hash()] = primitive
}

func (its *jsonPrimitive) getValue() types.JSONValue {
	panic("should be overridden")
}

func (its *jsonPrimitive) setValue(v types.JSONValue) {
	panic("should be overridden")
}

func (its *jsonPrimitive) getRoot() *jsonRoot {
	return its.root
}

func (its *jsonPrimitive) setRoot(r *jsonObject) {
	its.root.root = r
	its.root.nodeMap[r.getTime().Hash()] = r
}

func (its *jsonPrimitive) getParent() jsonType {
	return its.parent
}

func (its *jsonPrimitive) getParentAsJSONObject() *jsonObject {
	return its.parent.(*jsonObject)
}

func (its *jsonPrimitive) String() string {
	return fmt.Sprintf("%x", &its.parent)
}

func (its *jsonPrimitive) createJSONArray(parent jsonType, value interface{}, ts *model.Timestamp) *jsonArray {
	ja := newJSONArray(parent, ts.NextDeliminator())
	target := reflect.ValueOf(value)
	var appendValues []precededType
	for i := 0; i < target.Len(); i++ {
		field := target.Index(i)
		switch field.Kind() {
		case reflect.Slice, reflect.Array:
			ja := its.createJSONArray(ja, field.Interface(), ts)
			appendValues = append(appendValues, ja)
		case reflect.Struct, reflect.Map:
			childJO := its.createJSONObject(ja, field.Interface(), ts)
			appendValues = append(appendValues, childJO)
		case reflect.Ptr:
			val := field.Elem()
			its.createJSONArray(parent, val.Interface(), ts)
		default:
			element := newJSONElement(ja, types.ConvertToJSONSupportedValue(field.Interface()), ts.NextDeliminator())
			appendValues = append(appendValues, element)
		}
	}
	if appendValues != nil {
		ja.insertLocalWithPrecededTypes(0, appendValues...)
		for _, v := range appendValues {
			its.addToNodeMap(v.(jsonType))
		}
	}
	// log.Logger.Infof("%v", ja.String())
	return ja
}

func (its *jsonPrimitive) createJSONObject(parent jsonType, value interface{}, ts *model.Timestamp) *jsonObject {
	jo := newJSONObject(parent, ts.NextDeliminator())
	target := reflect.ValueOf(value)
	fields := reflect.TypeOf(value)

	if target.Kind() == reflect.Map {
		mapValue := value.(map[string]interface{})
		for k, v := range mapValue {
			val := reflect.ValueOf(v)
			its.addValueToJSONObject(jo, k, val, ts)
		}
	} else {
		for i := 0; i < target.NumField(); i++ {
			value := target.Field(i)
			its.addValueToJSONObject(jo, fields.Field(i).Name, value, ts)
		}
	}

	return jo
}

func (its *jsonPrimitive) addValueToJSONObject(jo *jsonObject, key string, value reflect.Value, ts *model.Timestamp) {
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		ja := its.createJSONArray(jo, value.Interface(), ts)
		jo.putCommonWithTimedValue(key, ja)
		its.addToNodeMap(ja)
	case reflect.Struct, reflect.Map:
		childJO := its.createJSONObject(jo, value.Interface(), ts)
		jo.putCommonWithTimedValue(key, childJO)
		its.addToNodeMap(childJO)
	case reflect.Ptr:
		val := value.Elem()
		its.createJSONObject(jo, val.Interface(), ts)
	default:
		element := newJSONElement(jo, types.ConvertToJSONSupportedValue(value.Interface()), ts.NextDeliminator())
		jo.putCommonWithTimedValue(key, element)
		its.addToNodeMap(element)
	}
}
