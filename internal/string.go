// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package internal

import (
	"reflect"
	"unsafe"
)

// StringToByte
// unsafe string to byte
// without memory copy
func StringToByte(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&*(*reflect.StringHeader)(unsafe.Pointer(&s))))
}
