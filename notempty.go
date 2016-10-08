package validator

import (
	"errors"
	"reflect"
)

func notEmpty(v reflect.Value, name, param string) error {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		if v.Len() == 0 {
			return errors.New(name + " must not be empty")
		}
		return nil
	case reflect.Bool:
		if !v.Bool() {
			return errors.New(name + " must not be false")
		}
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() == 0 {
			return errors.New(name + " must not be zero")
		}
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if v.Uint() == 0 {
			return errors.New(name + " must not be zero")
		}
		return nil
	case reflect.Float32, reflect.Float64:
		if v.Float() == 0 {
			return errors.New(name + " must not be zero")
		}
		return nil
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return errors.New(name + " must not be nil")
		}
		return nil
	}
	return UnsupportedError(name)
}
