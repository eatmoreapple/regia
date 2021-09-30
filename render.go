// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type Render interface {
	Render(writer http.ResponseWriter, data interface{}) error
}

type JsonRender struct {
	Encoder Encoder
}

func (j JsonRender) Render(writer http.ResponseWriter, v interface{}) error {
	writeContentType(writer, jsonContentType)
	return j.Encoder.Encode(writer, v)
}

type XmlRender struct {
	Encoder Encoder
}

func (j XmlRender) Render(writer http.ResponseWriter, v interface{}) error {
	writeContentType(writer, textXmlContentType)
	return j.Encoder.Encode(writer, v)
}

type StringRender struct {
	format string
	data   []interface{}
}

func (s StringRender) Render(writer http.ResponseWriter, data interface{}) (err error) {
	writeContentType(writer, textHtmlContentType)
	if len(s.data) > 0 {
		_, err = fmt.Fprintf(writer, s.format, s.data...)
	} else {
		_, err = writer.Write(stringToByte(s.format))
	}
	return err
}

type Encoder interface {
	Encode(writer io.Writer, v interface{}) error
}

type Decoder interface {
	Decode(reader io.Reader, v interface{}) error
}

type Marshaller interface {
	Encoder
	Decoder
}

type JsonMarshal struct{}

func (JsonMarshal) Decode(reader io.Reader, v interface{}) error {
	return json.NewDecoder(reader).Decode(v)
}

func (JsonMarshal) Encode(writer io.Writer, v interface{}) error {
	return json.NewEncoder(writer).Encode(v)
}

type XmlMarshaller struct{}

func (x XmlMarshaller) Decode(reader io.Reader, v interface{}) error {
	return xml.NewDecoder(reader).Decode(v)
}

func (x XmlMarshaller) Encode(writer io.Writer, v interface{}) error {
	return xml.NewEncoder(writer).Encode(v)
}

var (
	XML  Marshaller = XmlMarshaller{}
	JSON Marshaller = JsonMarshal{}

	JSONRender = JsonRender{Encoder: JSON}
	XMLRender  = XmlRender{Encoder: XML}
)
