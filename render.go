package regia

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type Render interface {
	Render(writer http.ResponseWriter, data interface{}) error
}

type JsonRender struct{}

func (j JsonRender) Render(writer http.ResponseWriter, v interface{}) error {
	writeContentType(writer, jsonContentType)
	data, err := JSON.Marshal(v)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

type XmlRender struct{}

func (j XmlRender) Render(writer http.ResponseWriter, v interface{}) error {
	writeContentType(writer, textXmlContentType)
	data, err := xml.Marshal(v)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

var (
	jsonRender = JsonRender{}
	xmlRender  = XmlRender{}
)

type Marshaller interface {
	Unmarshal(data []byte, v interface{}) error
	Marshal(v interface{}) ([]byte, error)
}

type JsonMarshal struct{}

func (JsonMarshal) Unmarshal(data []byte, v interface{}) error { return json.Unmarshal(data, v) }

func (JsonMarshal) Marshal(v interface{}) ([]byte, error) { return json.Marshal(v) }

var JSON Marshaller = JsonMarshal{}
