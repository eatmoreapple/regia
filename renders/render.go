// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package renders

import "net/http"

type Render interface {
	Render(writer http.ResponseWriter, data interface{}) error
	WriterHeader(writer http.ResponseWriter, code int)
}
