package internal

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

type Serializer interface {
	Encode(writer io.Writer, v interface{}) error
	Decode(reader io.Reader, v interface{}) error
}

type JsonSerializer struct{}

func (j JsonSerializer) Encode(writer io.Writer, v interface{}) error {
	return json.NewEncoder(writer).Encode(v)
}

func (j JsonSerializer) Decode(reader io.Reader, v interface{}) error {
	return json.NewDecoder(reader).Decode(v)
}

type XmlSerializer struct{}

func (x XmlSerializer) Encode(writer io.Writer, v interface{}) error {
	return xml.NewEncoder(writer).Encode(v)
}

func (x XmlSerializer) Decode(reader io.Reader, v interface{}) error {
	return xml.NewDecoder(reader).Decode(v)
}
