package regia

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
)

type Render interface {
	Render(writer http.ResponseWriter, data interface{}) error
}

type JsonRender struct{}

func (j JsonRender) Render(writer http.ResponseWriter, v interface{}) error {
	writeContentType(writer, jsonContentType)
	data, err := js.Marshal(v)
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

type jsonTransformer interface {
	Unmarshal(data []byte, v interface{}) error
	Marshal(v interface{}) ([]byte, error)
}

type defaultJsonTransformer struct{}

func (d defaultJsonTransformer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (d defaultJsonTransformer) Marshal(v interface{}) ([]byte, error) { return json.Marshal(v) }

var js jsonTransformer = defaultJsonTransformer{}

func SetJsonTransformer(j jsonTransformer) error {
	if j == nil {
		return errors.New("jsonTransformer can not be nil")
	}
	js = j
	return nil
}
