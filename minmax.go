package validator

import (
	"fmt"
	"reflect"
	"strconv"
)

func min(v reflect.Value, name, param string) error {
	// Resolve pointer
	for v.Kind() == reflect.Ptr {
		// Allow nil
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		if v.Len() < int(parseInt(param)) {
			return fmt.Errorf("%s must have length not less than %s (was %v)", name, param, v.Len())
		}
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() < parseInt(param) {
			return fmt.Errorf("%s must not be less than %s (was %v)", name, param, v.Int())
		}
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if v.Uint() < parseUint(param) {
			return fmt.Errorf("%s must not be less than %s (was %v)", name, param, v.Uint())
		}
		return nil
	case reflect.Float32, reflect.Float64:
		if v.Float() < parseFloat(param) {
			return fmt.Errorf("%s must not be less than %s (was %v)", name, param, v.Float())
		}
		return nil
	}
	return UnsupportedError(name)
}

func max(v reflect.Value, name, param string) error {
	// Resolve pointer
	for v.Kind() == reflect.Ptr {
		// Allow nil
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		if v.Len() > int(parseInt(param)) {
			return fmt.Errorf("%s must have length not greater than %s (was %v)", name, param, v.Len())
		}
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() > parseInt(param) {
			return fmt.Errorf("%s must not be greater than %s (was %v)", name, param, v.Int())
		}
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if v.Uint() > parseUint(param) {
			return fmt.Errorf("%s must not be greater than %s (was %v)", name, param, v.Uint())
		}
		return nil
	case reflect.Float32, reflect.Float64:
		if v.Float() > parseFloat(param) {
			return fmt.Errorf("%s must not be greater than %s (was %v)", name, param, v.Float())
		}
		return nil
	}
	return UnsupportedError(name)
}

var defaultIntCache intCache

func parseInt(s string) int64 {
	val, ok := defaultIntCache.get(s)
	if ok {
		return val
	}
	val, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		panic(err)
	}
	defaultIntCache.save(s, val)
	return val
}

var defaultUintCache uintCache

func parseUint(s string) uint64 {
	val, ok := defaultUintCache.get(s)
	if ok {
		return val
	}
	val, err := strconv.ParseUint(s, 0, 0)
	if err != nil {
		panic(err)
	}
	defaultUintCache.save(s, val)
	return val
}

var defaultFloatCache floatCache

func parseFloat(s string) float64 {
	val, ok := defaultFloatCache.get(s)
	if ok {
		return val
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	defaultFloatCache.save(s, val)
	return val
}
