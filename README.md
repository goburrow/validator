# Validator

[![GoDoc](https://godoc.org/github.com/goburrow/validator?status.svg)](https://godoc.org/github.com/goburrow/validator) [![Build Status](https://travis-ci.org/goburrow/validator.svg?branch=master)](https://travis-ci.org/goburrow/validator) 

Package validator implements value validation for struct fields.

The package was a fork of go-validator but was rewritten to:

- Reduce unneccesary memory allocations
- Support Slice/Array, Map, Interface,
- Simplify source code.

## Download
```
go get -u github.com/goburrow/validator
```

## Example
```go
package main

import (
	"errors"
	"fmt"

	"github.com/goburrow/validator"
)

type User struct {
	Name      string     `valid:"notempty"`
	Age       int        `valid:"min=13"`
	Addresses []*Address `valid:"min=1,max=2"`
}

type Address struct {
	Line1    string
	Line2    string
	PostCode int    `valid:"min=1"`
	Country  string `valid:"notempty,max=2"`
}

func (a *Address) Validate() error {
	if a.Line1 == "" && a.Line2 == "" {
		return errors.New("Either address Line1 or Line2 must be set")
	}
	return nil
}

func main() {
	u := &User{
		Addresses: []*Address{
			&Address{
				Line1:    "Somewhere",
				PostCode: 1000,
				Country:  "AU",
			},
			&Address{
				PostCode: -1,
				Country:  "US",
			},
			&Address{
				Line2:    "Here",
				PostCode: 1,
				Country:  "USA",
			},
		},
	}
	v := validator.Default()
	fmt.Println(v.Validate(u))
	// Output:
	// Name must not be empty,
	// Age must not be less than 13 (was 0),
	// Addresses must have length not greater than 2 (was 3),
	// Either address Line1 or Line2 must be set,
	// PostCode must not be less than 1 (was -1),
	// Country must have length not greater than 2 (was 3)
}
```
