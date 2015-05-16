package gotox

/*
#include <tox/tox.h>
#include "hooks.c"
*/
import "C"
import "unsafe"

func (t *Tox) CallbackSelfConnectionStatusChanges(f OnSelfConnectionStatusChanges) {
	if t.tox != nil {
		t.onSelfConnectionStatusChanges = f
		C.set_callback_self_connection_status(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFriendNameChanges(f OnFriendNameChanges) {
	if t.tox != nil {
		t.onFriendNameChanges = f
		C.set_callback_friend_name(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFriendStatusMessageChanges(f OnFriendStatusMessageChanges) {
	if t.tox != nil {
		t.onFriendStatusMessageChanges = f
		C.set_callback_friend_status_message(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFriendStatusChanges(f OnFriendStatusChanges) {
	if t.tox != nil {
		t.onFriendStatusChanges = f
		C.set_callback_friend_status(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFriendConnectionStatusChanges(f OnFriendConnectionStatusChanges) {
	if t.tox != nil {
		t.onFriendConnectionStatusChanges = f
		C.set_callback_friend_connection_status(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFriendTypingChanges(f OnFriendTypingChanges) {
	if t.tox != nil {
		t.onFriendTypingChanges = f
		C.set_callback_friend_typing(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFriendReadReceipt(f OnFriendReadReceipt) {
	if t.tox != nil {
		t.onFriendReadReceipt = f
		C.set_callback_friend_read_receipt(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFriendRequest(f OnFriendRequest) {
	if t.tox != nil {
		t.onFriendRequest = f
		C.set_callback_friend_request(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFriendMessage(f OnFriendMessage) {
	if t.tox != nil {
		t.onFriendMessage = f
		C.set_callback_friend_message(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFileRecvControl(f OnFileRecvControl) {
	if t.tox != nil {
		t.onFileRecvControl = f
		C.set_callback_file_recv_control(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFileChunkRequest(f OnFileChunkRequest) {
	if t.tox != nil {
		t.onFileChunkRequest = f
		C.set_callback_file_chunk_request(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFileRecv(f OnFileRecv) {
	if t.tox != nil {
		t.onFileRecv = f
		C.set_callback_file_recv(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFileRecvChunk(f OnFileRecvChunk) {
	if t.tox != nil {
		t.onFileRecvChunk = f
		C.set_callback_file_recv_chunk(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFriendLossyPacket(f OnFriendLossyPacket) {
	if t.tox != nil {
		t.onFriendLossyPacket = f
		C.set_callback_friend_lossy_packet(t.tox, unsafe.Pointer(t))
	}
}

func (t *Tox) CallbackFriendLosslessPacket(f OnFriendLosslessPacket) {
	if t.tox != nil {
		t.onFriendLosslessPacket = f
		C.set_callback_friend_lossless_packet(t.tox, unsafe.Pointer(t))
	}
}
