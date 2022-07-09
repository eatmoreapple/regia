// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package renders

import "net/http"

func writeContentType(writer http.ResponseWriter, t string) {
	writer.Header().Del("Content-Type")
	writer.Header().Set("Content-Type", t)
}
