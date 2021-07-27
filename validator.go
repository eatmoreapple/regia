package regia

import "github.com/eatmoreapple/validate"

type Validator interface {
	Validate(v interface{}) error
}

type DefaultValidator struct{}

func (d DefaultValidator) Validate(v interface{}) error {
	return validate.Validate(v)
}

var _ Validator = DefaultValidator{}
