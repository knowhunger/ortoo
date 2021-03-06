package utils

import (
	"encoding/json"
	"github.com/TylerBrock/colorjson"
	"github.com/knowhunger/ortoo/pkg/log"
)

func PrintMarshalDoc(l *log.OrtooLog, doc interface{}) {
	f := colorjson.NewFormatter()
	f.Indent = 2
	f.DisabledColor = true
	m, _ := json.Marshal(doc)
	var obj map[string]interface{}
	_ = json.Unmarshal(m, &obj)
	s, _ := f.Marshal(obj)
	l.Infof("%v", string(s))
}
