// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"errors"
	"github.com/eatmoreapple/regia/renders"
	"html/template"
)

type HTMLLoader interface {
	Load(name string) (renders.Render, error) // Load a template by name
	ParseGlob(pattern string) error           // ParseGlob templates by pattern
}

type TemplateLoader struct {
	*template.Template
}

func (h *TemplateLoader) Load(name string) (renders.Render, error) {
	t := h.Lookup(name)
	if t == nil {
		return nil, errors.New("template not found")
	}
	render := renders.Template{Template: t}
	return render, nil
}

func (h *TemplateLoader) ParseGlob(pattern string) error {
	var err error
	h.Template, err = template.ParseGlob(pattern)
	return err
}
