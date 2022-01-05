package binders

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"
	"reflect"
	"strings"
)

const (
	pass    = "-"
	formTag = "form"
	bindTag = "bind"
	fileTag = "file"
)

// EmptyMultipartFormError may be used outside
var EmptyMultipartFormError = errors.New("nil *multipart.Form got")

type BindMethod func(reflect.Value, []string) error

type formBinder interface {
	BindForm(values url.Values, v interface{}) error
}

type multipartFormBinder interface {
	BindMultipartForm(form *multipart.Form, v interface{}) error
}

type UrlFormBinder struct {
	TagName     string
	BindTagName string
	BindMethods map[string]BindMethod
}

func (f UrlFormBinder) BindForm(form url.Values, v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return errors.New("pointer type required")
	}
	value = value.Elem()
	t := reflect.TypeOf(v).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := t.Field(i)

		tag, exist := field.Tag.Lookup(f.TagName)
		if !exist {
			// set default tag value
			tag = field.Name
		}

		tags := strings.Split(tag, ",")

		var (
			formKey      string
			defFormValue []string
		)

		if len(tags) == 1 {
			formKey = tags[0]
		} else {
			formKey = tags[0]
			defFormValue = tags[1:]
		}

		if formKey == pass {
			continue
		}

		if formValue, exist := form[formKey]; exist {
			// if we need custom bind by self
			customBindTag, found := field.Tag.Lookup(f.BindTagName)
			if found {
				// custom bind by self
				if method := f.BindMethods[customBindTag]; method != nil {
					if err := method(value.Field(i), formValue); err != nil {
						return err
					}
				} else {
					return errors.New("no method named " + customBindTag)
				}
				continue
			}
			if err := bind(value.Field(i), formValue); err != nil {
				// try to bind default value
				if len(defFormValue) > 0 {
					if err = bind(value.Field(i), defFormValue); err != nil {
						return err
					}
				}
				return err
			}
		} else {
			// try to bind default value
			if len(defFormValue) > 0 {
				if err := bind(value.Field(i), defFormValue); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (f *UrlFormBinder) AddBindMethod(name string, method BindMethod) error {
	if _, exist := f.BindMethods[name]; exist {
		return fmt.Errorf("%s already exist", name)
	}
	if f.BindMethods == nil {
		f.BindMethods = make(map[string]BindMethod)
	}
	f.BindMethods[name] = method
	return nil
}

type HttpMultipartFormBinder struct {
	*UrlFormBinder
	FieldTag string
}

func (m *HttpMultipartFormBinder) BindMultipartForm(form *multipart.Form, v interface{}) error {
	if form == nil {
		return EmptyMultipartFormError
	}
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return errors.New("pointer type required")
	}
	value = value.Elem()
	t := reflect.TypeOf(v).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := t.Field(i)

		tag, exist := field.Tag.Lookup(m.TagName)
		if !exist {
			// set default tag value
			tag = field.Name
		}

		tags := strings.Split(tag, ",")

		var (
			formKey      string
			defFormValue []string
		)

		if len(tags) == 1 {
			formKey = tags[0]
		} else {
			formKey = tags[0]
			defFormValue = tags[1:]
		}

		if formKey == pass {
			continue
		}

		if formValue, exist := form.Value[formKey]; exist {
			// if we need custom bind by self
			customBindTag, found := field.Tag.Lookup(m.BindTagName)
			if found {
				// custom bind by self
				if method := m.BindMethods[customBindTag]; method != nil {
					if err := method(value.Field(i), formValue); err != nil {
						return err
					}
				}
				continue
			}
			if err := bind(value.Field(i), formValue); err != nil {
				// try to bind default value
				if len(defFormValue) > 0 {
					if err = bind(value.Field(i), defFormValue); err != nil {
						return err
					}
				}
				return err
			}
			continue
		} else {
			// try to bind default value
			if len(defFormValue) > 0 {
				if err := bind(value.Field(i), defFormValue); err != nil {
					return err
				}
			}
		}

		switch value.Field(i).Interface().(type) {
		case *multipart.FileHeader, []*multipart.FileHeader:
			fileTag, exist := field.Tag.Lookup(m.FieldTag)
			if !exist {
				fileTag = field.Name
			}
			if fileTag == pass {
				break
			}
			if files, exist := form.File[fileTag]; exist {
				if err := bindFile(value.Field(i), files); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

var (
	DefaultFormBinder                              = &UrlFormBinder{TagName: formTag, BindTagName: bindTag}
	DefaultMultipartFormBinder                     = &HttpMultipartFormBinder{UrlFormBinder: DefaultFormBinder, FieldTag: fileTag}
	_                          formBinder          = DefaultFormBinder
	_                          multipartFormBinder = DefaultMultipartFormBinder
)
