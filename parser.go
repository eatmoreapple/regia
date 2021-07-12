package regia

import (
	"errors"
	"net/http"
	"strings"
)

const (
	minePostForm          = "application/x-www-form-urlencoded"
	mimeJson              = "application/json"
	mimeMultipartPostForm = "multipart/form-data"
)

var parseError = errors.New(http.StatusText(http.StatusBadRequest))

type Parser interface {
	// Parse parse incoming bytestream and return a error if parse failed
	Parse(context *Context, v interface{}) error
	// Match define that if we should parse this request
	Match(context *Context) bool
}

type Parsers struct {
	ps  []Parser
	err error
}

// Parse start to parse request data
func (p *Parsers) Parse(context *Context, v interface{}) error {
	for _, parse := range p.ps {
		if match := parse.Match(context); match {
			return parse.Parse(context, v)
		}
	}
	if p.err != nil {
		return p.err
	}
	return parseError
}

func (p *Parsers) AddParser(parser ...Parser) {
	if p.ps == nil {
		p.ps = make([]Parser, 0)
	}
	p.ps = append(p.ps, parser...)
}

func (p *Parsers) SetParseError(e error) {
	p.err = e
}

// FormParser Parser for form data.
type FormParser struct{}

func (f FormParser) Parse(context *Context, v interface{}) error {
	return context.BindForm(v)
}

func (f FormParser) Match(context *Context) bool {
	return strings.ToLower(context.Request.Header.Get(contentType)) == minePostForm
}

// JsonParser Parses JSON-serialized data.
type JsonParser struct{}

func (j JsonParser) Parse(context *Context, v interface{}) error {
	return context.BindJSON(v)
}

func (j JsonParser) Match(context *Context) bool {
	return strings.ToLower(context.Request.Header.Get(contentType)) == mimeJson
}

// MultipartFormParser Parser for multipart form data, which may include file data.
type MultipartFormParser struct{}

func (m MultipartFormParser) Parse(context *Context, v interface{}) error {
	return context.BindMultipartForm(v)
}

func (m MultipartFormParser) Match(context *Context) bool {
	return strings.Contains(strings.ToLower(context.Request.Header.Get(contentType)), mimeMultipartPostForm)
}
