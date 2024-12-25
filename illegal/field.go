//go:build !solution

package illegal

import (
	"reflect"
	"unsafe"
)

func SetPrivateField(obj interface{}, name string, value interface{}) {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return
	}
	field := val.Elem().FieldByName(name)
	pointer := unsafe.Pointer(field.UnsafeAddr())
	newVal := reflect.ValueOf(value)
	reflect.NewAt(field.Type(), pointer).Elem().Set(newVal)
}
