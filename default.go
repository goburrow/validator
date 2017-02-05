package validator

// DefaultOption returns an Option which sets validator to use
// tag name 'valid' and support function 'notempty', 'min', 'max'.
func DefaultOption() Option {
	return func(v *Validator) {
		v.tagName = defaultTagName
		v.register("notempty", notEmpty)
		v.register("min", min)
		v.register("max", max)
		v.register("regex", regex)
		v.register("date", date)
	}
}

// Default returns a new default validator.
// See DefaultOption for the options used for the validator.
func Default() *Validator {
	return New(DefaultOption())
}
