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
	friendRequestFunc(goBytes((unsafe.Pointer)(publicKey), FRIEND_ADDRESS_SIZE), goBytes((unsafe.Pointer)(data), (C.int)(length)), uint16(length))
}

//export hook_callback_friend_message
func hook_callback_friend_message(t unsafe.Pointer, friendNumber C.int32_t, message *C.uint8_t, length C.uint16_t, userdata unsafe.Pointer) {
	friendMessageFunc(int32(friendNumber), goBytes((unsafe.Pointer)(message), (C.int)(length)), uint16(length))
}

//export hook_callback_friend_action
func hook_callback_friend_action(t unsafe.Pointer, friendNumber C.int32_t, action *C.uint8_t, length C.uint16_t, userdata unsafe.Pointer) {
	friendActionFunc(int32(friendNumber), goBytes((unsafe.Pointer)(action), (C.int)(length)), uint16(length))
}

//export hook_callback_name_change
func hook_callback_name_change(t unsafe.Pointer, friendNumber C.int32_t, newName *C.uint8_t, length C.uint16_t, userdata unsafe.Pointer) {
	nameChangeFunc(int32(friendNumber), goBytes((unsafe.Pointer)(newName), (C.int)(length)), uint16(length))
}

//export hook_callback_status_message
func hook_callback_status_message(t unsafe.Pointer, friendNumber C.int32_t, newStatus *C.uint8_t, length C.uint16_t, userdata unsafe.Pointer) {
	statusMessageFunc(int32(friendNumber), goBytes((unsafe.Pointer)(newStatus), (C.int)(length)), uint16(length))
}

//export hook_callback_user_status
func hook_callback_user_status(t unsafe.Pointer, friendNumber C.int32_t, status C.uint8_t, userdata unsafe.Pointer) {
	userStatusFunc(int32(friendNumber), UserStatus(status))
}

//export hook_callback_typing_change
func hook_callback_typing_change(t unsafe.Pointer, friendNumber C.int32_t, isTyping C.uint8_t, userdata unsafe.Pointer) {
	typing := false
	if isTyping == 1 {
		typing = true
	}
	typingChangeFunc(int32(friendNumber), typing)
}

//export hook_callback_read_receipt
func hook_callback_read_receipt(t unsafe.Pointer, friendNumber C.int32_t, receipt C.uint32_t, userdata unsafe.Pointer) {
	readReceiptFunc(int32(friendNumber), uint32(receipt))
}

//export hook_callback_connection_status
func hook_callback_connection_status(t unsafe.Pointer, friendNumber C.int32_t, status C.uint8_t, userdata unsafe.Pointer) {
	goStatus := false
	if status == 1 {
		goStatus = true
	}
	connectionStatusFunc(int32(friendNumber), goStatus)
}
