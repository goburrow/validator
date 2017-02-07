# Validator

Package validator implements value validation for struct fields.

The package was a fork of go-validator but was rewritten to:

- Reduce unneccesary memory allocations
- Support Slice/Array, Map, Interface,
- Simplify source code.

## Download
```
go get github.com/predixdeveloperACN/validator
```

## Example
```go
package main

import (
	"errors"
	"fmt"

	"github.com/predixdeveloperACN/validator"
)

type User struct {
	Name      string     `valid:"notempty"`
	Age       int        `valid:"min=13"`
	Addresses []*Address `valid:"min=1,max=2"`
	BirthDay  string     `valid:"date=MMDDYYYY"`
	Phone     int        `valid:"regex=^[0-9]*$"`
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
	// BirthDay is not a valid date,
	// Either address Line1 or Line2 must be set,
	// PostCode must not be less than 1 (was -1),
	// Country must have length not greater than 2 (was 3)
}
```

## Date Validator Placeholders
```
M    - month (1)
MM   - month (01)
MMM  - month (Jan)
MMMM - month (January)
D    - day (2)
DD   - day (02)
DDD  - day (Mon)
DDDD - day (Monday)
YY   - year (06)
YYYY - year (2006)
hh   - hours (15)
mm   - minutes (04)
ss   - seconds (05)

AM/PM hours: 'h' followed by optional 'mm' and 'ss' followed by 'pm', e.g.

hpm        - hours (03PM)
h:mmpm     - hours:minutes (03:04PM)
h:mm:sspm  - hours:minutes:seconds (03:04:05PM)

Time zones: a time format followed by 'ZZZZ', 'ZZZ' or 'ZZ', e.g.

hh:mm:ss ZZZZ (16:05:06 +0100)
hh:mm:ss ZZZ  (16:05:06 CET)
hh:mm:ss ZZ   (16:05:06 +01:00)
```