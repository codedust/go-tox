package golibtox

/*
#include <tox/tox.h>
*/
import "C"
import "unsafe"

//export hook_callback_friend_request
func hook_callback_friend_request(t unsafe.Pointer, publicKey *C.uint8_t, data *C.uint8_t, length C.uint16_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendRequest((*Tox)(tox), C.GoBytes((unsafe.Pointer)(publicKey), CLIENT_ID_SIZE), C.GoBytes((unsafe.Pointer)(data), (C.int)(length)), uint16(length))
}

//export hook_callback_friend_message
func hook_callback_friend_message(t unsafe.Pointer, friendnumber C.int32_t, message *C.uint8_t, length C.uint16_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendMessage((*Tox)(tox), int32(friendnumber), C.GoBytes((unsafe.Pointer)(message), (C.int)(length)), uint16(length))
}

//export hook_callback_friend_action
func hook_callback_friend_action(t unsafe.Pointer, friendnumber C.int32_t, action *C.uint8_t, length C.uint16_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendAction((*Tox)(tox), int32(friendnumber), C.GoBytes((unsafe.Pointer)(action), (C.int)(length)), uint16(length))
}

//export hook_callback_name_change
func hook_callback_name_change(t unsafe.Pointer, friendnumber C.int32_t, name *C.uint8_t, length C.uint16_t, tox unsafe.Pointer) {
	(*Tox)(tox).onNameChange((*Tox)(tox), int32(friendnumber), C.GoBytes((unsafe.Pointer)(name), (C.int)(length)), uint16(length))
}

//export hook_callback_status_message
func hook_callback_status_message(t unsafe.Pointer, friendnumber C.int32_t, status *C.uint8_t, length C.uint16_t, tox unsafe.Pointer) {
	(*Tox)(tox).onStatusMessage((*Tox)(tox), int32(friendnumber), C.GoBytes((unsafe.Pointer)(status), (C.int)(length)), uint16(length))
}

//export hook_callback_user_status
func hook_callback_user_status(t unsafe.Pointer, friendnumber C.int32_t, userstatus C.uint8_t, tox unsafe.Pointer) {
	(*Tox)(tox).onUserStatus((*Tox)(tox), int32(friendnumber), UserStatus(userstatus))
}

//export hook_callback_typing_change
func hook_callback_typing_change(t unsafe.Pointer, friendnumber C.int32_t, ctyping C.uint8_t, tox unsafe.Pointer) {
	typing := false
	if ctyping == 1 {
		typing = true
	}
	(*Tox)(tox).onTypingChange((*Tox)(tox), int32(friendnumber), typing)
}

//export hook_callback_read_receipt
func hook_callback_read_receipt(t unsafe.Pointer, friendnumber C.int32_t, receipt C.uint32_t, tox unsafe.Pointer) {
	(*Tox)(tox).onReadReceipt((*Tox)(tox), int32(friendnumber), uint32(receipt))
}

//export hook_callback_connection_status
func hook_callback_connection_status(t unsafe.Pointer, friendnumber C.int32_t, conline C.uint8_t, tox unsafe.Pointer) {
	online := false
	if conline == 1 {
		online = true
	}
	(*Tox)(tox).onConnectionStatus((*Tox)(tox), int32(friendnumber), online)
}

//export hook_callback_file_send_request
func hook_callback_file_send_request(t unsafe.Pointer, friendnumber C.int32_t, filenumber C.uint8_t, filesize C.uint64_t, filename unsafe.Pointer, filenameLength C.uint16_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFileSendRequest((*Tox)(tox), int32(friendnumber), uint8(filenumber), uint64(filesize), C.GoBytes(unsafe.Pointer(filename), C.int(filenameLength)), uint16(filenameLength))
}

//export hook_callback_file_control
func hook_callback_file_control(t unsafe.Pointer, friendnumber C.int32_t, csending C.uint8_t, filenumber C.uint8_t, fileControl C.uint8_t, data unsafe.Pointer, length C.uint16_t, tox unsafe.Pointer) {
	sending := false
	if csending == 1 {
		sending = true
	}
	(*Tox)(tox).onFileControl((*Tox)(tox), int32(friendnumber), sending, uint8(filenumber), FileControl(fileControl), C.GoBytes(unsafe.Pointer(data), C.int(length)), uint16(length))
}

//export hook_callback_file_data
func hook_callback_file_data(t unsafe.Pointer, friendnumber C.int32_t, filenumber C.uint8_t, data unsafe.Pointer, length C.uint16_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFileData((*Tox)(tox), int32(friendnumber), uint8(filenumber), C.GoBytes(unsafe.Pointer(data), C.int(length)), uint16(length))
}
