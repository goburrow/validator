# Validator

[![GoDoc](https://godoc.org/github.com/goburrow/validator?status.svg)](https://godoc.org/github.com/goburrow/validator) [![Build Status](https://travis-ci.org/goburrow/validator.svg?branch=master)](https://travis-ci.org/goburrow/validator) 

Package validator implements value validation for struct fields.

The package was a fork of go-validator but was rewritten to:

- Reduce unneccesary memory allocations
- Support Slice/Array, Map, Interface,
- Simplify source code.
