//go:build !solution

package structtags

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var cache sync.Map

func fieldMapping(t reflect.Type) map[string]int {
	if mapping, ok := cache.Load(t); ok {
		return mapping.(map[string]int)
	}

	fieldMap := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		fieldInfo := t.Field(i)
		tag := fieldInfo.Tag.Get("http")
		if tag == "" {
			tag = strings.ToLower(fieldInfo.Name)
		}
		fieldMap[tag] = i
	}
	cache.Store(t, fieldMap)
	return fieldMap
}

func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	v := reflect.ValueOf(ptr).Elem()
	typ := v.Type()
	fieldMapping := fieldMapping(typ)

	for name, values := range req.Form {
		fieldIndex, ok := fieldMapping[name]
		if !ok {
			continue
		}

		field := v.Field(fieldIndex)
		for _, value := range values {
			if field.Kind() == reflect.Slice {
				elem := reflect.New(field.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				field.Set(reflect.Append(field, elem))
			} else {
				if err := populate(field, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		v.SetInt(int64(i))
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)
	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}
