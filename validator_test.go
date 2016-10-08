package validator

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

var errNOK = errors.New("nok")

func ok(rv reflect.Value, name, param string) error {
	return nil
}

func nok(rv reflect.Value, name, param string) error {
	return errNOK
}

func strPtr(v string) *string {
	return &v
}

func intPtr(v int) *int {
	return &v
}

func newTestValidator() *Validator {
	return New(WithFunc("ok", ok), WithFunc("nok", nok))
}

func TestSimpleStruct(t *testing.T) {
	v := newTestValidator()

	s := struct {
		A string `valid:"ok"`
		B int    `valid:"nok=b,ok=ok"`
		C bool
		D float64 `valid:""`
		E *string `valid:"ok=c,nok"`
		F *int    `valid:"-"`
	}{
		"a",
		1, // err
		true,
		1.0,
		strPtr("e"), // err
		intPtr(2),
	}
	err := v.Validate(s)
	assertNOK(t, err, 2)
	err = v.Validate(&s)
	assertNOK(t, err, 2)
}

func TestStructOfStruct(t *testing.T) {
	v := newTestValidator()

	type s1 struct {
		A int `valid:"nok"`
		B *s1 `valid:"ok"`
	}
	type s2 struct {
		C string `valid:"nok"`
		D s1
		E *s1
		F *s1
	}

	s := s2{
		C: "c", // err
		D: s1{
			A: 1, // err
			B: &s1{
				A: 2, // err
				B: nil,
			},
		},
		E: nil,
		F: &s1{
			A: 3, // err
			B: nil,
		},
	}
	err := v.Validate(s)
	assertNOK(t, err, 4)
	err = v.Validate(&s)
	assertNOK(t, err, 4)
}

func TestSlice(t *testing.T) {
	v := newTestValidator()
	err := v.Validate([]string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	err = v.Validate([]byte("abc"))
	if err != nil {
		t.Fatal(err)
	}
	type s1 struct {
		A bool `valid:"ok"`
	}
	s := []s1{
		s1{true},
		s1{false},
	}
	err = v.Validate(s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSliceOfStruct(t *testing.T) {
	v := newTestValidator()

	type s1 struct {
		A bool `valid:"nok=2"`
	}
	s := []s1{}
	err := v.Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	err = v.Validate(&s)
	if err != nil {
		t.Fatal(err)
	}

	s = []s1{
		s1{false},
		s1{true},
		s1{false},
	}
	err = v.Validate(s)
	assertNOK(t, err, 3)
	err = v.Validate(&s)
	assertNOK(t, err, 3)
}

func TestSliceOfStructPointer(t *testing.T) {
	v := newTestValidator()

	type s1 struct {
		A uint `valid:"nok"`
	}
	s := []*s1{}
	err := v.Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	err = v.Validate(&s)
	if err != nil {
		t.Fatal(err)
	}

	s = []*s1{
		&s1{1},
		&s1{0},
	}
	err = v.Validate(s)
	assertNOK(t, err, 2)
	err = v.Validate(&s)
	assertNOK(t, err, 2)
}

func TestStructOfSliceOfStruct(t *testing.T) {
	v := newTestValidator()

	type s1 struct {
		A int `valid:"nok"`
	}
	type s2 struct {
		C []s1
		D []*s1
		E *[]s1
		F *[]*s1
	}
	s := s2{
		C: []s1{s1{1}},
		D: []*s1{&s1{2}, &s1{3}},
		E: &[]s1{s1{4}},
		F: &[]*s1{&s1{5}},
	}
	err := v.Validate(s)
	assertNOK(t, err, 5)
	err = v.Validate(&s)
	assertNOK(t, err, 5)
}

func TestStructOfInterface(t *testing.T) {
	v := newTestValidator()

	type s1 struct {
		A int `valid:"nok"`
	}
	type s2 struct {
		B interface{}
		C interface{}
		D interface{}
		E interface{}
		F interface{}
	}
	s := s2{
		B: s1{1},
		C: &s1{2},
		D: nil,
		E: []s1{s1{3}},
		F: []interface{}{
			&s2{
				B: &s1{4},
			},
			&s1{5},
		},
	}
	err := v.Validate(&s)
	assertNOK(t, err, 5)
}

func TestMap(t *testing.T) {
	type s1 struct {
		A int `valid:"nok"`
	}
	v := newTestValidator()
	m1 := map[int]string{
		1: "a",
		2: "b",
	}
	err := v.Validate(m1)
	if err != nil {
		t.Fatal(err)
	}
	m2 := map[string]*s1{
		"a": &s1{1},
		"b": &s1{2},
	}
	err = v.Validate(m2)
	assertNOK(t, err, 2)
}

func TestMapOfInterface(t *testing.T) {
	type s1 struct {
		A int `valid:"nok"`
	}
	v := newTestValidator()
	m := map[int]interface{}{
		1: &s1{1},
		2: map[string]string{
			"0": "a",
			"1": "b",
		},
		3: map[int]interface{}{
			3: &s1{2},
			4: s1{3},
		},
	}
	err := v.Validate(m)
	assertNOK(t, err, 3)
}

func assertNOK(t *testing.T, err error, n int) {
	if err == nil {
		t.Fatal("error expected")
	}
	errs := err.(Errors)
	if len(errs) != n {
		t.Fatalf("unexpected errors: %+v; want: %d", errs, n)
	}
	for _, e := range errs {
		if !strings.Contains(e.Error(), "nok") {
			t.Fatalf("unexpected error: %+v", e)
		}
	}
}

func BenchmarkSimple(b *testing.B) {
	b.ReportAllocs()
	v := newTestValidator()

	s := struct {
		A string `valid:"nok"`
		B int    `valid:"ok"`
	}{
		"a",
		1,
	}
	for i := 0; i < b.N; i++ {
		err := v.Validate(&s)
		if err == nil {
			b.Fatalf("error expected")
		}
	}
}

func BenchmarkNoTag(b *testing.B) {
	b.ReportAllocs()
	v := newTestValidator()

	type s1 struct {
		A int
		B uint
	}
	type s2 struct {
		B string
		C *s1
		D s1
		E []s1
	}

	s := s2{
		B: "nok",
		C: &s1{1, 0},
		D: s1{2, 0},
		E: []s1{
			s1{3, 0},
			s1{4, 0},
		},
	}
	for i := 0; i < b.N; i++ {
		err := v.Validate(&s)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWithTag(b *testing.B) {
	b.ReportAllocs()
	v := newTestValidator()

	type s1 struct {
		A int  `valid:"nok=1"`
		B uint `valid:"nok=2"`
	}
	type s2 struct {
		B string `valid:"ok"`
		C *s1
		D s1
		E []s1
	}

	s := s2{
		B: "nok",
		C: &s1{1, 0},
		D: s1{2, 0},
		E: []s1{
			s1{3, 0},
			s1{4, 0},
		},
	}
	for i := 0; i < b.N; i++ {
		err := v.Validate(&s)
		if err == nil {
			b.Fatal("error expected")
		}
	}
}
