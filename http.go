// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"net/http"
)

const (
	MethodGet     = "Get"
	MethodPost    = "Post"
	MethodPut     = "Put"
	MethodPatch   = "Patch"
	MethodDelete  = "Delete"
	MethodHead    = "Head"
	MethodOptions = "Options"
	MethodTrace   = "Trace"
	ALLMethods    = "*"
)

var httpMethods = [...]string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodOptions,
	http.MethodHead,
	http.MethodTrace,
}

// HttpRequestMethodMapping define request method mapping
var HttpRequestMethodMapping = map[string]string{
	MethodPost:    http.MethodPost,
	MethodGet:     http.MethodGet,
	MethodPut:     http.MethodPut,
	MethodPatch:   http.MethodPatch,
	MethodDelete:  http.MethodDelete,
	MethodHead:    http.MethodHead,
	MethodOptions: http.MethodOptions,
	MethodTrace:   http.MethodTrace,
}

// bodyAllowedForStatus reports whether a given response status code
// permits a body. See RFC 7230, section 3.3.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == 204:
		return false
	case status == 304:
		return false
	}
	return true
}
