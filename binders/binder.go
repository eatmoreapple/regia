// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package binders

import (
	"github.com/eatmoreapple/regia/internal"
	"net/http"
	"net/url"
)

type Binder interface {
	Bind(request *http.Request, v interface{}) error
}

type QueryBinder struct{}

func (QueryBinder) Bind(request *http.Request, v interface{}) error {
	binder := URLValueBinder{TagName: formTag, BindTagName: bindTag}
	return binder.BindForm(request.URL.Query(), v)
}

type FormBinder struct{}

func (FormBinder) Bind(request *http.Request, v interface{}) error {
	binder := URLValueBinder{TagName: formTag, BindTagName: bindTag}
	return binder.BindForm(request.Form, v)
}

type MultipartFormBodyBinder struct{}

func (MultipartFormBodyBinder) Bind(request *http.Request, v interface{}) error {
	urlValueBinder := URLValueBinder{TagName: formTag, BindTagName: bindTag}
	binder := HttpMultipartFormBinder{URLValueBinder: urlValueBinder, FieldTag: fileTag}
	return binder.BindMultipartForm(request.MultipartForm, v)
}

type JsonBodyBinder struct {
	Serializer internal.Serializer
}

func (j JsonBodyBinder) Bind(request *http.Request, v interface{}) error {
	return j.Serializer.Decode(request.Body, v)
}

type XmlBodyBinder struct {
	Serializer internal.Serializer
}

func (j XmlBodyBinder) Bind(request *http.Request, v interface{}) error {
	return j.Serializer.Decode(request.Body, v)
}

type HeaderBinder struct{}

func (h HeaderBinder) Bind(request *http.Request, v interface{}) error {
	values := url.Values(request.Header)
	binder := URLValueBinder{TagName: headerTag, BindTagName: bindTag}
	return binder.BindForm(values, v)
}

type URIParamContextKey struct{}

type URIBinder struct {
	Values url.Values
}

func (u URIBinder) Bind(request *http.Request, v interface{}) error {
	binder := URLValueBinder{TagName: uriTag, BindTagName: bindTag}
	return binder.BindForm(u.Values, v)
}
