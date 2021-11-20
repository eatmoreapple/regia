// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"fmt"
	"github.com/eatmoreapple/regia/internal"
	"net/http"
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
	writeContentType(writer, textXmlContentType)
}

// SetJsonSerializer is a setter for JSON Serializer
func SetJsonSerializer(serializer internal.Serializer) {
	internal.JSON = serializer
}

// SetXmlSerializer is a setter for XML Serializer
func SetXmlSerializer(serializer internal.Serializer) {
	internal.XML = serializer
}
