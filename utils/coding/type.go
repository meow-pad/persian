package coding

import "reflect"

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func IsFuncType(fnVal reflect.Value) bool {
	return fnVal.Kind() == reflect.Func
}

func IsErrorType(t reflect.Type) bool {
	return t == errorType || t.Implements(errorType)
}
