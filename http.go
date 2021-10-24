// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"net/http"
)

const (
	contentType         = "Content-Type"
	jsonContentType     = "application/json;charset=utf-8"
	textHtmlContentType = "text/html;charset=utf-8"
	textXmlContentType  = "text/xml;charset=utf-8"
	octetStream         = "application/octet-stream"
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

func writeContentType(writer http.ResponseWriter, cT string) {
	writer.Header().Del(contentType)
	writer.Header().Set(contentType, cT)
}
