// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

type Render interface {
	Render(writer http.ResponseWriter, data interface{}) error
}

type JsonRender struct {
	Marshaller Marshaller
}

func (j JsonRender) Render(writer http.ResponseWriter, v interface{}) error {
	writeContentType(writer, jsonContentType)
	return j.Marshaller.Encode(writer, v)
}

type XmlRender struct {
	Marshaller Marshaller
}

func (j XmlRender) Render(writer http.ResponseWriter, v interface{}) error {
	writeContentType(writer, textXmlContentType)
	return j.Marshaller.Encode(writer, v)
}

type Marshaller interface {
	Decode(reader io.Reader, v interface{}) error
	Encode(writer io.Writer, v interface{}) error
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

	JSONRender = JsonRender{Marshaller: JSON}
	XMLRender  = XmlRender{Marshaller: XML}
)
