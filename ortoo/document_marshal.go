package ortoo

import (
	"encoding/json"
	"github.com/knowhunger/ortoo/ortoo/log"
	"github.com/knowhunger/ortoo/ortoo/model"
)

type marshaledDocument struct {
	TS        map[string]*model.Timestamp   `json:"ts"`
	JSONTypes map[string]*marshaledJSONType `json:"js"`
	Cemetery  []string                      `json:"ce"` // store the hash of all deleted jsonType
}

type marshalJSONType string

const (
	mJSONElement marshalJSONType = "E"
	mJSONObject  marshalJSONType = "O"
	mJSONArray   marshalJSONType = "A"
)

type marshaledJSONType struct {
	F *model.Timestamp     `json:"f,omitempty"` // jsonPrimitive.parent (Forebear)
	K *model.Timestamp     `json:"k,omitempty"` // jsonPrimitive.K
	P *model.Timestamp     `json:"p,omitempty"` // jsonPrimitive.P
	D bool                 `json:"d,omitempty"` // jsonPrimitive.deleted
	T marshalJSONType      `json:"t,omitempty"` // type; "E": jsonElement, "O": jsonObject, "A": jsonArray
	E interface{}          `json:"e,omitempty"` // for jsonElement
	A *marshaledJSONArray  `json:"a,omitempty"` // for jsonArray
	O *marshaledJSONObject `json:"o,omitempty"` // for jsonObject
}

type marshaledJSONObject struct {
	M map[string]*model.Timestamp `json:"m"`
	S int                         `json:"s"`
}

type marshaledJSONArray struct {
	N []*model.Timestamp `json:"m"`
	S int                `json:"s"`
}

func (its *marshaledDocument) unifyTimestamp(ts *model.Timestamp) *model.Timestamp {
	if ts == nil {
		return nil
	}
	return its.TS[ts.Hash()]
}

func newMarshaledDocument() *marshaledDocument {
	return &marshaledDocument{
		Cemetery:  nil,
		TS:        make(map[string]*model.Timestamp),
		JSONTypes: make(map[string]*marshaledJSONType), // by marshal()
	}
}

// //////////////////////////////
// Marshal functions
// //////////////////////////////

// MarshalJSON returns marshaledDocument.
func (its *jsonObject) MarshalJSON() ([]byte, error) {
	marshalDoc := newMarshaledDocument()
	for k, v := range its.getRoot().nodeMap {
		t, p := v.getTime(), v.getPrecedence()
		if t != nil {
			marshalDoc.TS[k] = t
		}
		if p != nil {
			marshalDoc.TS[k] = p
		}
		marshaled := v.marshal()
		marshalDoc.JSONTypes[k] = marshaled
	}
	for k, _ := range its.getRoot().cemetery {
		marshalDoc.Cemetery = append(marshalDoc.Cemetery, k)
	}
	return json.Marshal(marshalDoc)
}

func (its *jsonPrimitive) marshal() *marshaledJSONType {
	var forebear *model.Timestamp = nil
	if its.parent != nil {
		forebear = its.parent.getTime()
	}
	return &marshaledJSONType{
		F: forebear,
		K: its.K,
		P: its.P,
		D: its.deleted,
	}
}

func (its *jsonElement) marshal() *marshaledJSONType {
	forMarshal := its.jsonType.marshal()
	forMarshal.T = mJSONElement
	forMarshal.E = its.getValue()
	return forMarshal
}

func (its *jsonObject) marshal() *marshaledJSONType {
	marshal := its.jsonType.marshal()
	marshal.T = mJSONObject
	value := &marshaledJSONObject{
		S: its.Size,
		M: make(map[string]*model.Timestamp),
	}
	for k, v := range its.hashMapSnapshot.Map {
		jsonP := v.(jsonType)
		value.M[k] = jsonP.getTime()
	}
	marshal.O = value
	return marshal
}

func (its *jsonArray) marshal() *marshaledJSONType {
	marshal := its.jsonType.marshal()
	marshal.T = mJSONArray
	value := &marshaledJSONArray{
		S: its.listSnapshot.size,
	}
	n := its.listSnapshot.head.getNext()
	for n != nil {
		value.N = append(value.N, n.getTime())
		n = n.getNext()
	}
	marshal.A = value
	return marshal
}

//
// func (its *jsonElement) MarshalJSON() ([]byte, error) {
// 	j := its.marshal()
// 	return json.Marshal(j)
// }
//
// func (its *jsonElement) UnmarshalJSON(bytes []byte) error {
// 	// forMarshal := marshaledJSONType{}
// 	// if err := json.Unmarshal(bytes, &forMarshal); err != nil {
// 	// 	return log.OrtooError(err)
// 	// }
// 	//
// 	// its.V = forMarshal.V
// 	// its.jsonType = &jsonPrimitive{
// 	// 	parent:  nil,
// 	// 	common:    nil,
// 	// 	K:       forMarshal.K,
// 	// 	P:       forMarshal.P,
// 	// 	deleted: false,
// 	// }
// 	return nil
// }
//
// func (its *jsonElement) marshalDocument(forMarshal *marshaledDocument) {
// 	k := its.jsonType.getTime()
// 	if k != nil {
// 		forMarshal.TS[k.Hash()] = k
// 	}
// 	p := its.jsonType.getPrecedence()
// 	if p != nil {
// 		forMarshal.TS[p.Hash()] = p
// 	}
// 	m := its.marshal()
// 	forMarshal.JSONTypes[m.K.Hash()] = m
// }
//
// func (its *jsonArray) marshalDocument(forMarshal *marshaledDocument) {
// 	var sb strings.Builder
// 	n := its.head
// 	for n != nil {
// 		if n.getPrev() != nil {
// 			sb.WriteString(n.getPrev().hash())
// 		}
// 		if n.getNext() != nil {
// 			sb.WriteString(n.getNext().hash())
// 		}
// 		// switch cast:= n.(type) {
// 		// case *jsonObject:
// 		// case *jsonElement:
// 		// case *jsonArray:
// 		// }
// 		n = n.getNext()
// 	}
// }
//
// func (its *jsonObject) marshalDocument(forMarshal *marshaledDocument) {
// 	var sb strings.Builder
//
// 	for k, v := range its.Map {
// 		if v != nil {
// 			switch cast := v.(type) {
// 			case *jsonObject:
// 				cast.marshalDocument(forMarshal)
// 			case *jsonElement:
// 				cast.marshalDocument(forMarshal)
// 				// kts := cast.jsonType.getKey()
// 				// pts := cast.jsonType.getTime()
// 				//
// 				sb.WriteString(k)
// 				sb.WriteString(" ")
// 				j, err := json.Marshal(cast)
// 				if err != nil {
// 					return
// 				}
// 				sb.WriteString(string(j))
// 			case *jsonArray:
// 				cast.marshalDocument(forMarshal)
//
// 			}
// 		}
// 		sb.WriteString("\n")
// 	}
// 	log.Logger.Infof("%v", sb.String())
// 	// b, err := json.Marshal(forMarshal)
// 	// if err != nil {
// 	// 	return
// 	// }
// 	return
// }

// //////////////////////////////
// Unmarshal functions
// //////////////////////////////

func (its *jsonObject) UnmarshalJSON(bytes []byte) error {
	var forUnmarshal marshaledDocument

	// jsonTypeMap := make(map[string]jsonType)

	if err := json.Unmarshal(bytes, &forUnmarshal); err != nil {
		return log.OrtooError(err)
	}

	root := forUnmarshal.JSONTypes[model.OldestTimestamp.Hash()]

	common := &jsonCommon{
		root:     its,
		nodeMap:  make(map[string]jsonType),
		cemetery: make(map[string]jsonType),
	}

	its.jsonType = &jsonPrimitive{
		common:  common,
		parent:  nil, // root has no parent
		deleted: root.D,
		K:       forUnmarshal.unifyTimestamp(root.K), // rootKey
		P:       forUnmarshal.unifyTimestamp(root.P),
	}
	rootKey := its.getTime().Hash()
	common.nodeMap[rootKey] = its

	// make all skeleton jsonTypes "in advance"
	for k, v := range forUnmarshal.JSONTypes {
		if k != rootKey {
			common.nodeMap[k] = v.unmarshalAsJSONType(&forUnmarshal, common, nil)
		}
	}

	// constitute the value for each jsonType
	for k, v := range forUnmarshal.JSONTypes {
		j := common.nodeMap[k] // real jsonType
		if v.F != nil {
			parent := common.nodeMap[v.F.Hash()]
			j.setParent(parent)
		}
		j.unmarshal(v, common.nodeMap) // j is a real jsonType and v is marshaledJSONType
	}
	return nil
}

func (its *jsonObject) unmarshal(marshaled *marshaledJSONType, jsonMap map[string]jsonType) {

	jo := marshaled.O

	its.hashMapSnapshot = &hashMapSnapshot{
		Map:  make(map[string]timedType),
		Size: jo.S,
	}
	for k, v := range jo.M {
		realInstance := jsonMap[v.Hash()]
		its.hashMapSnapshot.Map[k] = realInstance
	}
}

func (its *jsonArray) unmarshal(marshaled *marshaledJSONType, jsonMap map[string]jsonType) {
	ja := marshaled.A

	its.listSnapshot = newListSnapshot()

	prev := its.listSnapshot.head
	for _, ts := range ja.N {
		node := &orderedNode{
			precededType: jsonMap[ts.Hash()],
			prev:         prev,
			next:         nil,
		}
		prev.setNext(node)
		prev = node
	}
}

func (its *marshaledJSONType) unmarshalAsJSONPrimitive(doc *marshaledDocument, common *jsonCommon, parent jsonType) *jsonPrimitive {
	return &jsonPrimitive{
		common:  common,
		parent:  parent,
		deleted: its.D,
		K:       doc.unifyTimestamp(its.K),
		P:       doc.unifyTimestamp(its.P),
	}
}

func (its *marshaledJSONType) unmarshalAsJSONType(doc *marshaledDocument, common *jsonCommon, parent jsonType) jsonType {
	switch its.T {
	case mJSONElement:
		return &jsonElement{
			jsonType: its.unmarshalAsJSONPrimitive(doc, common, parent),
			V:        its.E,
		}
	case mJSONObject:
		return &jsonObject{
			jsonType: its.unmarshalAsJSONPrimitive(doc, common, parent),
		}
	case mJSONArray:
		return &jsonArray{
			jsonType: its.unmarshalAsJSONPrimitive(doc, common, parent),
		}
	}
	return nil
}
