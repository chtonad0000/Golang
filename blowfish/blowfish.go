//go:build !solution

package blowfish

/*
#cgo pkg-config: libcrypto
#cgo CFLAGS: -Wno-deprecated-declarations
#include <openssl/blowfish.h>
*/
import "C"
import (
	_ "crypto/cipher"
	"unsafe"
)

type Blowfish struct {
	key C.BF_KEY
}

func New(key []byte) *Blowfish {
	if len(key) == 0 || len(key) > 56 {
		panic("invalid key length")
	}
	if len(key) < 4 {
		key = append(key, make([]byte, 4-len(key))...)
	}
	bf := &Blowfish{}
	C.BF_set_key(&bf.key, C.int(len(key)), (*C.uchar)(unsafe.Pointer(&key[0])))
	return bf
}

func (bf *Blowfish) BlockSize() int {
	return 8
}

func (bf *Blowfish) Encrypt(dst, src []byte) {
	if len(src) != 8 {
		panic("src block must be 8 bytes")
	}
	if len(dst) != 8 {
		panic("dst block must be 8 bytes")
	}
	C.BF_ecb_encrypt((*C.uchar)(unsafe.Pointer(&src[0])), (*C.uchar)(unsafe.Pointer(&dst[0])), &bf.key, C.BF_ENCRYPT)
}

func (bf *Blowfish) Decrypt(dst, src []byte) {
	if len(src) != 8 {
		panic("src block must be 8 bytes")
	}
	if len(dst) != 8 {
		panic("dst block must be 8 bytes")
	}
	C.BF_ecb_encrypt((*C.uchar)(unsafe.Pointer(&src[0])), (*C.uchar)(unsafe.Pointer(&dst[0])), &bf.key, C.BF_DECRYPT)
}
