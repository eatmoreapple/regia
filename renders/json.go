// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package renders

import (
	"github.com/eatmoreapple/regia/internal"
	"net/http"
)

type JsonRender struct {
	Serializer internal.Serializer
}

func (j JsonRender) Render(writer http.ResponseWriter, data interface{}) error {
	writeContentType(writer, "application/json;charset=utf-8")
	return j.Serializer.Encode(writer, data)
}
