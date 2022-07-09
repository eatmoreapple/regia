// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package binders

import (
	"github.com/eatmoreapple/regia/internal"
	"mime/multipart"
	"net/http"
	"net/url"
)

func BindForm(form url.Values, v interface{}) error {
	return DefaultFormBinder.BindForm(form, v)
}

func BindMultipartForm(form *multipart.Form, v interface{}) error {
	return DefaultMultipartFormBinder.BindMultipartForm(form, v)
}

type Binder interface {
	Bind(request *http.Request, v interface{}) error
}

type QueryBinder struct{}

func (QueryBinder) Bind(request *http.Request, v interface{}) error {
	query := request.URL.Query()
	return BindForm(query, v)
}

type FormBinder struct{}

func (FormBinder) Bind(request *http.Request, v interface{}) error {
	return BindForm(request.Form, v)
}

type MultipartFormBodyBinder struct{}

func (MultipartFormBodyBinder) Bind(request *http.Request, v interface{}) error {
	return BindMultipartForm(request.MultipartForm, v)
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
	return DefaultHeaderBinder.BindForm(values, v)
}
