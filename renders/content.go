// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package renders

import (
	"io"
	"net/http"
	"time"
)

type ContentRender struct {
	Name    string
	ModTime time.Time
	Request *http.Request
	Content io.ReadSeeker
}

func (c ContentRender) WriterHeader(writer http.ResponseWriter, code int) {
	writeHeader(writer, code)
}

func (c ContentRender) Render(writer http.ResponseWriter, _ interface{}) error {
	http.ServeContent(writer, c.Request, c.Name, c.ModTime, c.Content)
	return nil
}
