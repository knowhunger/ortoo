package types

import (
	"fmt"
	"github.com/knowhunger/ortoo/ortoo/log"
	"math"
)

// JSONValue is an internal type used in storing various types, for converting any type to JSON supported type.
type JSONValue interface{}

func ConvertValueList(values []interface{}) ([]interface{}, error) {
	var jsonValues []interface{}
	for _, val := range values {
		if val == nil {
			return nil, fmt.Errorf("null value cannot be inserted")
		}
		jsonValues = append(jsonValues, ConvertToJSONSupportedValue(val))
	}
	return jsonValues, nil
}

// ConvertToJSONSupportedValue converts any type of Go into a type that is supported by JSON
func ConvertToJSONSupportedValue(t interface{}) JSONValue {
	switch v := t.(type) {
	// all number types are stored as float64, i.e., IEEE 754 64 bits floating point type.
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
		*int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64:
		var i64 int64
		switch vv := v.(type) {
		case int:
			i64 = int64(vv)
		case int8:
			i64 = int64(vv)
		case int16:
			i64 = int64(vv)
		case int32:
			i64 = int64(vv)
		case int64:
			i64 = vv
		case uint:
			i64 = int64(vv)
		case uint8:
			i64 = int64(vv)
		case uint16:
			i64 = int64(vv)
		case uint32:
			i64 = int64(vv)
		case uint64:
			if vv > math.MaxInt64 {
				log.Logger.Warnf("overflow: cannot store an integer more than int64.Max (%d)", math.MaxInt64)
			}
			i64 = int64(vv)
		case *int:
			i64 = int64(*vv)
		case *int8:
			i64 = int64(*vv)
		case *int16:
			i64 = int64(*vv)
		case *int32:
			i64 = int64(*vv)
		case *int64:
			i64 = *vv
		case *uint:
			i64 = int64(*vv)
		case *uint8:
			i64 = int64(*vv)
		case *uint16:
			i64 = int64(*vv)
		case *uint32:
			i64 = int64(*vv)
		case *uint64:
			if *vv > math.MaxInt64 {
				log.Logger.Warnf("overflow: cannot store an integer more than int64.Max (%d)", math.MaxInt64)
			}
			i64 = int64(*vv)
		}
		return i64
	case float32, float64, *float32, *float64:
		var f64 float64
		switch vv := v.(type) {
		case float32:
			f64 = float64(vv)
		case float64:
			f64 = vv
		case *float32:
			f64 = float64(*vv)
		case *float64:
			f64 = *vv
		}
		return f64
	case string:
		return v
	case *string:
		return *v
	case bool:
		return v
	case *bool:
		return *v
	default:
	}
	return t
}