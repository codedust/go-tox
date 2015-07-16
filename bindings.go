package gotox

//#include <tox/tox.h>
//#include <stdlib.h>
import "C"
import "time"
import "unsafe"

/* New returns a new Tox instance. */
func New(options *Options) (*Tox, error) {
	var ctox *C.Tox
	var toxErrNew C.TOX_ERR_NEW

	if options != nil {
		// Let's map options from Options to C.Tox_Options
		var cSaveData *C.uint8_t
		if len(options.SaveData) > 0 {
			cSaveData = (*C.uint8_t)(&options.SaveData[0])
		} else {
			cSaveData = nil
		}
		cIPv6Enabled := (C._Bool)(options.IPv6Enabled)
		cUDPEnabled := (C._Bool)(options.UDPEnabled)

		var cProxyType C.TOX_PROXY_TYPE = C.TOX_PROXY_TYPE_NONE
		if options.ProxyType == TOX_PROXY_TYPE_HTTP {
			cProxyType = C.TOX_PROXY_TYPE_HTTP
		} else if options.ProxyType == TOX_PROXY_TYPE_SOCKS5 {
			cProxyType = C.TOX_PROXY_TYPE_SOCKS5
		}

		var cSaveDataType C.TOX_SAVEDATA_TYPE = C.TOX_SAVEDATA_TYPE_NONE
		if options.SaveDataType == TOX_SAVEDATA_TYPE_TOX_SAVE {
			cSaveDataType = C.TOX_SAVEDATA_TYPE_TOX_SAVE
		} else if options.SaveDataType == TOX_SAVEDATA_TYPE_SECRET_KEY {
			cSaveDataType = C.TOX_SAVEDATA_TYPE_SECRET_KEY
		}

		// max ProxyHost length is 255
		if len(options.ProxyHost) > 255 {
			return nil, ErrArgs
		}
		cProxyHost := C.CString(options.ProxyHost)
		defer C.free(unsafe.Pointer(cProxyHost))

		cProxyPort := (C.uint16_t)(options.ProxyPort)
		cStartPort := (C.uint16_t)(options.StartPort)
		cEndPort := (C.uint16_t)(options.EndPort)

		cOptions := &C.struct_Tox_Options{
			ipv6_enabled:    cIPv6Enabled,
			udp_enabled:     cUDPEnabled,
			proxy_type:      cProxyType,
			proxy_host:      cProxyHost,
			proxy_port:      cProxyPort,
			start_port:      cStartPort,
			end_port:        cEndPort,
			tcp_port:        0,
			savedata_type:   cSaveDataType,
			savedata_data:   cSaveData,
			savedata_length: (C.size_t)(len(options.SaveData))}

		ctox = C.tox_new(cOptions, &toxErrNew)
	} else {
		ctox = C.tox_new(nil, &toxErrNew)
	}

	if ctox == nil || ToxErrNew(toxErrNew) != TOX_ERR_NEW_OK {
		return nil, ErrInit
	}

	t := &Tox{tox: ctox}

	return t, nil
}

/* Kill stops a Tox instance. */
func (t *Tox) Kill() error {
	if t.tox == nil {
		return ErrBadTox
	}
	C.tox_kill(t.tox)

	return nil
}

/* GetSaveDataSize returns the size of the save data returned by GetSavedata. */
func (t *Tox) GetSaveDataSize() (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	return uint32(C.tox_get_savedata_size(t.tox)), nil
}

/* GetSavedata returns a byte slice of all information associated with the tox
 * instance. */
func (t *Tox) GetSavedata() ([]byte, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}
	size, err := t.GetSaveDataSize()
	if err != nil {
		return nil, ErrFuncFail
	}

	data := make([]byte, size)

	if size > 0 {
		C.tox_get_savedata(t.tox, (*C.uint8_t)(&data[0]))
	}

	return data, nil
}

/* Bootstrap sends a "get nodes" request to the given bootstrap node with IP,
 * port, and public key to setup connections. */
func (t *Tox) Bootstrap(address string, port uint16, publickey []byte) error {
	if t.tox == nil {
		return ErrBadTox
	}

	if len(publickey) != TOX_PUBLIC_KEY_SIZE {
		return ErrArgs
	}

	caddr := C.CString(address)
	defer C.free(unsafe.Pointer(caddr))

	var toxErrBootstrap C.TOX_ERR_BOOTSTRAP
	C.tox_bootstrap(t.tox, caddr, (C.uint16_t)(port), (*C.uint8_t)(&publickey[0]), &toxErrBootstrap)

	var bootstrapError error

	switch ToxErrBootstrap(toxErrBootstrap) {
	case TOX_ERR_BOOTSTRAP_OK:
		bootstrapError = nil
	case TOX_ERR_BOOTSTRAP_NULL:
		bootstrapError = ErrArgs
	case TOX_ERR_BOOTSTRAP_BAD_HOST:
		bootstrapError = ErrFuncFail
	case TOX_ERR_BOOTSTRAP_BAD_PORT:
		bootstrapError = ErrFuncFail
	}

	return bootstrapError
}

/* SelfGetConnectionStatus returns true if Tox is connected to the DHT. */
func (t *Tox) SelfGetConnectionStatus() (ToxConnection, error) {
	if t.tox == nil {
		return TOX_CONNECTION_NONE, ErrBadTox
	}

	return ToxConnection(C.tox_self_get_connection_status(t.tox)), nil
}

/* IterationInterval returns the time in milliseconds before Iterate() should be
 * called again. */
func (t *Tox) IterationInterval() (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	ret := C.tox_iteration_interval(t.tox)

	return uint32(ret), nil
}

/* Iterate is the main loop. It needs to be called every IterationInterval()
 * milliseconds. */
func (t *Tox) Iterate() error {
	if t.tox == nil {
		return ErrBadTox
	}

	t.mtx.Lock()
	C.tox_iterate(t.tox)
	t.mtx.Unlock()

	return nil
}

/* SelfGetAddress returns the public address to give to others. */
func (t *Tox) SelfGetAddress() ([]byte, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}

	address := make([]byte, TOX_ADDRESS_SIZE)
	C.tox_self_get_address(t.tox, (*C.uint8_t)(&address[0]))

	return address, nil
}

/* SelfSetNospam sets the nospam of your ID. */
func (t *Tox) SelfSetNospam(nospam uint32) error {
	if t.tox == nil {
		return ErrBadTox
	}

	C.tox_self_set_nospam(t.tox, (C.uint32_t)(nospam))
	return nil
}

/* SelfGetNospam returns the nospam of your ID. */
func (t *Tox) SelfGetNospam() (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	n := C.tox_self_get_nospam(t.tox)
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
	if !bool(success) || ToxErrSetInfo(setInfoError) != TOX_ERR_SET_INFO_OK {
		return ErrFuncFail
	}

	return nil
}

/* SelfGetNameSize returns the length of your name. */
func (t *Tox) SelfGetNameSize() (int64, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	ret := C.tox_self_get_name_size(t.tox)

	return int64(ret), nil
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

/* SelfSetStatusMessage sets your status message.
 * The maximum status length is MAX_STATUS_MESSAGE_LENGTH. */
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

	if ToxErrSetInfo(setInfoError) != TOX_ERR_SET_INFO_OK {
		return ErrFuncFail
	}

	return nil
}

/* SelfGetStatusMessageSize returns the size of your status message. */
func (t *Tox) SelfGetStatusMessageSize() (int64, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	ret := C.tox_self_get_status_message_size(t.tox)

	return int64(ret), nil
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

/* SelfSetStatus sets your userstatus. */
func (t *Tox) SelfSetStatus(userstatus ToxUserStatus) error {
	if t.tox == nil {
		return ErrBadTox
	}

	C.tox_self_set_status(t.tox, (C.TOX_USER_STATUS)(userstatus))

	return nil
}

/* SelfGetStatus returns your status. */
func (t *Tox) SelfGetStatus() (ToxUserStatus, error) {
	if t.tox == nil {
		return TOX_USERSTATUS_NONE, ErrBadTox
	}

	n := C.tox_self_get_status(t.tox)

	return ToxUserStatus(n), nil
}

/* FriendAdd adds a friend by sending a friend request containing the given
 * message.
 * Returns the friend number on succes, or a ToxErrFriendAdd on failure.
 */
func (t *Tox) FriendAdd(address []byte, message string) (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	if len(address) != TOX_ADDRESS_SIZE || len(message) == 0 {
		return 0, ErrArgs
	}

	caddr := (*C.uint8_t)(&address[0])
	cmessage := (*C.uint8_t)(&[]byte(message)[0])

	var toxErrFriendAdd C.TOX_ERR_FRIEND_ADD
	ret := C.tox_friend_add(t.tox, caddr, cmessage, (C.size_t)(len(message)), &toxErrFriendAdd)

	var err error

	switch ToxErrFriendAdd(toxErrFriendAdd) {
	case TOX_ERR_FRIEND_ADD_OK:
		err = nil
	case TOX_ERR_FRIEND_ADD_NULL:
		err = ErrFriendAddNull
	case TOX_ERR_FRIEND_ADD_TOO_LONG:
		err = ErrFriendAddTooLong
	case TOX_ERR_FRIEND_ADD_NO_MESSAGE:
		err = ErrFriendAddNoMessage
	case TOX_ERR_FRIEND_ADD_OWN_KEY:
		err = ErrFriendAddOwnKey
	case TOX_ERR_FRIEND_ADD_ALREADY_SENT:
		err = ErrFriendAddAlreadySent
	case TOX_ERR_FRIEND_ADD_BAD_CHECKSUM:
		err = ErrFriendAddBadChecksum
	case TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM:
		err = ErrFriendAddSetNewNospam
	case TOX_ERR_FRIEND_ADD_MALLOC:
		err = ErrFriendAddNoMem
	default:
		err = ErrFriendAddUnkown
	}

	return uint32(ret), err
}

/* FriendAddNorequest adds a friend without sending a friend request.
 * Returns the friend number on success.
 */
func (t *Tox) FriendAddNorequest(publickey []byte) (uint32, error) {
	if t.tox == nil {
		return C.UINT32_MAX, ErrBadTox
	}

	if len(publickey) != TOX_PUBLIC_KEY_SIZE {
		return C.UINT32_MAX, ErrArgs
	}

	var toxErrFriendAdd C.TOX_ERR_FRIEND_ADD
	ret := C.tox_friend_add_norequest(t.tox, (*C.uint8_t)(&publickey[0]), &toxErrFriendAdd)
	if ret == C.UINT32_MAX {
		return C.UINT32_MAX, ErrFuncFail
	}

	var err error

	switch ToxErrFriendAdd(toxErrFriendAdd) {
	case TOX_ERR_FRIEND_ADD_OK:
		err = nil
	case TOX_ERR_FRIEND_ADD_NULL:
		err = ErrFriendAddNull
	case TOX_ERR_FRIEND_ADD_TOO_LONG:
		err = ErrFriendAddTooLong
	case TOX_ERR_FRIEND_ADD_NO_MESSAGE:
		err = ErrFriendAddNoMessage
	case TOX_ERR_FRIEND_ADD_OWN_KEY:
		err = ErrFriendAddOwnKey
	case TOX_ERR_FRIEND_ADD_ALREADY_SENT:
		err = ErrFriendAddAlreadySent
	case TOX_ERR_FRIEND_ADD_BAD_CHECKSUM:
		err = ErrFriendAddBadChecksum
	case TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM:
		err = ErrFriendAddSetNewNospam
	case TOX_ERR_FRIEND_ADD_MALLOC:
		err = ErrFriendAddNoMem
	default:
		err = ErrFriendAddUnkown
	}

	return uint32(ret), err
}

/* FriendDelete removes a friend. */
func (t *Tox) FriendDelete(friendnumber uint32) error {
	if t.tox == nil {
		return ErrBadTox
	}

	var toxErrFriendDelete C.TOX_ERR_FRIEND_DELETE = C.TOX_ERR_FRIEND_DELETE_OK
	ret := C.tox_friend_delete(t.tox, (C.uint32_t)(friendnumber), &toxErrFriendDelete)

	var err error
	switch ToxErrFriendDelete(toxErrFriendDelete) {
	case TOX_ERR_FRIEND_DELETE_OK:
		err = nil
	case TOX_ERR_FRIEND_DELETE_FRIEND_NOT_FOUND:
		err = ErrFuncFail
	default:
		err = ErrUnknown
	}

	if ret != true {
		return ErrFuncFail
	}

	return err
}

/* FriendByPublicKey returns the friend number associated to a given publickey. */
func (t *Tox) FriendByPublicKey(publickey []byte) (uint32, error) {
	if t.tox == nil {
		return C.UINT32_MAX, ErrBadTox
	}

	if len(publickey) != TOX_PUBLIC_KEY_SIZE {
		return C.UINT32_MAX, ErrArgs
	}

	var toxErrFriendByPublicKey C.TOX_ERR_FRIEND_BY_PUBLIC_KEY
	n := C.tox_friend_by_public_key(t.tox, (*C.uint8_t)(&publickey[0]), &toxErrFriendByPublicKey)

	var err error

	switch ToxErrFriendByPublicKey(toxErrFriendByPublicKey) {
	case TOX_ERR_FRIEND_BY_PUBLIC_KEY_OK:
		err = nil
	case TOX_ERR_FRIEND_BY_PUBLIC_KEY_NULL:
		err = ErrArgs
	case TOX_ERR_FRIEND_BY_PUBLIC_KEY_NOT_FOUND:
		err = ErrFuncFail
	default:
		err = ErrUnknown
	}

	return uint32(n), err
}

/* FriendExists returns true if a friend exists with given friendnumber. */
func (t *Tox) FriendExists(friendnumber uint32) (bool, error) {
	if t.tox == nil {
		return false, ErrBadTox
	}

	ret := C.tox_friend_exists(t.tox, (C.uint32_t)(friendnumber))

	return bool(ret), nil
}

/* SelfGetFriendlistSize returns the number of friends on the friendlist. */
func (t *Tox) SelfGetFriendlistSize() (int64, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}
	n := C.tox_self_get_friend_list_size(t.tox)

	return int64(n), nil
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

/* FriendGetPublickey returns the publickey associated to that friendnumber. */
func (t *Tox) FriendGetPublickey(friendnumber uint32) ([]byte, error) {
	if t.tox == nil {
		return nil, ErrBadTox
	}
	publickey := make([]byte, TOX_PUBLIC_KEY_SIZE)
	var toxErrFriendGetPublicKey C.TOX_ERR_FRIEND_GET_PUBLIC_KEY = C.TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK
	ret := C.tox_friend_get_public_key(t.tox, (C.uint32_t)(friendnumber), (*C.uint8_t)(&publickey[0]), &toxErrFriendGetPublicKey)

	var err error
	switch ToxErrFriendGetPublicKey(toxErrFriendGetPublicKey) {
	case TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK:
		err = nil
	case TOX_ERR_FRIEND_GET_PUBLIC_KEY_FRIEND_NOT_FOUND:
		err = ErrFuncFail
	default:
		err = ErrUnknown
	}

	if ret != true {
		return nil, ErrFuncFail
	}

	return publickey, err
}

/* FriendGetLastOnline returns the timestamp of the last time the friend with
 * the given friendnumber was seen online. */
func (t *Tox) FriendGetLastOnline(friendnumber uint32) (time.Time, error) {
	if t.tox == nil {
		return time.Time{}, ErrBadTox
	}

	var toxErrFriendGetLastOnline C.TOX_ERR_FRIEND_GET_LAST_ONLINE = C.TOX_ERR_FRIEND_GET_LAST_ONLINE_OK
	ret := C.tox_friend_get_last_online(t.tox, (C.uint32_t)(friendnumber), &toxErrFriendGetLastOnline)

	if ret == C.INT64_MAX || ToxErrFriendGetLastOnline(toxErrFriendGetLastOnline) != TOX_ERR_FRIEND_GET_LAST_ONLINE_OK {
		return time.Time{}, ErrFuncFail
	}

	last := time.Unix(int64(ret), 0)

	return last, nil
}

/* FriendGetNameSize returns the length of the name of friendnumber. */
func (t *Tox) FriendGetNameSize(friendnumber uint32) (int64, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	ret := C.tox_friend_get_name_size(t.tox, (C.uint32_t)(friendnumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return 0, ErrFuncFail
	}

	return int64(ret), nil
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
		var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
		success := C.tox_friend_get_name(t.tox, (C.uint32_t)(friendnumber), (*C.uint8_t)(&name[0]), &toxErrFriendQuery)

		if success != true || ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
			return "", ErrFuncFail
		}
	}

	return string(name), nil
}

/* FriendGetStatusMessageSize returns the size of the status of a friend with
 * the given friendnumber.
 */
func (t *Tox) FriendGetStatusMessageSize(friendnumber uint32) (int64, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	ret := C.tox_friend_get_status_message_size(t.tox, (C.uint32_t)(friendnumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return 0, ErrFuncFail
	}

	return int64(ret), nil
}

/* FriendGetStatusMessage returns the status message of friend with the given
 * friendnumber.
 */
func (t *Tox) FriendGetStatusMessage(friendnumber uint32) (string, error) {
	if t.tox == nil {
		return "", ErrBadTox
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK

	size, error := t.FriendGetStatusMessageSize(friendnumber)
	if error != nil {
		return "", ErrFuncFail
	}

	statusMessage := make([]byte, size)

	if size > 0 {
		toxErrFriendQuery = C.TOX_ERR_FRIEND_QUERY_OK
		n := C.tox_friend_get_status_message(t.tox, (C.uint32_t)(friendnumber), (*C.uint8_t)(&statusMessage[0]), &toxErrFriendQuery)

		if n != true || ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
			return "", ErrFuncFail
		}
	}

	return string(statusMessage), nil
}

/* FriendGetStatus returns the status of friendnumber. */
func (t *Tox) FriendGetStatus(friendnumber uint32) (ToxUserStatus, error) {
	if t.tox == nil {
		return TOX_USERSTATUS_NONE, ErrBadTox
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	status := C.tox_friend_get_status(t.tox, (C.uint32_t)(friendnumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return TOX_USERSTATUS_NONE, ErrFuncFail
	}

	return ToxUserStatus(status), nil
}

/* FriendGetConnectionStatus returns true if the friend is connected. */
func (t *Tox) FriendGetConnectionStatus(friendnumber uint32) (ToxConnection, error) {
	if t.tox == nil {
		return TOX_CONNECTION_NONE, ErrBadTox
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	status := C.tox_friend_get_connection_status(t.tox, (C.uint32_t)(friendnumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return TOX_CONNECTION_NONE, ErrFuncFail
	}

	return ToxConnection(status), nil
}

/* FriendGetTyping returns true if friendnumber is typing. */
func (t *Tox) FriendGetTyping(friendnumber uint32) (bool, error) {
	if t.tox == nil {
		return false, ErrBadTox
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	istyping := C.tox_friend_get_typing(t.tox, (C.uint32_t)(friendnumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return false, ErrFuncFail
	}

	return bool(istyping), nil
}

/* SelfSetTyping sets your typing status to a friend. */
func (t *Tox) SelfSetTyping(friendnumber uint32, typing bool) error {
	if t.tox == nil {
		return ErrBadTox
	}

	var toxErrSetTyping C.TOX_ERR_SET_TYPING = C.TOX_ERR_SET_TYPING_OK
	success := C.tox_self_set_typing(t.tox, (C.uint32_t)(friendnumber), (C._Bool)(typing), &toxErrSetTyping)

	if !bool(success) || ToxErrSetTyping(toxErrSetTyping) != TOX_ERR_SET_TYPING_OK {
		return ErrFuncFail
	}

	return nil
}

/* FriendSendMessage sends a message to a friend if he/she is online.
 * Maximum message length is MAX_MESSAGE_LENGTH.
 * messagetype is the type of the message (normal, action, ...).
 * Returns the message ID if successful, an error otherwise.
 */
func (t *Tox) FriendSendMessage(friendnumber uint32, messagetype ToxMessageType, message string) (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	if len(message) == 0 {
		return 0, ErrArgs
	}

	var cMessageType C.TOX_MESSAGE_TYPE
	if messagetype == TOX_MESSAGE_TYPE_ACTION {
		cMessageType = C.TOX_MESSAGE_TYPE_ACTION
	} else {
		cMessageType = C.TOX_MESSAGE_TYPE_NORMAL
	}

	cMessage := (*C.uint8_t)(&[]byte(message)[0])

	var toxFriendSendMessageError C.TOX_ERR_FRIEND_SEND_MESSAGE = C.TOX_ERR_FRIEND_SEND_MESSAGE_OK
	n := C.tox_friend_send_message(t.tox, (C.uint32_t)(friendnumber), cMessageType, cMessage, (C.size_t)(len(message)), &toxFriendSendMessageError)

	if ToxErrFriendSendMessage(toxFriendSendMessageError) != TOX_ERR_FRIEND_SEND_MESSAGE_OK {
		return 0, ErrFuncFail
	}

	return uint32(n), nil
}

/* FileControl sends a FileControl to a friend with the given friendnumber. */
func (t *Tox) FileControl(friendnumber uint32, receiving bool, filenumber uint32, fileControl ToxFileControl, data []byte) error {
	if t.tox == nil {
		return ErrBadTox
	}

	var cFileControl C.TOX_FILE_CONTROL
	switch ToxFileControl(fileControl) {
	case TOX_FILE_CONTROL_RESUME:
		cFileControl = C.TOX_FILE_CONTROL_RESUME
	case TOX_FILE_CONTROL_PAUSE:
		cFileControl = C.TOX_FILE_CONTROL_PAUSE
	case TOX_FILE_CONTROL_CANCEL:
		cFileControl = C.TOX_FILE_CONTROL_CANCEL
	}

	var toxErrFileControl C.TOX_ERR_FILE_CONTROL
	success := C.tox_file_control(t.tox, (C.uint32_t)(friendnumber), (C.uint32_t)(filenumber), cFileControl, &toxErrFileControl)

	if !bool(success) || ToxErrFileControl(toxErrFileControl) != TOX_ERR_FILE_CONTROL_OK {
		return ErrFuncFail
	}

	return nil
}

/* FileSend sends a file transmission request. */
func (t *Tox) FileSend(friendnumber uint32, fileKind ToxFileKind, fileLength uint64, fileID []byte, fileName string) (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	var cFileKind = C.TOX_FILE_KIND_DATA
	switch ToxFileKind(fileKind) {
	case TOX_FILE_KIND_AVATAR:
		cFileKind = C.TOX_FILE_KIND_AVATAR
	case TOX_FILE_KIND_DATA:
		cFileKind = C.TOX_FILE_KIND_DATA
	}

	if fileID != nil && len(fileID) != TOX_FILE_ID_LENGTH {
		return 0, ErrFileSendInvalidFileID
	}

	var cFileID *C.uint8_t

	if fileID == nil {
		cFileID = nil
	} else {
		cFileID = (*C.uint8_t)(&[]byte(fileID)[0])
	}

	cFileName := (*C.uint8_t)(&[]byte(fileName)[0])

	var toxErrFileSend C.TOX_ERR_FILE_SEND
	n := C.tox_file_send(t.tox, (C.uint32_t)(friendnumber), (C.uint32_t)(cFileKind), (C.uint64_t)(fileLength), cFileID, cFileName, (C.size_t)(len(fileName)), &toxErrFileSend)

	if n == C.UINT32_MAX || ToxErrFileSend(toxErrFileSend) != TOX_ERR_FILE_SEND_OK {
		return 0, ErrFuncFail
	}
	return uint32(n), nil
}

/* FileSendChunk sends a chunk of file data to a friend. */
func (t *Tox) FileSendChunk(friendnumber uint32, fileNumber uint32, position uint64, data []byte) error {
	if t.tox == nil {
		return ErrBadTox
	}

	var cData *C.uint8_t

	if len(data) == 0 {
		cData = nil
	} else {
		cData = (*C.uint8_t)(&data[0])
	}

	var toxErrFileSendChunk C.TOX_ERR_FILE_SEND_CHUNK
	success := C.tox_file_send_chunk(t.tox, (C.uint32_t)(friendnumber), (C.uint32_t)(fileNumber), (C.uint64_t)(position), cData, (C.size_t)(len(data)), &toxErrFileSendChunk)

	if !bool(success) || ToxErrFileSendChunk(toxErrFileSendChunk) != TOX_ERR_FILE_SEND_CHUNK_OK {
		return ErrFuncFail
	}
	return nil
}
