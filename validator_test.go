package validator

import (
	"errors"
	"reflect"
	"testing"
)

var errNOK = errors.New("nok")

func ok(rv reflect.Value, name, param string) error {
	return nil
}

func nok(rv reflect.Value, name, param string) error {
	if param != "" {
		return errors.New(param)
	}
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
		B int    `valid:"nok=B,ok=ok"`
		C bool
		D float64 `valid:""`
		E *string `valid:"ok,nok=E"`
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
	assertNOK(t, err, "B", "E")
	err = v.Validate(&s)
	assertNOK(t, err, "B", "E")
}

func TestStructOfStruct(t *testing.T) {
	v := newTestValidator()

	type s1 struct {
		A int `valid:"nok=A"`
		B *s1 `valid:"ok"`
	}
	type s2 struct {
		C string `valid:"nok=C"`
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
	assertNOK(t, err, "C", "A", "A", "A")
	err = v.Validate(&s)
	assertNOK(t, err, "C", "A", "A", "A")
}

func TestPrimitiveSlice(t *testing.T) {
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
		A bool `valid:"nok=A"`
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
	assertNOK(t, err, "A", "A", "A")
	err = v.Validate(&s)
	assertNOK(t, err, "A", "A", "A")
}

func TestSliceOfStructPointer(t *testing.T) {
	v := newTestValidator()

	type s1 struct {
		A uint `valid:"nok=A"`
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
	assertNOK(t, err, "A", "A")
	err = v.Validate(&s)
	assertNOK(t, err, "A", "A")
}

func TestStructOfSliceOfStruct(t *testing.T) {
	v := newTestValidator()

	type s1 struct {
		A int `valid:"nok=A"`
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
	assertNOK(t, err, "A", "A", "A", "A", "A")
	err = v.Validate(&s)
	assertNOK(t, err, "A", "A", "A", "A", "A")
}

func TestStructOfInterface(t *testing.T) {
	v := newTestValidator()

	type s1 struct {
		A int `valid:"nok=A"`
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
	assertNOK(t, err, "A", "A", "A", "A", "A")
}

func TestMap(t *testing.T) {
	type s1 struct {
		A int `valid:"nok=A"`
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
	assertNOK(t, err, "A", "A")
}

func TestMapOfInterface(t *testing.T) {
	type s1 struct {
		A int `valid:"nok=A"`
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
	assertNOK(t, err, "A", "A", "A")
}

func TestEmbedded(t *testing.T) {
	type s1 struct {
		A int  `valid:"nok=eA"`
		B bool `valid:"nok=eB"`
	}
	type s2 struct {
		s1
		S1 struct {
			C string `valid:"nok=C"`
			D uint   `valid:"ok"`
			s1
		}
		E string `valid:"nok=E"`
		B bool   `valid:"nok=B"`
	}
	s := s2{}
	s.A = 0
	s.B = true
	s.E = "e"
	s.S1.C = "c"
	s.S1.D = 1
	v := newTestValidator()
	err := v.Validate(&s)
	assertNOK(t, err, "eA", "eB", "C", "eA", "eB", "E", "B")
}

func TestUnexportField(t *testing.T) {
	type s1 struct {
		A int `valid:"nok=A"`
	}
	type s2 struct {
		s1 `valid:"ok=s1"`
		b  int `valid:"nok=B"`
	}
	s := s2{
		s1: s1{0},
		b:  1,
	}
	v := newTestValidator()
	err := v.Validate(&s)
	assertNOK(t, err, "A")
}

func assertNOK(t *testing.T, err error, msg ...string) {
	if err == nil {
		t.Fatal("error expected")
	}
	errs := err.(Errors)
	if len(errs) != len(msg) {
		t.Fatalf("unexpected errors: %v; want: %v", errs, msg)
	}
	for i, e := range errs {
		if e.Error() != msg[i] {
			t.Fatalf("unexpected error: %+v; want: %v", e, msg)
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
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			v.Validate(&s)
		}
	})
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
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			v.Validate(&s)
		}
	})
}

func BenchmarkWithTag(b *testing.B) {
	b.ReportAllocs()
	v := newTestValidator()

	type s1 struct {
		A int  `valid:"nok"`
		B uint `valid:"nok"`
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
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			v.Validate(&s)
		}
	})
}
