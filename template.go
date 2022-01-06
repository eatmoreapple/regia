package regia

import (
	"errors"
	"github.com/eatmoreapple/regia/renders"
	"html/template"
)

type TemplateLoader interface {
	Load(name string) (renders.Render, error) // Load a template by name
	ParseGlob(pattern string) error           // ParseGlob templates by pattern
}

type HTMLLoader struct {
	*template.Template
}

func (h *HTMLLoader) Load(name string) (renders.Render, error) {
	t := h.Lookup(name)
	if t == nil {
		return nil, errors.New("template not found")
	}
	render := renders.Template{Template: t}
	return render, nil
}

func (h *HTMLLoader) ParseGlob(pattern string) error {
	var err error
	h.Template, err = template.ParseGlob(pattern)
	return err
}
