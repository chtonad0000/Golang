//go:build !solution

package jsonlist

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

func Marshal(w io.Writer, slice interface{}) error {
	val := reflect.ValueOf(slice)

	if val.Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}

	encoder := json.NewEncoder(w)
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i).Interface()
		if err := encoder.Encode(elem); err != nil {
			return err
		}
	}

	return nil
}
func Unmarshal(r io.Reader, slice interface{}) error {

	sliceValue := reflect.ValueOf(slice)

	if sliceValue.Kind() != reflect.Ptr || sliceValue.Elem().Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}

	sliceValue = sliceValue.Elem()
	decoder := json.NewDecoder(r)

	for {
		var item interface{}
		if err := decoder.Decode(&item); err != nil {
			if err.Error() == "EOF" {
				return nil
			}
			return err
		}

		itemValue := reflect.ValueOf(item)
		elemType := sliceValue.Type().Elem()
		fmt.Println(itemValue)
		if itemValue.Kind() == reflect.Map && elemType.Kind() == reflect.Struct {
			elemValue := reflect.New(elemType).Elem()
			for _, key := range itemValue.MapKeys() {
				mapValue := itemValue.MapIndex(key)
				if !mapValue.IsValid() || mapValue.IsNil() {
					continue
				}
				mapValue = reflect.ValueOf(itemValue.MapIndex(key).Interface())
				fieldValue := elemValue.FieldByName(key.String())
				if fieldValue.IsValid() && fieldValue.CanSet() {
					if mapValue.Type() != fieldValue.Type() {
						mapValue = mapValue.Convert(fieldValue.Type())
					}
				}
				fieldValue.Set(mapValue)
			}
			sliceValue.Set(reflect.Append(sliceValue, elemValue))
		} else {
			if itemValue.Type() != elemType {
				if !itemValue.Type().ConvertibleTo(elemType) {
					return fmt.Errorf("cannot convert value of type %s to %s", itemValue.Type(), elemType)
				}
				itemValue = itemValue.Convert(elemType)
			}

			sliceValue.Set(reflect.Append(sliceValue, itemValue))
		}
	}
}
