package validator

import (
	"errors"
	"reflect"
)

func notNil(v reflect.Value, name, param string) error {
	switch v.Kind() {
	default:
		if v.IsNil() {
			return errors.New(name + " must not be nil")
		}
		return nil
	}
	return UnsupportedError(name)
}
