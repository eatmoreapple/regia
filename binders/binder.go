// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package binders

import (
	"mime/multipart"
	"net/url"
)

func BindForm(form url.Values, v interface{}) error {
	return DefaultFormBinder.BindForm(form, v)
}

func BindMultipartForm(form *multipart.Form, v interface{}) error {
	return DefaultMultipartFormBinder.BindMultipartForm(form, v)
}
