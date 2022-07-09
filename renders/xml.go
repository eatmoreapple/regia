// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package renders

import (
	"github.com/eatmoreapple/regia/internal"
	"net/http"
)

type XmlRender struct {
	Serializer internal.Serializer
}

func (x XmlRender) Render(writer http.ResponseWriter, data interface{}) error {
	writeContentType(writer, "text/xml;charset=utf-8")
	return x.Serializer.Encode(writer, data)
}
