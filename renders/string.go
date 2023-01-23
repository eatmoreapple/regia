// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package renders

import (
	"fmt"
	"net/http"

	"github.com/eatmoreapple/regia/internal"
)

type StringRender struct {
	Format string
	Data   []interface{}
}

func (s StringRender) WriterHeader(writer http.ResponseWriter, code int) {
	writeContentType(writer, "text/html; charset=utf-8")
	writeHeader(writer, code)
}

func (s StringRender) Render(writer http.ResponseWriter, v interface{}) (err error) {
	if len(s.Data) > 0 {
		_, err = fmt.Fprintf(writer, s.Format, s.Data...)
	} else {
		_, err = writer.Write(internal.StringToByte(s.Format))
	}
	return err
}
