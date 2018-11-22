package tools

import (
	"reflect"
)

//return the i is nil or not
func IsNil(object interface{}) bool {
	if object == nil {
		return true
	}
	value := reflect.ValueOf(object)
	kind := value.Kind()
	value.Interface()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}
	return false
}
