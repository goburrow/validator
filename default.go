package validator

// DefaultOption sets validator to use tag name 'valid'
// and support function 'notempty', 'min', 'max'.
var DefaultOption Option = func(v *Validator) {
	v.tagName = defaultTagName
	v.register("notempty", notEmpty)
	v.register("min", min)
	v.register("max", max)
}

// Default returns a new default validator.
// See DefaultOption for the options used for the validator.
func Default() *Validator {
	return New(DefaultOption)
}
