// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"fmt"
	"github.com/eatmoreapple/regia/internal"
	"io"
	"net/http"
	"time"
)

// Render define Render to write response data
type Render interface {
	Render(writer http.ResponseWriter, data interface{}) error
	WriteContentType(writer http.ResponseWriter)
}

type JsonRender struct {
	Serializer internal.Serializer
}

func (j JsonRender) Render(writer http.ResponseWriter, v interface{}) error {
	return j.Serializer.Encode(writer, v)
}

func (j JsonRender) WriteContentType(writer http.ResponseWriter) {
	writeContentType(writer, jsonContentType)
}

type XmlRender struct {
	Serializer internal.Serializer
}

func (x XmlRender) Render(writer http.ResponseWriter, v interface{}) error {
	return x.Serializer.Encode(writer, v)
}

func (x XmlRender) WriteContentType(writer http.ResponseWriter) {
	writeContentType(writer, textXmlContentType)
}

type StringRender struct {
	format string
	data   []interface{}
}

func (s StringRender) Render(writer http.ResponseWriter, v interface{}) (err error) {
	if len(s.data) > 0 {
		_, err = fmt.Fprintf(writer, s.format, s.data...)
	} else {
		_, err = writer.Write(stringToByte(s.format))
	}
	return err
}

func (s StringRender) WriteContentType(writer http.ResponseWriter) {
	writeContentType(writer, textHtmlContentType)
}

// SetJsonSerializer is a setter for JSON Serializer
func SetJsonSerializer(serializer internal.Serializer) {
	internal.JSON = serializer
}

// SetXmlSerializer is a setter for XML Serializer
func SetXmlSerializer(serializer internal.Serializer) {
	internal.XML = serializer
}

type FileAttachmentRender struct {
	Filename string
	FilePath string
	Request  *http.Request
}

func (f FileAttachmentRender) Render(writer http.ResponseWriter, data interface{}) error {
	http.ServeFile(writer, f.Request, f.FilePath)
	return nil
}

func (f FileAttachmentRender) WriteContentType(writer http.ResponseWriter) {
	writer.Header().Set("Content-Disposition", "attachment; filename=\""+f.Filename+"\"")
}

type RedirectRender struct {
	RedirectURL string
	Code        int
	Request     *http.Request
}

func (r RedirectRender) Render(writer http.ResponseWriter, data interface{}) error {
	if r.Code == 0 {
		r.Code = http.StatusFound
	}
	http.Redirect(writer, r.Request, r.RedirectURL, r.Code)
	return nil
}

func (r RedirectRender) WriteContentType(writer http.ResponseWriter) {}

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

func (c ContentRender) WriteContentType(writer http.ResponseWriter) {}
