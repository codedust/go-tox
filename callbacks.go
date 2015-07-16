package gotox

/*
#include <tox/tox.h>
#include "hooks.c"
*/
import "C"
import "unsafe"

/* This event is triggered whenever there is a change in the DHT connection
 * state. When disconnected, a client may choose to call tox_bootstrap again, to
 * reconnect to the DHT. Note that this state may frequently change for short
 * amounts of time. Clients should therefore not immediately bootstrap on
 * receiving a disconnect. */
type OnSelfConnectionStatusChanges func(tox *Tox, status ToxConnection)

/* This event is triggered when a friend changes their name. */
type OnFriendNameChanges func(tox *Tox, friendnumber uint32, name string)

/* This event is triggered when a friend changes their status message. */
type OnFriendStatusMessageChanges func(tox *Tox, friendnumber uint32, message string)

/* This event is triggered when a friend changes their user status. */
type OnFriendStatusChanges func(tox *Tox, friendnumber uint32, userstatus ToxUserStatus)

/* This event is triggered when a friend goes offline after having been online,
 * or when a friend goes online.
 *
 * This callback is not called when adding friends. It is assumed that when
 * adding friends, their connection status is initially offline. */
type OnFriendConnectionStatusChanges func(tox *Tox, friendnumber uint32, connectionstatus ToxConnection)

/* This event is triggered when a friend starts or stops typing. */
type OnFriendTypingChanges func(tox *Tox, friendnumber uint32, istyping bool)

/* This event is triggered when the friend receives the message with the
 * corresponding message ID. */
type OnFriendReadReceipt func(tox *Tox, friendnumber uint32, messageid uint32)

/* This event is triggered when a friend request is received. */
type OnFriendRequest func(tox *Tox, publickey []byte, message string)

/* This event is triggered when a message from a friend is received. */
type OnFriendMessage func(tox *Tox, friendnumber uint32, messagetype ToxMessageType, message string)

/* This event is triggered when a file control command is received from a
 * friend. */
type OnFileRecvControl func(tox *Tox, friendnumber uint32, filenumber uint32, filecontrol ToxFileControl)

/* This event is triggered when Core is ready to send more file data. */
type OnFileChunkRequest func(tox *Tox, friendnumber uint32, filenumber uint32, position uint64, length uint64)

/* This event is triggered when a file transfer request is received. */
type OnFileRecv func(tox *Tox, friendnumber uint32, filenumber uint32, kind ToxFileKind, filesize uint64, filename string)

/* This event is first triggered when a file transfer request is received, and
 * subsequently when a chunk of file data for an accepted request was received.
 */
type OnFileRecvChunk func(tox *Tox, friendnumber uint32, filenumber uint32, position uint64, data []byte)

/* This event is triggered when a lossy packet is received from a friend. */
type OnFriendLossyPacket func(tox *Tox, friendnumber uint32, data []byte)

/* This event is triggered when a lossless packet is received from a friend. */
type OnFriendLosslessPacket func(tox *Tox, friendnumber uint32, data []byte)

/*
 * Functions to register the callbacks.
 */

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
