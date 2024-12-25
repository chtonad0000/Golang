//go:build !solution

package reversemap

import "reflect"

func ReverseMap(forward interface{}) interface{} {
	v := reflect.ValueOf(forward)

	if v.Kind() == reflect.Map {
		keyType := v.Type().Key()
		valueType := v.Type().Elem()
		result := reflect.MakeMap(reflect.MapOf(valueType, keyType))
		for _, k := range v.MapKeys() {
			result.SetMapIndex(v.MapIndex(k), k)
		}

		return result.Interface()
	}

	panic("Not a map")
}
