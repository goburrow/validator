package validator

import (
	"errors"
	"regexp"
	"reflect"
)

func regex(v reflect.Value, name, param string) error {

	// Resolve pointer
	for v.Kind() == reflect.Ptr {
		// Allow nil
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		return errors.New(name + " must not be of type slice, array or map.")
	default:
		s := v.String()
		re, err := regexp.Compile(param)
		if err != nil {
			return errors.New(name + " regex is not valid.")
		}
		if !re.MatchString(s) {
			return errors.New(name + " is not a valid value.")
		}

	}
	return nil
}