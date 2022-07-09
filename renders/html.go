// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package renders

import (
	"html/template"
	"net/http"
)

type Template struct {
	*template.Template
}

func (t Template) Render(writer http.ResponseWriter, data interface{}) error {
	writeContentType(writer, "text/html; charset=utf-8")
	return t.Execute(writer, data)
}
