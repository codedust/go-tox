package golibtox

// Ideas and some funcs from https://code.google.com/p/go-sqlite/source/browse/go1/sqlite3/util.go - Copyright to them...

/*
#include <tox/tox.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

// goBytes returns a Go representation of an n-byte C array.
func goBytes(p unsafe.Pointer, n C.int) (b []byte) {
	if n > 0 {
		h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
		h.Data = uintptr(p)
		h.Len = int(n)
		h.Cap = int(n)
	}
	return
}

//export hook_callback_friend_request
func hook_callback_friend_request(t unsafe.Pointer, publicKey *C.uint8_t, data *C.uint8_t, length C.uint16_t, userdata unsafe.Pointer) {
	friendRequestFunc(goBytes((unsafe.Pointer)(publicKey), C.TOX_FRIEND_ADDRESS_SIZE), goBytes((unsafe.Pointer)(data), (C.int)(length)), (uint16)(length))
}

//export hook_callback_friend_message
func hook_callback_friend_message(t unsafe.Pointer, friendNumber C.int, message *C.uint8_t, length C.uint16_t, userdata unsafe.Pointer) {
	friendMessageFunc((int)(friendNumber), goBytes((unsafe.Pointer)(message), (C.int)(length)), (uint16)(length))
}
