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
