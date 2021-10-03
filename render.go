// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"fmt"
	"github.com/eatmoreapple/regia/internal"
	"net/http"
)

type Render interface {
	Render(writer http.ResponseWriter, data interface{}) error
}

type JsonRender struct {
	Serializer internal.Serializer
}

func (j JsonRender) Render(writer http.ResponseWriter, v interface{}) error {
	writeContentType(writer, jsonContentType)
	return j.Serializer.Encode(writer, v)
}

type XmlRender struct {
	Serializer internal.Serializer
}

func (j XmlRender) Render(writer http.ResponseWriter, v interface{}) error {
	writeContentType(writer, textXmlContentType)
	return j.Serializer.Encode(writer, v)
}

type StringRender struct {
	format string
	data   []interface{}
}

func (s StringRender) Render(writer http.ResponseWriter, v interface{}) (err error) {
	writeContentType(writer, textHtmlContentType)
	if len(s.data) > 0 {
		_, err = fmt.Fprintf(writer, s.format, s.data...)
	} else {
		_, err = writer.Write(stringToByte(s.format))
	}
	return err
}
