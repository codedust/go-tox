package golibtox

// WIP - organ 2014

/*
#cgo LDFLAGS: -ltoxcore

#include <tox/tox.h>
#include <stdlib.h>

// Convenient macro:
// Creates the C function to directly register a given callback
#define HOOK(x) \
static void set_##x(Tox *tox, void *t) { \
	tox_##x(tox, hook_##x, t); \
}

void hook_callback_friend_request(Tox*, const uint8_t*, const uint8_t*, uint16_t, void*);
void hook_callback_friend_message(Tox*, int32_t, uint8_t*, uint16_t, void*);
void hook_callback_friend_action(Tox*, int32_t, uint8_t*, uint16_t, void*);
void hook_callback_name_change(Tox*, int32_t, uint8_t*, uint16_t, void*);
void hook_callback_status_message(Tox*, int32_t, uint8_t*, uint16_t, void*);
void hook_callback_user_status(Tox*, int32_t, uint8_t, void*);
void hook_callback_typing_change(Tox*, int32_t, uint8_t, void*);
void hook_callback_read_receipt(Tox*, int32_t, uint32_t, void*);
void hook_callback_connection_status(Tox*, int32_t, uint8_t, void*);
void hook_callback_group_invite(Tox*, int32_t, uint8_t*, void*);
void hook_callback_group_message(Tox*, int, int, uint8_t*, uint16_t, void*);
void hook_callback_group_action(Tox*, int, int, uint8_t*, uint16_t, void*);
void hook_callback_group_namelist_change(Tox*, int, int, uint8_t, void*);
void hook_callback_file_send_request(Tox*, int32_t, uint8_t, uint64_t, uint8_t*, uint16_t, void*);
void hook_callback_file_control(Tox*, int32_t, uint8_t, uint8_t, uint8_t, uint8_t*, uint16_t, void*);
void hook_callback_file_data(Tox*, int32_t, uint8_t, uint8_t*, uint16_t, void*);

HOOK(callback_friend_request)
HOOK(callback_friend_message)
HOOK(callback_friend_action)
HOOK(callback_name_change)
HOOK(callback_status_message)
HOOK(callback_user_status)
HOOK(callback_typing_change)
HOOK(callback_read_receipt)
HOOK(callback_connection_status)
HOOK(callback_group_invite)
HOOK(callback_group_message)
HOOK(callback_group_action)
HOOK(callback_group_namelist_change)
HOOK(callback_file_send_request)
HOOK(callback_file_control)
HOOK(callback_file_data)

*/
import "C"

import (
	"encoding/binary"
	"encoding/hex"
	"sync"
	"time"
	"unsafe"
)

// OnFriendRequest is a callback function called by Tox when receiving a friend request.
type OnFriendRequest func(tox *Tox, publicKey []byte, data []byte, length uint16)

// OnFriendMessage is a callback function called by Tox when receiving a message.
type OnFriendMessage func(tox *Tox, friendnumber int32, message []byte, length uint16)

// OnFriendAction is a callback function called by Tox when receiving an action.
type OnFriendAction func(tox *Tox, friendnumber int32, action []byte, length uint16)

// OnNameChange is a callback function called by Tox when a friend changes his name.
type OnNameChange func(tox *Tox, friendnumber int32, name []byte, length uint16)

// OnStatusMessage is a callback function called by Tox when a friend changes his status message.
type OnStatusMessage func(tox *Tox, friendnumber int32, status []byte, length uint16)

// OnUserStatus is a callback function called by Tox when a friend changes status.
type OnUserStatus func(tox *Tox, friendnumber int32, userstatus UserStatus)

// OnTypingChange is a callback function called by Tox when a friend begins/ends typing.
type OnTypingChange func(tox *Tox, friendnumber int32, typing bool)

// OnReadReceipt is a callback function called by Tox when receiving a read receipt.
type OnReadReceipt func(tox *Tox, friendnumber int32, receipt uint32)

// OnConnectionStatus is a callback function called by Tox when a friend comes online/offline.
type OnConnectionStatus func(tox *Tox, friendnumber int32, online bool)

// OnGroupInvite is a callback function called by Tox when receiving a Groupchat invite.
type OnGroupInvite func(tox *Tox, friendnumber int32, groupPublicKey []byte)

// OnGroupMessage is a callback function called by Tox when receiving a Groupchat message.
type OnGroupMessage func(tox *Tox, groupnumber int, friendgroupnumber int, message []byte, length uint16)

// OnGroupAction is a callback function called by Tox when receiving an action from a Groupchat.
type OnGroupAction func(tox *Tox, groupnumber int, friendgroupnumber int, action []byte, length uint16)

// OnGroupNamelistChange is a callback function called by Tox when a peer connects/disconnects/change name in a Groupchat.
type OnGroupNamelistChange func(tox *Tox, groupnumber int, peernumber int, change ChatChange)

// OnFileSendRequest is a callback function called by Tox when receiving a file send request.
type OnFileSendRequest func(tox *Tox, friendnumber int32, filenumber uint8, filesize uint64, filename []byte, filenameLength uint16)

// OnFileControl is a callback function called by Tox when receiving a file control flag for a given file transfer.
type OnFileControl func(tox *Tox, friendnumber int32, sending bool, filenumber uint8, fileControl FileControl, data []byte, length uint16)

// OnFileData is a callback function called by Tox when data is received for a given file transfer.
type OnFileData func(tox *Tox, friendnumber int32, filenumber uint8, data []byte, length uint16)

// Tox is the main struct.
type Tox struct {
	tox *C.struct_Tox
	mtx sync.Mutex
	// Callbacks
	onFriendRequest       OnFriendRequest
	onFriendMessage       OnFriendMessage
	onFriendAction        OnFriendAction
	onNameChange          OnNameChange
	onStatusMessage       OnStatusMessage
	onUserStatus          OnUserStatus
	onTypingChange        OnTypingChange
	onReadReceipt         OnReadReceipt
	onConnectionStatus    OnConnectionStatus
	onGroupInvite         OnGroupInvite
	onGroupMessage        OnGroupMessage
	onGroupAction         OnGroupAction
	onGroupNamelistChange OnGroupNamelistChange
	onFileSendRequest     OnFileSendRequest
	onFileControl         OnFileControl
	onFileData            OnFileData
}

// New returns a new Tox instance.
func New() (*Tox, error) {
	ctox := C.tox_new(ENABLE_IPV6_DEFAULT)
	if ctox == nil {
		return nil, ErrInit
	}

	t := &Tox{tox: ctox}

	return t, nil
}

// Kill stops a Tox instance.
func (t *Tox) Kill() error {
	if t.tox == nil {
		return ErrBadTox
	}
	C.tox_kill(t.tox)

	return nil
}

// DoInterval returns the time in milliseconds before Do() should be called again.
func (t *Tox) DoInterval() (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	ret := C.tox_do_interval(t.tox)

	return uint32(ret), nil
}

// Do is the main loop needs to be called every DoInterval() milliseconds.
func (t *Tox) Do() error {
	if t.tox == nil {
		return ErrBadTox
	}

	t.mtx.Lock()
	C.tox_do(t.tox)
	t.mtx.Unlock()

	return nil
}

// BootstrapFromAddress resolves address into an IP address. If successful, sends a request to the given node to setup connection.
func (t *Tox) BootstrapFromAddress(address string, port uint16, hexPublicKey string) error {
	if t.tox == nil {
		return ErrBadTox
	}

	caddr := C.CString(address)
	defer C.free(unsafe.Pointer(caddr))

	pubkey, err := hex.DecodeString(hexPublicKey)

	if err != nil {
		return err
	}

	// BigEndian int conversion (Network Byte Order)
	var cport uint16
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, port)
	cport = binary.BigEndian.Uint16(b)

	C.tox_bootstrap_from_address(t.tox, caddr, ENABLE_IPV6_DEFAULT, (C.uint16_t)(cport), (*C.uint8_t)(&pubkey[0]))

	return nil
}

// IsConnected returns true if Tox is connected to the DHT.
func (t *Tox) IsConnected() (bool, error) {
	if t.tox == nil {
		return false, ErrBadTox
	}

	return (C.tox_isconnected(t.tox) == 1), nil
}

// GetAddress returns the public address to give to others.
func (t *Tox) GetAddress() ([]byte, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}

	address := make([]byte, FRIEND_ADDRESS_SIZE)
	C.tox_get_address(t.tox, (*C.uint8_t)(&address[0]))

	return address, nil
}

// AddFriend adds a friend. Data contains a message that is sent along with the request.
// Returns the friend number on succes, or a FriendAddError on failure.
func (t *Tox) AddFriend(address []byte, data []byte) (int32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	if len(address) != FRIEND_ADDRESS_SIZE {
		return 0, ErrArgs
	}

	ret := C.tox_add_friend(t.tox, (*C.uint8_t)(&address[0]), (*C.uint8_t)(&data[0]), (C.uint16_t)(len(data)))

	var faerr error

	switch FriendAddError(ret) {
	case FAERR_TOOLONG:
		faerr = FaerrTooLong
	case FAERR_NOMESSAGE:
		faerr = FaerrNoMessage
	case FAERR_OWNKEY:
		faerr = FaerrOwnKey
	case FAERR_ALREADYSENT:
		faerr = FaerrAlreadySent
	case FAERR_UNKNOWN:
		faerr = FaerrUnkown
	case FAERR_BADCHECKSUM:
		faerr = FaerrBadChecksum
	case FAERR_SETNEWNOSPAM:
		faerr = FaerrSetNewNospam
	case FAERR_NOMEM:
		faerr = FaerrNoMem
	}

	return int32(ret), faerr
}

// AddFriendNorequest adds a friend without sending a request.
// Returns the friend number on success.
func (t *Tox) AddFriendNorequest(clientId []byte) (int32, error) {
	if t.tox == nil {
		return -1, ErrBadTox
	}

	if len(clientId) != CLIENT_ID_SIZE {
		return -1, ErrArgs
	}

	n := C.tox_add_friend_norequest(t.tox, (*C.uint8_t)(&clientId[0]))
	if n == -1 {
		return -1, ErrFuncFail
	}

	return int32(n), nil
}

// GetFriendnumber returns the friend number associated to that clientId.
func (t *Tox) GetFriendNumber(clientId []byte) (int32, error) {
	if t.tox == nil {
		return -1, ErrBadTox
	}
	n := C.tox_get_friend_number(t.tox, (*C.uint8_t)(&clientId[0]))

	return int32(n), nil
}

// GetClientId returns the public key associated to that friendnumber.
func (t *Tox) GetClientId(friendnumber int32) ([]byte, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}
	clientId := make([]byte, CLIENT_ID_SIZE)
	ret := C.tox_get_client_id(t.tox, (C.int32_t)(friendnumber), (*C.uint8_t)(&clientId[0]))

	if ret != 0 {
		return nil, ErrFuncFail
	}

	return clientId, nil
}

// DelFriend removes a friend.
func (t *Tox) DelFriend(friendnumber int32) error {
	if t.tox == nil {
		return ErrBadTox
	}
	ret := C.tox_del_friend(t.tox, (C.int32_t)(friendnumber))

	if ret != 0 {
		return ErrFuncFail
	}

	return nil
}

// GetFriendConnectionStatus returns true if the friend is connected.
func (t *Tox) GetFriendConnectionStatus(friendnumber int32) (bool, error) {
	if t.tox == nil {
		return false, ErrBadTox
	}
	ret := C.tox_get_friend_connection_status(t.tox, (C.int32_t)(friendnumber))
	if ret == -1 {
		return false, ErrFuncFail
	}

	return (int(ret) == 1), nil
}

// FriendExists returns true if a friend exists with given friendnumber.
func (t *Tox) FriendExists(friendnumber int32) (bool, error) {
	if t.tox == nil {
		return false, ErrBadTox
	}
	//int tox_friend_exists(Tox *tox, int32_t friendnumber);
	ret := C.tox_friend_exists(t.tox, (C.int32_t)(friendnumber))

	return (int(ret) == 1), nil
}

// SendMessage sends a message to an online friend.
// Maximum message length is MAX_MESSAGE_LENGTH.
// Returns the message ID if successful, an error otherwise.
func (t *Tox) SendMessage(friendnumber int32, message []byte) (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	n := C.tox_send_message(t.tox, (C.int32_t)(friendnumber), (*C.uint8_t)(&message[0]), (C.uint32_t)(len(message)))
	if n == 0 {
		return 0, ErrFuncFail
	}

	return uint32(n), nil
}

// SendMessageWithID sends a message with a given ID to an online friend.
// Maximum message length is MAX_MESSAGE_LENGTH.
// Returns the message ID if successful, an error otherwise.
func (t *Tox) SendMessageWithID(friendnumber int32, id uint32, message []byte) (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	n := C.tox_send_message_withid(t.tox, (C.int32_t)(friendnumber), (C.uint32_t)(id), (*C.uint8_t)(&message[0]), (C.uint32_t)(len(message)))
	if n == 0 {
		return 0, ErrFuncFail
	}

	return uint32(n), nil
}

// SendAction sends an action to an online friend.
// Maximum action length is MAX_MESSAGE_LENGTH.
// Returns the message ID if successful, an error otherwise.
func (t *Tox) SendAction(friendnumber int32, action []byte) (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	n := C.tox_send_action(t.tox, (C.int32_t)(friendnumber), (*C.uint8_t)(&action[0]), (C.uint32_t)(len(action)))
	if n == 0 {
		return 0, ErrFuncFail
	}

	return uint32(n), nil
}

// SendActionActionWithID sends an action with a given ID to an online friend.
// Maximum action length is MAX_MESSAGE_LENGTH.
// Returns the message ID if successful, an error otherwise.
func (t *Tox) SendActionWithID(friendnumber int32, id uint32, action []byte) (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	n := C.tox_send_message_withid(t.tox, (C.int32_t)(friendnumber), (C.uint32_t)(id), (*C.uint8_t)(&action[0]), (C.uint32_t)(len(action)))
	if n == 0 {
		return 0, ErrFuncFail
	}

	return uint32(n), nil
}

// SetName sets your nickname.
// Maximum name length is MAX_NAME_LENGTH
func (t *Tox) SetName(name string) error {
	if t.tox == nil {
		return ErrBadTox
	}

	ret := C.tox_set_name(t.tox, (*C.uint8_t)(&[]byte(name)[0]), (C.uint16_t)(len(name)))
	if ret != 0 {
		return ErrFuncFail
	}

	return nil
}

// GetSelfName returns your nickname.
func (t *Tox) GetSelfName() (string, error) {
	if t.tox == nil {
		return "", ErrBadTox
	}

	cname := make([]byte, MAX_NAME_LENGTH)

	n := C.tox_get_self_name(t.tox, (*C.uint8_t)(&cname[0]))
	if n == 0 {
		return "", ErrFuncFail
	}

	name := string(cname[:n])

	return name, nil
}

// GetName returns the name of friendnumber.
func (t *Tox) GetName(friendnumber int32) (string, error) {
	if t.tox == nil {
		return "", ErrBadTox
	}

	cname := make([]byte, MAX_NAME_LENGTH)

	n := C.tox_get_name(t.tox, (C.int32_t)(friendnumber), (*C.uint8_t)(&cname[0]))
	if n == -1 {
		return "", ErrFuncFail
	}

	name := string(cname[:n])

	return name, nil
}

// GetNameSize returns the length of the name of friendnumber.
func (t *Tox) GetNameSize(friendnumber int32) (int, error) {
	if t.tox == nil {
		return -1, ErrBadTox
	}

	ret := C.tox_get_name_size(t.tox, (C.int32_t)(friendnumber))
	if ret == -1 {
		return -1, ErrFuncFail
	}

	return int(ret), nil
}

// GetSelfNameSize returns the length of your name.
func (t *Tox) GetSelfNameSize() (int, error) {
	if t.tox == nil {
		return -1, ErrBadTox
	}

	ret := C.tox_get_self_name_size(t.tox)
	if ret == -1 {
		return -1, ErrFuncFail
	}

	return int(ret), nil
}

// SetStatusMessage sets your status message.
// Maximum status length is MAX_STATUSMESSAGE_LENGTH.
func (t *Tox) SetStatusMessage(status []byte) error {
	if t.tox == nil {
		return ErrBadTox
	}

	ret := C.tox_set_status_message(t.tox, (*C.uint8_t)(&status[0]), (C.uint16_t)(len(status)))
	if ret != 0 {
		return ErrFuncFail
	}

	return nil
}

// SetUserStatus sets your userstatus.
func (t *Tox) SetUserStatus(userstatus UserStatus) error {
	if t.tox == nil {
		return ErrBadTox
	}

	ret := C.tox_set_user_status(t.tox, (C.uint8_t)(userstatus))
	if ret != 0 {
		return ErrFuncFail
	}

	return nil
}

// GetStatusMessageSize returns the size of the status of friendnumber.
func (t *Tox) GetStatusMessageSize(friendnumber int32) (int, error) {
	if t.tox == nil {
		return -1, ErrBadTox
	}

	ret := C.tox_get_status_message_size(t.tox, (C.int32_t)(friendnumber))
	if ret == -1 {
		return -1, ErrFuncFail
	}

	return int(ret), nil
}

// GetStatusMessageSize returns the size of your status.
func (t *Tox) GetSelfStatusMessageSize() (int, error) {
	if t.tox == nil {
		return -1, ErrBadTox
	}

	ret := C.tox_get_self_status_message_size(t.tox)
	if ret == -1 {
		return -1, ErrFuncFail
	}

	return int(ret), nil
}

// GetStatusMessage returns the status message of friendnumber.
func (t *Tox) GetStatusMessage(friendnumber int32) ([]byte, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}

	status := make([]byte, MAX_STATUSMESSAGE_LENGTH)

	n := C.tox_get_status_message(t.tox, (C.int32_t)(friendnumber), (*C.uint8_t)(&status[0]), MAX_STATUSMESSAGE_LENGTH)
	if n == -1 {
		return nil, ErrFuncFail
	}

	// Truncate status to n-byte read
	status = status[:n]

	return status, nil
}

// GetSelfStatusMessage returns your status message.
func (t *Tox) GetSelfStatusMessage() ([]byte, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}

	status := make([]byte, MAX_STATUSMESSAGE_LENGTH)

	n := C.tox_get_self_status_message(t.tox, (*C.uint8_t)(&status[0]), MAX_STATUSMESSAGE_LENGTH)
	if n == -1 {
		return nil, ErrFuncFail
	}

	// Truncate status to n-byte read
	status = status[:n]

	return status, nil
}

// GetUserStatus returns the status of friendnumber.
func (t *Tox) GetUserStatus(friendnumber int32) (UserStatus, error) {
	if t.tox == nil {
		return USERSTATUS_INVALID, ErrBadTox
	}
	n := C.tox_get_user_status(t.tox, (C.int32_t)(friendnumber))

	return UserStatus(n), nil
}

// GetSelfUserStatus returns your status.
func (t *Tox) GetSelfUserStatus() (UserStatus, error) {
	if t.tox == nil {
		return USERSTATUS_INVALID, ErrBadTox
	}
	n := C.tox_get_self_user_status(t.tox)

	return UserStatus(n), nil
}

// GetLastOnline returns the timestamp of the last time friendnumber was seen online.
func (t *Tox) GetLastOnline(friendnumber int32) (time.Time, error) {
	if t.tox == nil {
		return time.Time{}, ErrBadTox
	}

	ret := C.tox_get_last_online(t.tox, (C.int32_t)(friendnumber))

	if int(ret) == -1 {
		return time.Time{}, ErrFuncFail
	}

	last := time.Unix(int64(ret), 0)

	return last, nil
}

// SetUserIsTyping sets your typing status to a friend.
func (t *Tox) SetUserIsTyping(friendnumber int32, typing bool) error {
	if t.tox == nil {
		return ErrBadTox
	}
	ctyping := 0
	if typing {
		ctyping = 1
	}

	ret := C.tox_set_user_is_typing(t.tox, (C.int32_t)(friendnumber), (C.uint8_t)(ctyping))

	if ret != 0 {
		return ErrFuncFail
	}

	return nil
}

// GetIsTyping returns true if friendnumber is typing.
func (t *Tox) GetIsTyping(friendnumber int32) (bool, error) {
	if t.tox == nil {
		return false, ErrBadTox
	}

	ret := C.tox_get_is_typing(t.tox, (C.int32_t)(friendnumber))

	return (ret == 1), nil
}

// SetSendsReceipts sets whether we send read receipts to friendnumber.
func (t *Tox) SetSendsReceipts(friendnumber int32, send bool) error {
	if t.tox == nil {
		return ErrBadTox
	}
	csend := 0
	if send {
		csend = 1
	}

	C.tox_set_sends_receipts(t.tox, (C.int32_t)(friendnumber), (C.int)(csend))

	return nil
}

// CountFriendList returns the number of friends.
func (t *Tox) CountFriendlist() (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}
	n := C.tox_count_friendlist(t.tox)

	return uint32(n), nil
}

// GetNumOnlineFriends returns the number of online friends.
func (t *Tox) GetNumOnlineFriends() (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}
	n := C.tox_get_num_online_friends(t.tox)

	return uint32(n), nil
}

// GetFriendList returns a slice of int32 containing the friendnumbers.
func (t *Tox) GetFriendlist() ([]int32, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}

	size, _ := t.CountFriendlist()
	cfriendlist := make([]int32, size)

	n := C.tox_get_friendlist(t.tox, (*C.int32_t)(&cfriendlist[0]), (C.uint32_t)(size))

	friendlist := cfriendlist[:n]

	return friendlist, nil
}

// GetNoSpam returns the nospam of your ID.
func (t *Tox) GetNospam() (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	n := C.tox_get_nospam(t.tox)

	return uint32(n), nil
}

// SetNospam sets the nospam of your ID.
func (t *Tox) SetNospam(nospam uint32) error {
	if t.tox == nil {
		return ErrBadTox
	}

	C.tox_set_nospam(t.tox, (C.uint32_t)(nospam))

	return nil
}

// AddGroupchat creates a new groupchat and returns the groupnumber.
func (t *Tox) AddGroupchat() (int, error) {
	if t.tox == nil {
		return -1, ErrBadTox
	}

	ret := C.tox_add_groupchat(t.tox)

	if ret == -1 {
		return -1, ErrFuncFail
	}

	return int(ret), nil
}

// DelGroupchat deletes a groupchat identified by groupnumber.
func (t *Tox) DelGroupchat(groupnumber int) error {
	if t.tox == nil {
		return ErrBadTox
	}

	ret := C.tox_del_groupchat(t.tox, (C.int)(groupnumber))

	if ret == -1 {
		return ErrFuncFail
	}

	return nil
}

// GroupPeername returns the name of peernumber in groupnumber.
func (t *Tox) GroupPeername(groupnumber int, peernumber int) ([]byte, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}

	cname := make([]byte, MAX_NAME_LENGTH)
	ret := C.tox_group_peername(t.tox, (C.int)(groupnumber), (C.int)(peernumber), (*C.uint8_t)(&cname[0]))

	if ret == -1 {
		return nil, ErrFuncFail
	}

	name := cname[:ret]

	return name, nil
}

// InviteFriend invites friendnumber to groupnumber.
func (t *Tox) InviteFriend(friendnumber int32, groupnumber int) error {
	if t.tox == nil {
		return ErrBadTox
	}

	ret := C.tox_invite_friend(t.tox, (C.int32_t)(friendnumber), (C.int)(groupnumber))

	if ret == -1 {
		return ErrFuncFail
	}

	return nil
}

// JoinGroupchat joins the groupchat after having been invited by friendnumber.
func (t *Tox) JoinGroupchat(friendnumber int32, friendGroupPublicKey []byte) (int, error) {
	if t.tox == nil {
		return -1, ErrBadTox
	}

	if len(friendGroupPublicKey) == 0 {
		return -1, ErrArgs
	}

	ret := C.tox_join_groupchat(t.tox, (C.int32_t)(friendnumber), (*C.uint8_t)(&friendGroupPublicKey[0]))

	if ret == -1 {
		return -1, ErrFuncFail
	}

	return int(ret), nil
}

// GroupMessageSend sends a message to groupnumber.
func (t *Tox) GroupMessageSend(groupnumber int, message []byte) error {
	if t.tox == nil {
		return ErrBadTox
	}

	if len(message) == 0 {
		return ErrArgs
	}

	ret := C.tox_group_message_send(t.tox, (C.int)(groupnumber), (*C.uint8_t)(&message[0]), (C.uint32_t)(len(message)))

	if ret == -1 {
		return ErrFuncFail
	}

	return nil
}

// GroupActionSend sends an action to groupnumber.
func (t *Tox) GroupActionSend(groupnumber int, action []byte) error {
	if t.tox == nil {
		return ErrBadTox
	}

	if len(action) == 0 {
		return ErrArgs
	}

	ret := C.tox_group_action_send(t.tox, (C.int)(groupnumber), (*C.uint8_t)(&action[0]), (C.uint32_t)(len(action)))

	if ret == -1 {
		return ErrFuncFail
	}

	return nil
}

// GroupNumberPeers returns the number of peers in groupnumber.
func (t *Tox) GroupNumberPeers(groupnumber int) (int, error) {
	if t.tox == nil {
		return -1, ErrBadTox
	}

	ret := C.tox_group_number_peers(t.tox, (C.int)(groupnumber))

	if ret == -1 {
		return -1, ErrFuncFail
	}

	return int(ret), nil
}

//TODO
/* List all the peers in the group chat.
*
* Copies the names of the peers to the name[length][TOX_MAX_NAME_LENGTH] array.
*
* Copies the lengths of the names to lengths[length]
*
* returns the number of peers on success.
*
* return -1 on failure.
 */
//int tox_group_get_names(Tox *tox, int groupnumber, uint8_t names[][TOX_MAX_NAME_LENGTH], uint16_t lengths[],uint16_t length);

// CountChatlist returns the number of chats.
func (t *Tox) CountChatlist() (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}
	n := C.tox_count_friendlist(t.tox)

	return uint32(n), nil
}

// GetChatlist returns a slice of chat IDs.
func (t *Tox) GetChatlist() ([]int, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}

	size, _ := t.CountChatlist()
	cchatlist := make([]int, size)

	n := C.tox_get_chatlist(t.tox, (*C.int)(unsafe.Pointer(&cchatlist[0])), (C.uint32_t)(size))

	chatlist := cchatlist[:n]

	return chatlist, nil
}

// NewFileSender sends a file send request to friendnumber.
// Returns the filenumber on success.
func (t *Tox) NewFileSender(friendnumber int32, filesize uint64, filename []byte) (int, error) {
	if t.tox == nil {
		return -1, ErrBadTox
	}

	if len(filename) > 255 {
		return -1, ErrArgs
	}

	n := C.tox_new_file_sender(t.tox, (C.int32_t)(friendnumber), (C.uint64_t)(filesize), (*C.uint8_t)(&filename[0]), (C.uint16_t)(len(filename)))

	if n == -1 {
		return -1, ErrFuncFail
	}

	return int(n), nil
}

// FileSendControl sends a FileControl to friendnumber.
func (t *Tox) FileSendControl(friendnumber int32, receiving bool, filenumber uint8, messageId FileControl, data []byte) error {
	if t.tox == nil {
		return ErrBadTox
	}

	cReceiving := 0
	if receiving {
		cReceiving = 1
	}

	// Stupid workaround to prevent index out of range when using &data[0] if data == nil
	var cdata *C.uint8_t
	var clen C.uint16_t

	if data == nil {
		cdata = nil
		clen = 0
	} else {
		cdata = (*C.uint8_t)(&data[0])
		clen = (C.uint16_t)(len(data))
	}
	// End of stupid workaround

	n := C.tox_file_send_control(t.tox, (C.int32_t)(friendnumber), (C.uint8_t)(cReceiving), (C.uint8_t)(filenumber), (C.uint8_t)(messageId), cdata, clen)

	if n == -1 {
		return ErrFuncFail
	}

	return nil
}

// FileSendData sends file data of filenumber to friendnumber.
func (t *Tox) FileSendData(friendnumber int32, filenumber uint8, data []byte) error {
	if t.tox == nil {
		return ErrBadTox
	}

	if len(data) == 0 {
		return ErrArgs

	}

	n := C.tox_file_send_data(t.tox, (C.int32_t)(friendnumber), (C.uint8_t)(filenumber), (*C.uint8_t)(&data[0]), (C.uint16_t)(len(data)))

	if n == -1 {
		return ErrFuncFail
	}

	return nil
}

// FileDataSize returns the recommended/maximum size of the data chunks for FileSendData.
func (t *Tox) FileDataSize(friendnumber int32) (int, error) {
	if t.tox == nil {
		return -1, ErrBadTox
	}

	n := C.tox_file_data_size(t.tox, (C.int32_t)(friendnumber))

	if n == -1 {
		return -1, ErrFuncFail
	}

	return int(n), nil
}

// FileDataRemaining returns the number of bytes left to be transfered.
func (t *Tox) FileDataRemaining(friendnumber int32, filenumber uint8, receiving bool) (uint64, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	cReceiving := 0
	if receiving {
		cReceiving = 1
	}

	n := C.tox_file_data_remaining(t.tox, (C.int32_t)(friendnumber), (C.uint8_t)(filenumber), (C.uint8_t)(cReceiving))

	if n == 0 {
		return 0, ErrFuncFail
	}

	return uint64(n), nil
}

// Size returns the size of the save data returned by Save.
func (t *Tox) Size() (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	return uint32(C.tox_size(t.tox)), nil
}

// Save returns a byte slice of the save data.
func (t *Tox) Save() ([]byte, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}
	size, _ := t.Size()

	data := make([]byte, size)
	C.tox_save(t.tox, (*C.uint8_t)(&data[0]))

	return data, nil
}

// Load loads a save data.
func (t *Tox) Load(data []byte) error {
	if t.tox == nil {
		return ErrBadTox
	}

	ret := C.tox_load(t.tox, (*C.uint8_t)(&data[0]), (C.uint32_t)(len(data)))

	if ret == -1 {
		return ErrFuncFail
	}

	return nil
}

// CallbackFriendRequest sets the function to be called when receiving a friend request.
func (t *Tox) CallbackFriendRequest(f OnFriendRequest) {
	if t.tox != nil {
		t.onFriendRequest = f
		C.set_callback_friend_request(t.tox, unsafe.Pointer(t))
	}
}

// CallbackFriendMessage sets the function to be called when receiving a friend message.
func (t *Tox) CallbackFriendMessage(f OnFriendMessage) {
	if t.tox != nil {
		t.onFriendMessage = f
		C.set_callback_friend_message(t.tox, unsafe.Pointer(t))
	}
}

// CallbackFriendAction sets the function to be called when receiving a friend action.
func (t *Tox) CallbackFriendAction(f OnFriendAction) {
	if t.tox != nil {
		t.onFriendAction = f
		C.set_callback_friend_action(t.tox, unsafe.Pointer(t))
	}
}

// CallbackNameChange sets the callback for name changes.
func (t *Tox) CallbackNameChange(f OnNameChange) {
	if t.tox != nil {
		t.onNameChange = f
		C.set_callback_name_change(t.tox, unsafe.Pointer(t))
	}
}

// CallbackStatusMessage sets the callback for status message changes.
func (t *Tox) CallbackStatusMessage(f OnStatusMessage) {
	if t.tox != nil {
		t.onStatusMessage = f
		C.set_callback_status_message(t.tox, unsafe.Pointer(t))
	}
}

// CallbackUserStatus sets the callback for user status changes.
func (t *Tox) CallbackUserStatus(f OnUserStatus) {
	if t.tox != nil {
		t.onUserStatus = f
		C.set_callback_user_status(t.tox, unsafe.Pointer(t))
	}
}

// CallbackTypingChange sets the callback for typing changes.
func (t *Tox) CallbackTypingChange(f OnTypingChange) {
	if t.tox != nil {
		t.onTypingChange = f
		C.set_callback_typing_change(t.tox, unsafe.Pointer(t))
	}
}

// CallbackReadReceipt sets the function to be called when receiving read receipts.
func (t *Tox) CallbackReadReceipt(f OnReadReceipt) {
	if t.tox != nil {
		t.onReadReceipt = f
		C.set_callback_read_receipt(t.tox, unsafe.Pointer(t))
	}
}

// CallbackConnectionStatus sets the callback for connection status changes.
func (t *Tox) CallbackConnectionStatus(f OnConnectionStatus) {
	if t.tox != nil {
		t.onConnectionStatus = f
		C.set_callback_connection_status(t.tox, unsafe.Pointer(t))
	}
}

// CallbackFileSendRequest sets the callback for file send requests.
func (t *Tox) CallbackFileSendRequest(f OnFileSendRequest) {
	if t.tox != nil {
		t.onFileSendRequest = f
		C.set_callback_file_send_request(t.tox, unsafe.Pointer(t))
	}
}

// CallbackFileControl sets the callback for file control requests.
func (t *Tox) CallbackFileControl(f OnFileControl) {
	if t.tox != nil {
		t.onFileControl = f
		C.set_callback_file_control(t.tox, unsafe.Pointer(t))
	}
}

// CallbackFileData sets the callback for file data.
func (t *Tox) CallbackFileData(f OnFileData) {
	if t.tox != nil {
		t.onFileData = f
		C.set_callback_file_data(t.tox, unsafe.Pointer(t))
	}
}

// CallbackGroupInvite sets the callback for group invites.
func (t *Tox) CallbackGroupInvite(f OnGroupInvite) {
	if t.tox != nil {
		t.onGroupInvite = f
		C.set_callback_group_invite(t.tox, unsafe.Pointer(t))
	}
}

// CallbackGroupMessage sets the callback for group messages.
func (t *Tox) CallbackGroupMessage(f OnGroupMessage) {
	if t.tox != nil {
		t.onGroupMessage = f
		C.set_callback_group_message(t.tox, unsafe.Pointer(t))
	}
}

// CallbackGroupAction sets the callback for group actions.
func (t *Tox) CallbackGroupAction(f OnGroupAction) {
	if t.tox != nil {
		t.onGroupAction = f
		C.set_callback_group_action(t.tox, unsafe.Pointer(t))
	}
}

// CallbackGroupNamelistChange sets the callback for peer name list changes.
func (t *Tox) CallbackGroupNamelistChange(f OnGroupNamelistChange) {
	if t.tox != nil {
		t.onGroupNamelistChange = f
		C.set_callback_group_namelist_change(t.tox, unsafe.Pointer(t))
	}
}
