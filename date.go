package validator

import (
	"errors"
	"reflect"
	"github.com/metakeule/fmtdate"
)

func date(v reflect.Value, name, param string) error {

	// Resolve pointer
	for v.Kind() == reflect.Ptr {
		// Allow nil
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String:
		_, err := fmtdate.Parse(param, v.String())
		if err != nil {
			return errors.New(name + " is not a valid date. " + err.Error())
		}
		return nil
	}
	return errors.New(name + " is not a valid date.")
}