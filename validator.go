// Package validator provides validation for structs and fields.
package validator

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

const (
	defaultTagName = "valid"
)

// UnsupportedError is a generic error returned when validation function
// is applied on fields it does not support or the function has not been
// registered to Validator.
type UnsupportedError string

func (e UnsupportedError) Error() string {
	return "validator: unsupported: " + string(e)
}

// Errors is a list of error.
type Errors []error

// Error returns formatted string of all underlying errors.
func (e Errors) Error() string {
	var buf bytes.Buffer
	for i, err := range e {
		if i > 0 {
			buf.WriteString("; ")
		}
		buf.WriteString(err.Error())
	}
	return buf.String()
}

// Func validates field with value v, field name and parameter p.
type Func func(v reflect.Value, name, param string) error

// Option sets options for the validator.
type Option func(v *Validator)

// Validatable is an interface implemented by types that can
// validate themselves.
type Validatable interface {
	Validate() error
}

var validatableType = reflect.TypeOf(new(Validatable)).Elem()

// Validator implements value validation for structs and fields.
type Validator struct {
	tagName string
	funcs   map[string]Func

	fieldCache fieldCache
}

// New allocates and returns a new Validator with given options.
// To create a new Validator with default options, use Default instead.
func New(options ...Option) *Validator {
	v := &Validator{
		tagName: defaultTagName,
		funcs:   make(map[string]Func),
	}
	for _, opt := range options {
		opt(v)
	}
	return v
}

// WithTagName returns an Option which sets tagName to the validator.
func WithTagName(tagName string) Option {
	return func(v *Validator) {
		v.tagName = tagName
	}
}

// WithFunc returns an Option which adds a new function handler.
// If a function with same name existed, it will be overriden by the given one.
// It panics if name is empty or handler is nil.
func WithFunc(name string, fn Func) Option {
	if name == "" || fn == nil {
		panic("validator: invalid handler " + name)
	}
	return func(v *Validator) {
		v.register(name, fn)
	}
}

func (a *Validator) register(name string, fn Func) {
	a.funcs[name] = fn
}

// Validate validates given value. Value v is usually a pointer to
// the struct to validate, but it can also be a struct, slice or array.
func (a *Validator) Validate(v interface{}) (err error) {
	s := state{validator: a}
	defer func() {
		if r := recover(); r != nil {
			if rerr, ok := r.(error); ok {
				s.addError(rerr)
			} else {
				s.addError(fmt.Errorf("%v", r))
			}
			err = Errors(s.errors)
		}
	}()
	s.validateInterface(v)
	if len(s.errors) == 0 {
		return nil
	}
	return Errors(s.errors)
}

func (a *Validator) getFields(rt reflect.Type) []field {
	fields, ok := a.fieldCache.get(rt)
	if ok {
		return fields
	}
	fields = make([]field, 0, 10)

	n := rt.NumField()
	for i := 0; i < n; i++ {
		ft := rt.Field(i)
		if !ft.Anonymous {
			// Ignore unexported but allow embedded fields
			if unicode.IsLower([]rune(ft.Name)[0]) {
				continue
			}
		}
		// Explicitly ignored
		tags := ft.Tag.Get(a.tagName)
		if tags == "-" {
			continue
		}
		if tags == "" {
			if !supported(ft.Type) {
				continue
			}
		}
		fields = append(fields, field{
			idx:  i,
			name: ft.Name,
			tags: tags,
		})
	}
	a.fieldCache.save(rt, fields)
	return fields
}

func supported(rt reflect.Type) bool {
	switch rt.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map, reflect.Ptr, reflect.Interface:
		return true
	}
	return rt.Implements(validatableType)
}

type state struct {
	validator *Validator

	errors []error
}

func (s *state) validateValue(rv reflect.Value) {
	// Call Validate method if this value implements Validatable
	s.validateValidatable(rv)

	// Resolve pointer
	for rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Struct:
		s.validateStruct(rv)
	case reflect.Slice, reflect.Array:
		s.validateSlice(rv)
	case reflect.Map:
		s.validateMap(rv)
	case reflect.Interface:
		s.validateInterface(rv.Interface())
	}
}

func (s *state) validateInterface(v interface{}) {
	rv := reflect.ValueOf(v)
	s.validateValue(rv)
}

func (s *state) validateStruct(rv reflect.Value) {
	rt := rv.Type()

	fields := s.validator.getFields(rt)
	n := len(fields)
	for i := 0; i < n; i++ {
		ft := &fields[i]
		fv := rv.Field(ft.idx)
		if ft.tags != "" {
			// Validate this field
			s.validateField(fv, ft.name, ft.tags)
		}
		s.validateValue(fv)
	}
}

func (s *state) validateSlice(rv reflect.Value) {
	rt := rv.Type()
	if !supported(rt.Elem()) {
		return
	}
	n := rv.Len()
	for i := 0; i < n; i++ {
		fv := rv.Index(i)
		s.validateValue(fv)
	}
}

func (s *state) validateMap(rv reflect.Value) {
	rt := rv.Type()
	if !supported(rt.Elem()) {
		return
	}
	if rv.Len() == 0 {
		return
	}
	for _, k := range rv.MapKeys() {
		fv := rv.MapIndex(k)
		s.validateValue(fv)
	}
}

func (s *state) validateField(fv reflect.Value, name, tags string) {
	for tags != "" {
		var tag string
		i := strings.Index(tags, ",")
		if i < 0 {
			tag = tags
			tags = ""
		} else {
			tag = tags[:i]
			tags = tags[i+1:]
		}
		fn, param := parseTag(tag)
		f, ok := s.validator.funcs[fn]
		if ok {
			err := f(fv, name, param)
			if err != nil {
				s.addError(err)
			}
		} else {
			s.addError(UnsupportedError(fn))
		}
	}
}

func (s *state) validateValidatable(rv reflect.Value) {
	if !rv.IsValid() || rv.Type().NumMethod() == 0 {
		return
	}
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return
	}
	if f, ok := rv.Interface().(Validatable); ok {
		err := f.Validate()
		if err != nil {
			s.addError(err)
		}
	}
}

func (s *state) addError(err error) {
	s.errors = append(s.errors, err)
}

// parseTag returns function name and parameter.
func parseTag(tag string) (name, param string) {
	i := strings.Index(tag, "=")
	if i < 0 {
		name = tag
	} else {
		name = tag[:i]
		param = tag[i+1:]
	}
	return
}
