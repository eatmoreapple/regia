// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package renders

import (
	"net/http"

	"github.com/eatmoreapple/regia/internal"
)

type XmlRender struct {
	Serializer internal.Serializer
}

func (x XmlRender) WriterHeader(writer http.ResponseWriter, code int) {
	writeContentType(writer, "text/xml; charset=utf-8")
	writeHeader(writer, code)
}

func (x XmlRender) Render(writer http.ResponseWriter, data interface{}) error {
	return x.Serializer.Encode(writer, data)
}
