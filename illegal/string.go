//go:build !solution

package illegal

import "unsafe"

func StringFromBytes(b []byte) string {
	if len(b) == 0 {
		return ""
	}

	return *(*string)(unsafe.Pointer(&b))
}
