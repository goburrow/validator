package validator

import "testing"

func mIntVal(v int) *int {
	return &v
}

func mUintVal(v uint) *uint {
	return &v
}

func TestMin(t *testing.T) {
	type s1 struct {
		A string         `valid:"min=3"`
		Z string         `valid:"min=1"`
		B uint           `valid:"min=100"`
		Y int            `valid:"min=-0xf"`
		C float64        `valid:"min=3.21"`
		X float32        `valid:"min=-1.1"`
		D *uint          `valid:"min=22222"`
		W *int           `valid:"min=077"`
		E []byte         `valid:"min=2"`
		V []byte         `valid:"min=1"`
		F map[int]uint   `valid:"min=5"`
		U map[string]int `valid:"min=2"`
	}
	s := s1{
		A: "ab",
		Z: "a",
		B: 99,
		Y: -10,
		C: 3.2099999,
		X: -1.0,
		D: mUintVal(22221),
		W: mIntVal(11112),
		E: []byte{1},
		V: []byte{0},
		F: map[int]uint{
			1: 2,
			3: 4,
			5: 6,
		},
		U: map[string]int{
			"a": 1,
			"b": 2,
		},
	}
	v := New(WithFunc("min", min))
	err := v.Validate(&s)
	if err == nil {
		t.Fatal("error expected")
	}
	errs := err.(Errors)
	if len(errs) != 6 {
		t.Fatalf("unexpected errors length: %+v; want: %d", err, 6)
	}
	if errs[0].Error() != "A must have length not less than 3 (was 2)" {
		t.Fatalf("unexpected error: %+v", errs[0])
	}
	if errs[1].Error() != "B must not be less than 100 (was 99)" {
		t.Fatalf("unexpected error: %+v", errs[1])
	}
	if errs[2].Error() != "C must not be less than 3.21 (was 3.2099999)" {
		t.Fatalf("unexpected error: %+v", errs[2])
	}
	if errs[3].Error() != "D must not be less than 22222 (was 22221)" {
		t.Fatalf("unexpected error: %+v", errs[3])
	}
	if errs[4].Error() != "E must have length not less than 2 (was 1)" {
		t.Fatalf("unexpected error: %+v", errs[4])
	}
	if errs[5].Error() != "F must have length not less than 5 (was 3)" {
		t.Fatalf("unexpected error: %+v", errs[5])
	}
}

func TestMinUnsupported(t *testing.T) {
	type s1 struct {
		A int
	}
	s := struct {
		A s1   `valid:"min=1"`
		B bool `valid:"min=0"`
	}{
		s1{0},
		false,
	}
	v := New(WithFunc("min", min))
	err := v.Validate(&s)
	if err == nil {
		t.Fatal("error expected")
	}
	errs := err.(Errors)
	if len(errs) != 2 {
		t.Fatalf("unexpected errors length: %+v; want: %d", errs, 2)
	}
	if errs[0].Error() != "validator: unsupported: A" {
		t.Fatalf("unexpected error: %+v", errs[0])
	}
	if errs[1].Error() != "validator: unsupported: B" {
		t.Fatalf("unexpected error: %+v", errs[1])
	}
}

func TestMax(t *testing.T) {
	type s1 struct {
		A string         `valid:"max=3"`
		Z string         `valid:"max=0"`
		B int            `valid:"max=-1"`
		Y uint           `valid:"max=99"`
		C float64        `valid:"max=-3.211984"`
		X float32        `valid:"max=19.54"`
		D *uint          `valid:"max=0xFF"`
		W *int           `valid:"max=19780601"`
		E []byte         `valid:"max=3"`
		V []byte         `valid:"max=1"`
		F map[int]uint   `valid:"max=2"`
		U map[string]int `valid:"max=2"`
	}
	s := s1{
		A: "abcd",
		Z: "",
		B: 0,
		Y: 99,
		C: -3.21,
		X: 19.53,
		D: mUintVal(0x100),
		W: mIntVal(19780601),
		E: []byte{1, 2, 3, 4},
		V: nil,
		F: map[int]uint{
			1: 2,
			3: 4,
			5: 6,
		},
		U: map[string]int{
			"a": 1,
			"b": 2,
		},
	}
	v := New(WithFunc("max", max))
	err := v.Validate(&s)
	if err == nil {
		t.Fatal("error expected")
	}
	errs := err.(Errors)
	if len(errs) != 6 {
		t.Fatalf("unexpected errors length: %+v; want: %d", err, 6)
	}
	if errs[0].Error() != "A must have length not greater than 3 (was 4)" {
		t.Fatalf("unexpected error: %+v", errs[0])
	}
	if errs[1].Error() != "B must not be greater than -1 (was 0)" {
		t.Fatalf("unexpected error: %+v", errs[1])
	}
	if errs[2].Error() != "C must not be greater than -3.211984 (was -3.21)" {
		t.Fatalf("unexpected error: %+v", errs[2])
	}
	if errs[3].Error() != "D must not be greater than 0xFF (was 256)" {
		t.Fatalf("unexpected error: %+v", errs[3])
	}
	if errs[4].Error() != "E must have length not greater than 3 (was 4)" {
		t.Fatalf("unexpected error: %+v", errs[4])
	}
	if errs[5].Error() != "F must have length not greater than 2 (was 3)" {
		t.Fatalf("unexpected error: %+v", errs[5])
	}
}

func TestMaxUnsupported(t *testing.T) {
	type s1 struct {
		A int
	}
	s := struct {
		A s1   `valid:"max=-1"`
		B bool `valid:"max=99"`
	}{
		s1{0},
		false,
	}
	v := New(WithFunc("max", max))
	err := v.Validate(&s)
	if err == nil {
		t.Fatal("error expected")
	}
	errs := err.(Errors)
	if len(errs) != 2 {
		t.Fatalf("unexpected errors length: %+v; want: %d", errs, 2)
	}
	if errs[0].Error() != "validator: unsupported: A" {
		t.Fatalf("unexpected error: %+v", errs[0])
	}
	if errs[1].Error() != "validator: unsupported: B" {
		t.Fatalf("unexpected error: %+v", errs[1])
	}
}
