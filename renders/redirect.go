// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package renders

import "net/http"

type RedirectRender struct {
	RedirectURL string
	Code        int
	Request     *http.Request
}

func (r RedirectRender) WriterHeader(_ http.ResponseWriter, _ int) {}

func (r RedirectRender) Render(writer http.ResponseWriter, data interface{}) error {
	if r.Code == 0 {
		r.Code = http.StatusFound
	}
	http.Redirect(writer, r.Request, r.RedirectURL, r.Code)
	return nil
}
