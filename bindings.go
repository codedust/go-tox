package gotox

//#include <tox/tox.h>
import "C"
import "time"

/* SelfGetConnectionStatus returns true if Tox is connected to the DHT. */
func (t *Tox) SelfGetConnectionStatus() (ConnectionStatus, error) {
	if t.tox == nil {
		return CONNECTION_NONE, ErrBadTox
	}

	return ConnectionStatus(C.tox_self_get_connection_status(t.tox)), nil
}

/* SelfGetAddress returns the public address to give to others. */
func (t *Tox) SelfGetAddress() ([]byte, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}

	address := make([]byte, ADDRESS_SIZE)
	C.tox_self_get_address(t.tox, (*C.uint8_t)(&address[0]))

	return address, nil
}

/* FriendAdd adds a friend by sending a friend request containing the given
 * message.
 * Returns the friend number on succes, or a FriendAddError on failure.
 */
func (t *Tox) FriendAdd(address []byte, message string) (int32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	if len(address) != ADDRESS_SIZE || len(message) == 0 {
		return 0, ErrArgs
	}

	caddr := (*C.uint8_t)(&address[0])
	cmessage := (*C.uint8_t)(&[]byte(message)[0])

	var friendAddError C.TOX_ERR_FRIEND_ADD
	ret := C.tox_friend_add(t.tox, caddr, cmessage, (C.size_t)(len(message)), &friendAddError)

	var err error

	switch FriendAddError(friendAddError) {
	case ERR_FRIEND_ADD_OK:
		err = nil
	case ERR_FRIEND_ADD_NULL:
		err = FaerrNull
	case ERR_FRIEND_ADD_TOO_LONG:
		err = FaerrTooLong
	case ERR_FRIEND_ADD_NO_MESSAGE:
		err = FaerrNoMessage
	case ERR_FRIEND_ADD_OWN_KEY:
		err = FaerrOwnKey
	case ERR_FRIEND_ADD_ALREADY_SENT:
		err = FaerrAlreadySent
	case ERR_FRIEND_ADD_BAD_CHECKSUM:
		err = FaerrBadChecksum
	case ERR_FRIEND_ADD_SET_NEW_NOSPAM:
		err = FaerrSetNewNospam
	case ERR_FRIEND_ADD_MALLOC:
		err = FaerrNoMem
	default:
		err = FaerrUnkown
	}

	return int32(ret), err
}

/* FriendAddNorequest adds a friend without sending a friend request.
 * Returns the friend number on success.
 */
func (t *Tox) FriendAddNorequest(publickey []byte) (int32, error) {
	if t.tox == nil {
		return -1, ErrBadTox
	}

	if len(publickey) != PUBLIC_KEY_SIZE {
		return -1, ErrArgs
	}

	var friendAddError C.TOX_ERR_FRIEND_ADD
	ret := C.tox_friend_add_norequest(t.tox, (*C.uint8_t)(&publickey[0]), &friendAddError)
	if ret == C.UINT32_MAX {
		return -1, ErrFuncFail
	}

	var err error

	switch FriendAddError(friendAddError) {
	case ERR_FRIEND_ADD_OK:
		err = nil
	case ERR_FRIEND_ADD_NULL:
		err = FaerrNull
	case ERR_FRIEND_ADD_TOO_LONG:
		err = FaerrTooLong
	case ERR_FRIEND_ADD_NO_MESSAGE:
		err = FaerrNoMessage
	case ERR_FRIEND_ADD_OWN_KEY:
		err = FaerrOwnKey
	case ERR_FRIEND_ADD_ALREADY_SENT:
		err = FaerrAlreadySent
	case ERR_FRIEND_ADD_BAD_CHECKSUM:
		err = FaerrBadChecksum
	case ERR_FRIEND_ADD_SET_NEW_NOSPAM:
		err = FaerrSetNewNospam
	case ERR_FRIEND_ADD_MALLOC:
		err = FaerrNoMem
	default:
		err = FaerrUnkown
	}

	return int32(ret), err
}

/* FriendGetNumber returns the friend number associated to a given publickey. */
func (t *Tox) FriendGetNumber(publickey []byte) (int32, error) {
	if t.tox == nil {
		return -1, ErrBadTox
	}

	if len(publickey) != PUBLIC_KEY_SIZE {
		return -1, ErrArgs
	}

	var friendByPublicKeyError C.TOX_ERR_FRIEND_BY_PUBLIC_KEY
	n := C.tox_friend_by_public_key(t.tox, (*C.uint8_t)(&publickey[0]), &friendByPublicKeyError)

	var err error

	switch FriendByPublicKeyError(friendByPublicKeyError) {
	case TOX_ERR_FRIEND_BY_PUBLIC_KEY_OK:
		err = nil
	case TOX_ERR_FRIEND_BY_PUBLIC_KEY_NULL:
		err = ErrArgs
	case TOX_ERR_FRIEND_BY_PUBLIC_KEY_NOT_FOUND:
		err = ErrFuncFail
	default:
		err = ErrUnknown
	}

	return int32(n), err
}

/* FriendGetPublickey returns the publickey associated to that friendnumber. */
func (t *Tox) FriendGetPublickey(friendnumber uint32) ([]byte, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}
	publickey := make([]byte, PUBLIC_KEY_SIZE)
	var friendGetPublicKeyError C.TOX_ERR_FRIEND_GET_PUBLIC_KEY = C.TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK
	ret := C.tox_friend_get_public_key(t.tox, (C.uint32_t)(friendnumber), (*C.uint8_t)(&publickey[0]), &friendGetPublicKeyError)

	var err error
	switch FriendGetPublicKeyError(friendGetPublicKeyError) {
	case ERR_FRIEND_GET_PUBLIC_KEY_OK:
		err = nil
	case ERR_FRIEND_GET_PUBLIC_KEY_FRIEND_NOT_FOUND:
		err = ErrFuncFail
	default:
		err = ErrUnknown
	}

	if ret != true {
		return nil, ErrFuncFail
	}

	return publickey, err
}

/* FriendDelete removes a friend. */
func (t *Tox) FriendDelete(friendnumber uint32) error {
	if t.tox == nil {
		return ErrBadTox
	}

	var friendDeleteError C.TOX_ERR_FRIEND_DELETE = C.TOX_ERR_FRIEND_DELETE_OK
	ret := C.tox_friend_delete(t.tox, (C.uint32_t)(friendnumber), &friendDeleteError)

	var err error
	switch FriendDeleteError(friendDeleteError) {
	case ERR_FRIEND_DELETE_OK:
		err = nil
	case ERR_FRIEND_DELETE_FRIEND_NOT_FOUND:
		err = ErrFuncFail
	default:
		err = ErrUnknown
	}

	if ret != true {
		return ErrFuncFail
	}

	return err
}

/* FriendGetConnectionStatus returns true if the friend is connected. */
func (t *Tox) FriendGetConnectionStatus(friendnumber uint32) (ConnectionStatus, error) {
	if t.tox == nil {
		return CONNECTION_NONE, ErrBadTox
	}

	var friendQueryError C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	status := C.tox_friend_get_connection_status(t.tox, (C.uint32_t)(friendnumber), &friendQueryError)

	if FriendQueryError(friendQueryError) != ERR_FRIEND_QUERY_OK {
		return CONNECTION_NONE, ErrFuncFail
	}

	return ConnectionStatus(status), nil
}

/* FriendExists returns true if a friend exists with given friendnumber. */
func (t *Tox) FriendExists(friendnumber uint32) (bool, error) {
	if t.tox == nil {
		return false, ErrBadTox
	}

	ret := C.tox_friend_exists(t.tox, (C.uint32_t)(friendnumber))

	return bool(ret), nil
}

/* FriendSendMessage sends a message to a friend if he/she is online.
 * Maximum message length is MAX_MESSAGE_LENGTH.
 * messagetype is the type of the message (normal, action, ...).
 * Returns the message ID if successful, an error otherwise.
 */
func (t *Tox) FriendSendMessage(friendnumber uint32, messagetype MessageType, message string) (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	if len(message) == 0 {
		return 0, ErrArgs
	}

	var cMessageType C.TOX_MESSAGE_TYPE
	if messagetype == MESSAGE_TYPE_ACTION {
		cMessageType = C.TOX_MESSAGE_TYPE_ACTION
	} else {
		cMessageType = C.TOX_MESSAGE_TYPE_NORMAL
	}

	cMessage := (*C.uint8_t)(&[]byte(message)[0])

	var friendSendMessageError C.TOX_ERR_FRIEND_SEND_MESSAGE = C.TOX_ERR_FRIEND_SEND_MESSAGE_OK
	n := C.tox_friend_send_message(t.tox, (C.uint32_t)(friendnumber), cMessageType, cMessage, (C.size_t)(len(message)), &friendSendMessageError)

	if FriendSendMessageError(friendSendMessageError) != ERR_FRIEND_SEND_MESSAGE_OK {
		return 0, ErrFuncFail
	}

	return uint32(n), nil
}

/* SelfSetName sets your nickname. The maximum name length is MAX_NAME_LENGTH. */
func (t *Tox) SelfSetName(name string) error {
	if t.tox == nil {
		return ErrBadTox
	}

	if len(name) == 0 {
		return ErrArgs
	}

	cName := (*C.uint8_t)(&[]byte(name)[0])

	var setInfoError C.TOX_ERR_SET_INFO = C.TOX_ERR_SET_INFO_OK
	success := C.tox_self_set_name(t.tox, cName, (C.size_t)(len(name)), &setInfoError)
	if !success || SetInfoError(setInfoError) != ERR_SET_INFO_OK {
		return ErrFuncFail
	}

	return nil
}

/* SelfGetNameSize returns the length of your name. */
func (t *Tox) SelfGetNameSize() (int, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	ret := C.tox_self_get_name_size(t.tox)

	return int(ret), nil
}

/* SelfGetName returns your nickname. */
func (t *Tox) SelfGetName() (string, error) {
	if t.tox == nil {
		return "", ErrBadTox
	}

	length, err := t.SelfGetNameSize()
	if err != nil {
		return "", ErrFuncFail
	}

	name := make([]byte, length)

	if length > 0 {
		C.tox_self_get_name(t.tox, (*C.uint8_t)(&name[0]))
	}

	return string(name), nil
}

/* FriendGetNameSize returns the length of the name of friendnumber. */
func (t *Tox) FriendGetNameSize(friendnumber uint32) (int, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	var friendQueryError C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	ret := C.tox_friend_get_name_size(t.tox, (C.uint32_t)(friendnumber), &friendQueryError)

	if FriendQueryError(friendQueryError) != ERR_FRIEND_QUERY_OK {
		return 0, ErrFuncFail
	}

	return int(ret), nil
}

/* FriendGetName returns the name of friendnumber. */
func (t *Tox) FriendGetName(friendnumber uint32) (string, error) {
	if t.tox == nil {
		return "", ErrBadTox
	}

	length, err := t.FriendGetNameSize(friendnumber)
	if err != nil {
		return "", ErrFuncFail
	}

	name := make([]byte, length)

	if length > 0 {
		var friendQueryError C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
		success := C.tox_friend_get_name(t.tox, (C.uint32_t)(friendnumber), (*C.uint8_t)(&name[0]), &friendQueryError)

		if success != true || FriendQueryError(friendQueryError) != ERR_FRIEND_QUERY_OK {
			return "", ErrFuncFail
		}
	}

	return string(name), nil
}

/* SelfSetStatusMessage sets your status message.
 * The maximum status length is MAX_STATUS_MESSAGE_LENGTH.
 */
func (t *Tox) SelfSetStatusMessage(status string) error {
	if t.tox == nil {
		return ErrBadTox
	}

	if len(status) == 0 {
		return ErrArgs
	}

	cStatus := (*C.uint8_t)(&[]byte(status)[0])

	var setInfoError C.TOX_ERR_SET_INFO = C.TOX_ERR_SET_INFO_OK
	C.tox_self_set_status_message(t.tox, cStatus, (C.size_t)(len(status)), &setInfoError)

	if SetInfoError(setInfoError) != ERR_SET_INFO_OK {
		return ErrFuncFail
	}

	return nil
}

/* SelfSetStatus sets your userstatus. */
func (t *Tox) SelfSetStatus(userstatus UserStatus) error {
	if t.tox == nil {
		return ErrBadTox
	}

	C.tox_self_set_status(t.tox, (C.TOX_USER_STATUS)(userstatus))

	return nil
}

/* SelfGetStatusMessageSize returns the size of your status message. */
func (t *Tox) SelfGetStatusMessageSize() (int, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	ret := C.tox_self_get_status_message_size(t.tox)

	return int(ret), nil
}

/* SelfGetStatusMessage returns your status message. */
func (t *Tox) SelfGetStatusMessage() (string, error) {
	if t.tox == nil {
		return "", ErrBadTox
	}

	length, err := t.SelfGetStatusMessageSize()
	if err != nil {
		return "", ErrFuncFail
	}

	statusMessage := make([]byte, length)

	if length > 0 {
		C.tox_self_get_status_message(t.tox, (*C.uint8_t)(&statusMessage[0]))
	}

	return string(statusMessage), nil
}

/* FriendGetStatusMessageSize returns the size of the status of a friend with
 * the given friendnumber.
 */
func (t *Tox) FriendGetStatusMessageSize(friendnumber uint32) (int, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	var friendQueryError C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	ret := C.tox_friend_get_status_message_size(t.tox, (C.uint32_t)(friendnumber), &friendQueryError)

	if FriendQueryError(friendQueryError) != ERR_FRIEND_QUERY_OK {
		return 0, ErrFuncFail
	}

	return int(ret), nil
}

/* FriendGetStatusMessage returns the status message of friend with the given
 * friendnumber.
 */
func (t *Tox) FriendGetStatusMessage(friendnumber uint32) (string, error) {
	if t.tox == nil {
		return "", ErrBadTox
	}

	var friendQueryError C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK

	size, error := t.FriendGetStatusMessageSize(friendnumber)
	if error != nil {
		return "", ErrFuncFail
	}

	statusMessage := make([]byte, size)

	if size > 0 {
		friendQueryError = C.TOX_ERR_FRIEND_QUERY_OK
		n := C.tox_friend_get_status_message(t.tox, (C.uint32_t)(friendnumber), (*C.uint8_t)(&statusMessage[0]), &friendQueryError)

		if n != true || FriendQueryError(friendQueryError) != ERR_FRIEND_QUERY_OK {
			return "", ErrFuncFail
		}
	}

	return string(statusMessage), nil
}

/* FriendGetStatus returns the status of friendnumber. */
func (t *Tox) FriendGetStatus(friendnumber uint32) (UserStatus, error) {
	if t.tox == nil {
		return USERSTATUS_NONE, ErrBadTox
	}

	var friendQueryError C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	status := C.tox_friend_get_status(t.tox, (C.uint32_t)(friendnumber), &friendQueryError)

	if FriendQueryError(friendQueryError) != ERR_FRIEND_QUERY_OK {
		return USERSTATUS_NONE, ErrFuncFail
	}

	return UserStatus(status), nil
}

/* SelfGetStatus returns your status. */
func (t *Tox) SelfGetStatus() (UserStatus, error) {
	if t.tox == nil {
		return USERSTATUS_NONE, ErrBadTox
	}

	n := C.tox_self_get_status(t.tox)

	return UserStatus(n), nil
}

/* FriendGetLastOnline returns the timestamp of the last time the friend with
 * the given friendnumber was seen online.
 */
func (t *Tox) FriendGetLastOnline(friendnumber uint32) (time.Time, error) {
	if t.tox == nil {
		return time.Time{}, ErrBadTox
	}

	var friendGetLastOnlineError C.TOX_ERR_FRIEND_GET_LAST_ONLINE = C.TOX_ERR_FRIEND_GET_LAST_ONLINE_OK
	ret := C.tox_friend_get_last_online(t.tox, (C.uint32_t)(friendnumber), &friendGetLastOnlineError)

	if int(ret) == -1 || FriendGetLastOnlineError(friendGetLastOnlineError) != ERR_FRIEND_GET_LAST_ONLINE_OK {
		return time.Time{}, ErrFuncFail
	}

	last := time.Unix(int64(ret), 0)

	return last, nil
}

/* SelfSetTyping sets your typing status to a friend. */
func (t *Tox) SelfSetTyping(friendnumber uint32, typing bool) error {
	if t.tox == nil {
		return ErrBadTox
	}

	var setTypingError C.TOX_ERR_SET_TYPING = C.TOX_ERR_SET_TYPING_OK
	success := C.tox_self_set_typing(t.tox, (C.uint32_t)(friendnumber), (C._Bool)(typing), &setTypingError)

	if !success || SetTypingError(setTypingError) != ERR_SET_TYPING_OK {
		return ErrFuncFail
	}

	return nil
}

/* FriendGetTyping returns true if friendnumber is typing. */
func (t *Tox) FriendGetTyping(friendnumber uint32) (bool, error) {
	if t.tox == nil {
		return false, ErrBadTox
	}

	var friendQueryError C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	istyping := C.tox_friend_get_typing(t.tox, (C.uint32_t)(friendnumber), &friendQueryError)

	if FriendQueryError(friendQueryError) != ERR_FRIEND_QUERY_OK {
		return false, ErrFuncFail
	}

	return bool(istyping), nil
}

/* SelfGetFriendlistSize returns the number of friends on the friendlist. */
func (t *Tox) SelfGetFriendlistSize() (int32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}
	n := C.tox_self_get_friend_list_size(t.tox)

	return int32(n), nil
}

/* SelfGetFriendlist returns a slice of uint32 containing the friendnumbers. */
func (t *Tox) SelfGetFriendlist() ([]uint32, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}

	size, err := t.SelfGetFriendlistSize()
	if err != nil {
		return nil, ErrFuncFail
	}

	friendlist := make([]uint32, size)

	if size > 0 {
		C.tox_self_get_friend_list(t.tox, (*C.uint32_t)(&friendlist[0]))
	}

	return friendlist, nil
}

/* SelfGetNospam returns the nospam of your ID. */
func (t *Tox) SelfGetNospam() (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	n := C.tox_self_get_nospam(t.tox)
	return uint32(n), nil
}

/* SelfSetNospam sets the nospam of your ID. */
func (t *Tox) SelfSetNospam(nospam uint32) error {
	if t.tox == nil {
		return ErrBadTox
	}

	C.tox_self_set_nospam(t.tox, (C.uint32_t)(nospam))
	return nil
}

/* SendFileControl sends a FileControl to a friend with the given friendnumber. */
func (t *Tox) SendFileControl(friendnumber uint32, receiving bool, filenumber uint32, fileControl FileControl, data []byte) error {
	if t.tox == nil {
		return ErrBadTox
	}

	var cFileControl C.TOX_FILE_CONTROL
	switch FileControl(fileControl) {
	case FILE_CONTROL_RESUME:
		cFileControl = C.TOX_FILE_CONTROL_RESUME
	case FILE_CONTROL_PAUSE:
		cFileControl = C.TOX_FILE_CONTROL_PAUSE
	case FILE_CONTROL_CANCEL:
		cFileControl = C.TOX_FILE_CONTROL_CANCEL
	}

	var fileControlError C.TOX_ERR_FILE_CONTROL
	success := C.tox_file_control(t.tox, (C.uint32_t)(friendnumber), (C.uint32_t)(filenumber), cFileControl, &fileControlError)

	if !success || FileControlError(fileControlError) != ERR_FILE_CONTROL_OK {
		return ErrFuncFail
	}

	return nil
}
