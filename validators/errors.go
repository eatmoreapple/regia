// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package validators

import "errors"

type ValidationError struct {
	error
	FieldName string `json:"field_name"`
	Rule      string `json:"rule"`
}

// Unwrap use with errors.Unwrap to get raw error
func (e *ValidationError) Unwrap() error {
	return e.error
}

func NewValidationError(err error, fieldName string, rule string) *ValidationError {
	return &ValidationError{error: err, FieldName: fieldName, Rule: rule}
}

var (
	messageParamRequiredError = errors.New("param error, param `message` is required")
	valueParamRequiredError   = errors.New("param error, param `value` is required")
	unsupportedError          = errors.New("unsupported error")
)

func IsValidationError(err error) bool {
	var ok bool
	if err != nil {
		_, ok = err.(*ValidationError)
	}
	return ok
}
