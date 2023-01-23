// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package renders

import (
	"net/http"

	"github.com/eatmoreapple/regia/internal"
)

type JsonRender struct {
	Serializer internal.Serializer
}

func (j JsonRender) WriterHeader(writer http.ResponseWriter, code int) {
	writeContentType(writer, "application/json; charset=utf-8")
	writeHeader(writer, code)
}

func (j JsonRender) Render(writer http.ResponseWriter, data interface{}) error {
	return j.Serializer.Encode(writer, data)
}
