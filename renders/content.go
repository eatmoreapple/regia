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

func (c ContentRender) Render(writer http.ResponseWriter, data interface{}) error {
	http.ServeContent(writer, c.Request, c.Name, c.ModTime, c.Content)
	return nil
}
