// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

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
