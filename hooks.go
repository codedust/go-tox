package golibtox

// Ideas from https://code.google.com/p/go-sqlite/source/browse/go1/sqlite3/util.go

/*
#include <tox/tox.h>
*/
import "C"
import "unsafe"

//export hook_callback_friend_request
func hook_callback_friend_request(t unsafe.Pointer, publicKey *C.uint8_t, data *C.uint8_t, length C.uint16_t, userdata unsafe.Pointer) {
	friendRequestFunc(C.GoBytes((unsafe.Pointer)(publicKey), FRIEND_ADDRESS_SIZE), C.GoBytes((unsafe.Pointer)(data), (C.int)(length)), uint16(length))
}

//export hook_callback_friend_message
func hook_callback_friend_message(t unsafe.Pointer, friendNumber C.int32_t, message *C.uint8_t, length C.uint16_t, userdata unsafe.Pointer) {
	friendMessageFunc(int32(friendNumber), C.GoBytes((unsafe.Pointer)(message), (C.int)(length)), uint16(length))
}

//export hook_callback_friend_action
func hook_callback_friend_action(t unsafe.Pointer, friendNumber C.int32_t, action *C.uint8_t, length C.uint16_t, userdata unsafe.Pointer) {
	friendActionFunc(int32(friendNumber), C.GoBytes((unsafe.Pointer)(action), (C.int)(length)), uint16(length))
}

//export hook_callback_name_change
func hook_callback_name_change(t unsafe.Pointer, friendNumber C.int32_t, newName *C.uint8_t, length C.uint16_t, userdata unsafe.Pointer) {
	nameChangeFunc(int32(friendNumber), C.GoBytes((unsafe.Pointer)(newName), (C.int)(length)), uint16(length))
}

//export hook_callback_status_message
func hook_callback_status_message(t unsafe.Pointer, friendNumber C.int32_t, newStatus *C.uint8_t, length C.uint16_t, userdata unsafe.Pointer) {
	statusMessageFunc(int32(friendNumber), C.GoBytes((unsafe.Pointer)(newStatus), (C.int)(length)), uint16(length))
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

//export hook_callback_file_send_request
func hook_callback_file_send_request(t unsafe.Pointer, friendNumber C.int32_t, filenumber C.uint8_t, filesize C.uint64_t, filename unsafe.Pointer, filenameLength C.uint16_t, userdata unsafe.Pointer) {
	fileSendRequestFunc(int32(friendNumber), uint8(filenumber), uint64(filesize), C.GoBytes(unsafe.Pointer(filename), C.int(filenameLength)), uint16(filenameLength))
}

//export hook_callback_file_control
func hook_callback_file_control(t unsafe.Pointer, friendNumber C.int32_t, sending C.uint8_t, filenumber C.uint8_t, fileControl C.uint8_t, data unsafe.Pointer, length C.uint16_t, userdata unsafe.Pointer) {
	goSending := false
	if sending == 1 {
		goSending = true
	}
	fileControlFunc(int32(friendNumber), goSending, uint8(filenumber), FileControl(fileControl), C.GoBytes(unsafe.Pointer(data), C.int(length)), uint16(length))
}

//export hook_callback_file_data
func hook_callback_file_data(t unsafe.Pointer, friendNumber C.int32_t, filenumber C.uint8_t, data unsafe.Pointer, length C.uint16_t, userdata unsafe.Pointer) {
	fileDataFunc(int32(friendNumber), uint8(filenumber), C.GoBytes(unsafe.Pointer(data), C.int(length)), uint16(length))
}
