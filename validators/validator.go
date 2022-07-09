// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package validators

import (
	"regexp"
)

var validateParamRegexp = regexp.MustCompile(`(?P<validator>\w+)\((?P<params>.*)\)`)

type Param struct {
	Message string
	Value   string
}

func Validate(value interface{}) error {
	return validatorLibrary.Validate(value)
}

type DefaultValidator struct{}

func (d DefaultValidator) Validate(v interface{}) error {
	return Validate(v)
}

var _ Validator = DefaultValidator{}
