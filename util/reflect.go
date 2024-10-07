package util

import (
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"strings"
)

func DeRef(v any) any {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		return val.Elem().Interface()
	}
	return v
}

func MapMySQL(v any, colType string) (result string) {
	v = DeRef(v)
	switch tp := v.(type) {
	case uint:
	case uint8:
	case uint16:
	case uint32:
	case uint64:
	case int:
	case int8:
	case int16:
	case int32:
	case int64:
		result = fmt.Sprintf("%d", tp)
	case float32:
	case float64:
		result = fmt.Sprintf("%f", tp)
	case []byte:
		if len(tp) == 1 && strings.Contains(colType, "bit") {
			result = fmt.Sprintf("b'%d'", tp[0])
		} else if strings.Contains(colType, "blob") {
			result = "0x" + hex.EncodeToString(tp)
		} else {
			result = fmt.Sprintf("'%s'", string(tp))
		}
	case string:
		result = fmt.Sprintf("'%s'", tp)
	case nil:
		result = "NULL"
	default:
		log.Fatalf("table_field_value map error:%s", tp)
	}
	return result
}
