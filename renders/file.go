// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package renders

import "net/http"

type FileAttachmentRender struct {
	Filename string
	FilePath string
	Request  *http.Request
}

func (f FileAttachmentRender) WriterHeader(writer http.ResponseWriter, code int) {
	writer.Header().Set("Content-Disposition", "attachment; filename=\""+f.Filename+"\"")
	writeHeader(writer, code)
}

func (f FileAttachmentRender) Render(writer http.ResponseWriter, data interface{}) error {
	http.ServeFile(writer, f.Request, f.FilePath)
	return nil
}
