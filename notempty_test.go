package validator

import "testing"

func TestNotEmpty(t *testing.T) {
	type s1 struct {
		A string      `valid:"notempty"`
		Z string      `valid:"notempty"`
		B int         `valid:"notempty"`
		Y int         `valid:"notempty"`
		C bool        `valid:"notempty"`
		X bool        `valid:"notempty"`
		D *testing.T  `valid:"notempty"`
		W *testing.T  `valid:"notempty"`
		E []byte      `valid:"notempty"`
		V []byte      `valid:"notempty"`
		F interface{} `valid:"notempty"`
		U interface{} `valid:"notempty"`
	}
	s := s1{
		A: "",
		Z: "a",
		B: 0,
		Y: 1,
		C: false,
		X: true,
		D: nil,
		W: t,
		E: []byte{},
		V: []byte("a"),
		F: nil,
		U: t,
	}
	v := New(WithFunc("notempty", notEmpty))
	err := v.Validate(&s)
	if err == nil {
		t.Fatal("error expected")
	}
	errs := err.(Errors)
	if len(errs) != 6 {
		t.Fatalf("unexpected errors length: %+v; want: %d", err, 6)
	}
	if errs[0].Error() != "A must not be empty" {
		t.Fatalf("unexpected error: %+v", errs[0])
	}
	if errs[1].Error() != "B must not be zero" {
		t.Fatalf("unexpected error: %+v", errs[1])
	}
	if errs[2].Error() != "C must not be false" {
		t.Fatalf("unexpected error: %+v", errs[2])
	}
	if errs[3].Error() != "D must not be nil" {
		t.Fatalf("unexpected error: %+v", errs[3])
	}
	if errs[4].Error() != "E must not be empty" {
		t.Fatalf("unexpected error: %+v", errs[4])
	}
	if errs[5].Error() != "F must not be nil" {
		t.Fatalf("unexpected error: %+v", errs[5])
	}
}

func TestNotEmptyUnsupported(t *testing.T) {
	type s1 struct {
		A int
	}
	s := struct {
		A s1 `valid:"notempty"`
	}{
		s1{0},
	}
	v := New(WithFunc("notempty", notEmpty))
	err := v.Validate(&s)
	if err == nil {
		t.Fatal("error expected")
	}
	errs := err.(Errors)
	if len(errs) != 1 {
		t.Fatalf("unexpected errors length: %+v; want: %d", errs, 1)
	}
	if errs[0].Error() != "validator: unsupported: A" {
		t.Fatalf("unexpected error: %+v", errs[0])
	}
}
