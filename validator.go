package regia

import "github.com/eatMoreApple/validate"

type Validator interface {
	Validate(v interface{}) error
}

type DefaultValidator struct{}

func (d DefaultValidator) Validate(v interface{}) error {
	return validate.Validate(v)
}

var defaultValidator Validator = DefaultValidator{}