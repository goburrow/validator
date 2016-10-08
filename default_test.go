package validator

import "fmt"

func ExampleValidator() {
	type data struct {
		A string       `valid:"notempty"`
		B int          `valid:"min=1,max=10"`
		C map[int]bool `valid:"notempty,max=2"`
		D []data
	}
	d := &data{
		A: "",
		B: 11,
		C: map[int]bool{
			1: false,
			2: true,
			3: false,
		},
		D: []data{
			data{
				A: "ab",
				B: 2,
			},
			data{
				A: "cd",
				B: 0,
				C: map[int]bool{
					0: false,
				},
			},
		},
	}
	v := Default()
	err := v.Validate(d)
	if err != nil {
		for _, e := range err.(Errors) {
			fmt.Println(e)
		}
	}
	// Output:
	// A must not be empty
	// B must not be greater than 10 (was 11)
	// C must have length not greater than 2 (was 3)
	// C must not be empty
	// B must not be less than 1 (was 0)
}
