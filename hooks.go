package gotox

//#include <tox/tox.h>
import "C"
import "unsafe"

//export hook_callback_self_connection_status
func hook_callback_self_connection_status(t unsafe.Pointer, status C.enum_TOX_CONNECTION, tox unsafe.Pointer) {
	(*Tox)(tox).onSelfConnectionStatusChanges((*Tox)(tox), ToxConnection(status))
}

//export hook_callback_friend_name
func hook_callback_friend_name(t unsafe.Pointer, friendnumber C.uint32_t, name *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendNameChanges((*Tox)(tox), uint32(friendnumber), string(C.GoBytes((unsafe.Pointer)(name), (C.int)(length))))
}

//export hook_callback_friend_status_message
func hook_callback_friend_status_message(t unsafe.Pointer, friendnumber C.uint32_t, message *C.uint8_t, length C.uint16_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendStatusMessageChanges((*Tox)(tox), uint32(friendnumber), string(C.GoBytes((unsafe.Pointer)(message), (C.int)(length))))
}

//export hook_callback_friend_status
func hook_callback_friend_status(t unsafe.Pointer, friendnumber C.uint32_t, status C.enum_TOX_USER_STATUS, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendStatusChanges((*Tox)(tox), uint32(friendnumber), ToxUserStatus(status))
}

//export hook_callback_friend_connection_status
func hook_callback_friend_connection_status(t unsafe.Pointer, friendnumber C.uint32_t, status C.enum_TOX_CONNECTION, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendConnectionStatusChanges((*Tox)(tox), uint32(friendnumber), ToxConnection(status))
}

//export hook_callback_friend_typing
func hook_callback_friend_typing(t unsafe.Pointer, friendnumber C.uint32_t, istyping C._Bool, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendTypingChanges((*Tox)(tox), uint32(friendnumber), bool(istyping))
}

//export hook_callback_friend_read_receipt
func hook_callback_friend_read_receipt(t unsafe.Pointer, friendnumber C.uint32_t, messageid C.uint32_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendReadReceipt((*Tox)(tox), uint32(friendnumber), uint32(messageid))
}

//export hook_callback_friend_request
func hook_callback_friend_request(t unsafe.Pointer, publicKey *C.uint8_t, message *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendRequest((*Tox)(tox), C.GoBytes((unsafe.Pointer)(publicKey), TOX_PUBLIC_KEY_SIZE), string(C.GoBytes((unsafe.Pointer)(message), (C.int)(length))))
}

//export hook_callback_friend_message
func hook_callback_friend_message(t unsafe.Pointer, friendnumber C.uint32_t, messagetype C.enum_TOX_MESSAGE_TYPE, message *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendMessage((*Tox)(tox), uint32(friendnumber), ToxMessageType(messagetype), string(C.GoBytes((unsafe.Pointer)(message), (C.int)(length))))
}

//export hook_callback_file_recv_control
func hook_callback_file_recv_control(t unsafe.Pointer, friendnumber C.uint32_t, filenumber C.uint32_t, control C.enum_TOX_FILE_CONTROL, tox unsafe.Pointer) {
	(*Tox)(tox).onFileRecvControl((*Tox)(tox), uint32(friendnumber), uint32(filenumber), ToxFileControl(control))
}

//export hook_callback_file_chunk_request
func hook_callback_file_chunk_request(t unsafe.Pointer, friendnumber C.uint32_t, filenumber C.uint32_t, position C.uint64_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFileChunkRequest((*Tox)(tox), uint32(friendnumber), uint32(filenumber), uint64(position), uint64(length))
}

//export hook_callback_file_recv
func hook_callback_file_recv(t unsafe.Pointer, friendnumber C.uint32_t, filenumber C.uint32_t, kind C.uint32_t, filesize C.uint64_t, filename *C.uint8_t, filenamelength C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFileRecv((*Tox)(tox), uint32(friendnumber), uint32(filenumber), ToxFileKind(kind), uint64(filesize), string(C.GoBytes((unsafe.Pointer)(filename), (C.int)(filenamelength))))
}

//export hook_callback_file_recv_chunk
func hook_callback_file_recv_chunk(t unsafe.Pointer, friendnumber C.uint32_t, filenumber C.uint32_t, position C.uint64_t, data *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFileRecvChunk((*Tox)(tox), uint32(friendnumber), uint32(filenumber), uint64(position), C.GoBytes((unsafe.Pointer)(data), (C.int)(length)))
}

//export hook_callback_friend_lossy_packet
func hook_callback_friend_lossy_packet(t unsafe.Pointer, friendnumber C.uint32_t, data *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendLossyPacket((*Tox)(tox), uint32(friendnumber), C.GoBytes((unsafe.Pointer)(data), (C.int)(length)))
}

//export hook_callback_friend_lossless_packet
func hook_callback_friend_lossless_packet(t unsafe.Pointer, friendnumber C.uint32_t, data *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendLosslessPacket((*Tox)(tox), uint32(friendnumber), C.GoBytes((unsafe.Pointer)(data), (C.int)(length)))
}
