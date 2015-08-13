package gotox

//#include <tox/tox.h>
//#include <stdlib.h>
import "C"
import "time"
import "unsafe"

/* VersionMajor returns the major version number of the used Tox library */
func VersionMajor() uint32 {
	return uint32(C.tox_version_major())
}

/* VersionMinor returns the minor version number of the used Tox library */
func VersionMinor() uint32 {
	return uint32(C.tox_version_minor())
}

/* VersionPatch returns the patch number of the used Tox library */
func VersionPatch() uint32 {
	return uint32(C.tox_version_patch())
}

/* VersionIsCompatible returns whether the compiled Tox library version is
 * compatible with the passed version numbers. */
func VersionIsCompatible(major uint32, minor uint32, patch uint32) bool {
	return bool(C.tox_version_is_compatible((C.uint32_t)(major), (C.uint32_t)(minor), (C.uint32_t)(patch)))
}

/* New creates and initialises a new Tox instance and returns the corresponding
 * gotox instance. */
func New(options *Options) (*Tox, error) {
	var cTox *C.Tox
	var toxErrNew C.TOX_ERR_NEW
	var toxErrOptionsNew C.TOX_ERR_OPTIONS_NEW

	var cOptions *C.struct_Tox_Options = C.tox_options_new(&toxErrOptionsNew)
	if cOptions == nil || ToxErrOptionsNew(toxErrOptionsNew) != TOX_ERR_OPTIONS_NEW_OK {
		return nil, ErrFuncFail
	}

	if options == nil {
		cOptions = nil
	} else {
		// map options from Options to C.Tox_Options
		cOptions.ipv6_enabled = C.bool(options.IPv6Enabled)
		cOptions.udp_enabled = C.bool(options.UDPEnabled)

		var cProxyType C.TOX_PROXY_TYPE = C.TOX_PROXY_TYPE_NONE
		if options.ProxyType == TOX_PROXY_TYPE_HTTP {
			cProxyType = C.TOX_PROXY_TYPE_HTTP
		} else if options.ProxyType == TOX_PROXY_TYPE_SOCKS5 {
			cProxyType = C.TOX_PROXY_TYPE_SOCKS5
		}
		cOptions.proxy_type = cProxyType

		// max ProxyHost length is 255
		if len(options.ProxyHost) > 255 {
			return nil, ErrArgs
		}
		cProxyHost := C.CString(options.ProxyHost)
		cOptions.proxy_host = cProxyHost
		defer C.free(unsafe.Pointer(cProxyHost))

		cOptions.proxy_port = C.uint16_t(options.ProxyPort)
		cOptions.start_port = C.uint16_t(options.StartPort)
		cOptions.end_port = C.uint16_t(options.EndPort)
		cOptions.tcp_port = C.uint16_t(options.TcpPort)

		if options.SaveDataType == TOX_SAVEDATA_TYPE_TOX_SAVE {
			cOptions.savedata_type = C.TOX_SAVEDATA_TYPE_TOX_SAVE
		} else if options.SaveDataType == TOX_SAVEDATA_TYPE_SECRET_KEY {
			cOptions.savedata_type = C.TOX_SAVEDATA_TYPE_SECRET_KEY
		}

		if len(options.SaveData) > 0 {
			cOptions.savedata_data = (*C.uint8_t)(&options.SaveData[0])
		} else {
			cOptions.savedata_data = nil
		}

		cOptions.savedata_length = C.size_t(len(options.SaveData))
	}

	cTox = C.tox_new(cOptions, &toxErrNew)
	if cTox == nil || ToxErrNew(toxErrNew) != TOX_ERR_NEW_OK {
		C.tox_options_free(cOptions)
		switch ToxErrNew(toxErrNew) {
		case TOX_ERR_NEW_NULL:
			return nil, ErrArgs
		case TOX_ERR_NEW_MALLOC:
			return nil, ErrNewMalloc
		case TOX_ERR_NEW_PORT_ALLOC:
			return nil, ErrNewPortAlloc
		case TOX_ERR_NEW_PROXY_BAD_TYPE:
			return nil, ErrNewProxy
		case TOX_ERR_NEW_PROXY_BAD_HOST:
			return nil, ErrNewProxy
		case TOX_ERR_NEW_PROXY_BAD_PORT:
			return nil, ErrNewProxy
		case TOX_ERR_NEW_PROXY_NOT_FOUND:
			return nil, ErrNewProxy
		case TOX_ERR_NEW_LOAD_ENCRYPTED:
			return nil, ErrNewLoadEnc
		case TOX_ERR_NEW_LOAD_BAD_FORMAT:
			return nil, ErrNewLoadBadFormat
		}

		if cTox == nil {
			return nil, ErrToxNew
		}

		return nil, ErrUnknown
	}

	t := &Tox{tox: cTox, cOptions: cOptions}
	return t, nil
}

/* Kill releases all resources associated with the Tox instance and disconnects
 * from the network.
 * After calling this function `t *TOX` becomes invalid. Do not use it again! */
func (t *Tox) Kill() error {
	if t.tox == nil {
		return ErrToxInit
	}

	C.tox_options_free(t.cOptions)
	C.tox_kill(t.tox)

	return nil
}

/* GetSaveDataSize returns the size of the savedata returned by GetSavedata. */
func (t *Tox) GetSaveDataSize() (uint32, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	return uint32(C.tox_get_savedata_size(t.tox)), nil
}

/* GetSavedata returns a byte slice of all information associated with the tox
 * instance. */
func (t *Tox) GetSavedata() ([]byte, error) {
	if t.tox == nil {
		return nil, ErrToxInit
	}
	size, err := t.GetSaveDataSize()
	if err != nil || size == 0 {
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
		return ErrToxInit
	}

	if len(publickey) != TOX_PUBLIC_KEY_SIZE {
		return ErrArgs
	}

	caddr := C.CString(address)
	defer C.free(unsafe.Pointer(caddr))

	var toxErrBootstrap C.TOX_ERR_BOOTSTRAP
	success := C.tox_bootstrap(t.tox, caddr, (C.uint16_t)(port), (*C.uint8_t)(&publickey[0]), &toxErrBootstrap)

	switch ToxErrBootstrap(toxErrBootstrap) {
	case TOX_ERR_BOOTSTRAP_OK:
		return nil
	case TOX_ERR_BOOTSTRAP_NULL:
		return ErrArgs
	case TOX_ERR_BOOTSTRAP_BAD_HOST:
		return ErrFuncFail
	case TOX_ERR_BOOTSTRAP_BAD_PORT:
		return ErrFuncFail
	}

	if !bool(success) {
		return ErrFuncFail
	}

	return ErrUnknown
}

/* AddTCPRelay adds the given node with IP, port, and public key without using
 * it as a boostrap node. */
func (t *Tox) AddTCPRelay(address string, port uint16, publickey []byte) error {
	if t.tox == nil {
		return ErrToxInit
	}

	if len(publickey) != TOX_PUBLIC_KEY_SIZE {
		return ErrArgs
	}

	caddr := C.CString(address)
	defer C.free(unsafe.Pointer(caddr))

	var toxErrBootstrap C.TOX_ERR_BOOTSTRAP
	success := C.tox_add_tcp_relay(t.tox, caddr, (C.uint16_t)(port), (*C.uint8_t)(&publickey[0]), &toxErrBootstrap)

	switch ToxErrBootstrap(toxErrBootstrap) {
	case TOX_ERR_BOOTSTRAP_OK:
		return nil
	case TOX_ERR_BOOTSTRAP_NULL:
		return ErrArgs
	case TOX_ERR_BOOTSTRAP_BAD_HOST:
		return ErrFuncFail
	case TOX_ERR_BOOTSTRAP_BAD_PORT:
		return ErrFuncFail
	}

	if !bool(success) {
		return ErrFuncFail
	}

	return ErrUnknown
}

/* SelfGetConnectionStatus returns true if Tox is connected to the DHT. */
func (t *Tox) SelfGetConnectionStatus() (ToxConnection, error) {
	if t.tox == nil {
		return TOX_CONNECTION_NONE, ErrToxInit
	}

	return ToxConnection(C.tox_self_get_connection_status(t.tox)), nil
}

/* IterationInterval returns the time in milliseconds before Iterate() should be
 * called again. */
func (t *Tox) IterationInterval() (uint32, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	ret := C.tox_iteration_interval(t.tox)

	return uint32(ret), nil
}

/* Iterate is the main loop. It needs to be called every IterationInterval()
 * milliseconds. */
func (t *Tox) Iterate() error {
	if t.tox == nil {
		return ErrToxInit
	}

	t.mtx.Lock()
	C.tox_iterate(t.tox)
	t.mtx.Unlock()

	return nil
}

/* SelfGetAddress returns the public address to give to others. */
func (t *Tox) SelfGetAddress() ([]byte, error) {
	if t.tox == nil {
		return nil, ErrToxInit
	}

	address := make([]byte, TOX_ADDRESS_SIZE)
	C.tox_self_get_address(t.tox, (*C.uint8_t)(&address[0]))

	return address, nil
}

/* SelfSetNospam sets the nospam of your ID. */
func (t *Tox) SelfSetNospam(nospam uint32) error {
	if t.tox == nil {
		return ErrToxInit
	}

	C.tox_self_set_nospam(t.tox, (C.uint32_t)(nospam))
	return nil
}

/* SelfGetNospam returns the nospam of your ID. */
func (t *Tox) SelfGetNospam() (uint32, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	n := C.tox_self_get_nospam(t.tox)
	return uint32(n), nil
}

/* SelfGetPublicKey returns the publickey of your profile. */
func (t *Tox) SelfGetPublicKey() ([]byte, error) {
	if t.tox == nil {
		return nil, ErrToxInit
	}

	publickey := make([]byte, TOX_PUBLIC_KEY_SIZE)

	C.tox_self_get_public_key(t.tox, (*C.uint8_t)(&publickey[0]))
	return publickey, nil
}

/* SelfGetSecretKey returns the secretkey of your profile. */
func (t *Tox) SelfGetSecretKey() ([]byte, error) {
	if t.tox == nil {
		return nil, ErrToxInit
	}

	secretkey := make([]byte, TOX_SECRET_KEY_SIZE)

	C.tox_self_get_secret_key(t.tox, (*C.uint8_t)(&secretkey[0]))
	return secretkey, nil
}

/* SelfSetName sets your nickname. The maximum name length is MAX_NAME_LENGTH. */
func (t *Tox) SelfSetName(name string) error {
	if t.tox == nil {
		return ErrToxInit
	}

	var cName (*C.uint8_t)

	if len(name) == 0 {
		cName = nil
	} else {
		cName = (*C.uint8_t)(&[]byte(name)[0])
	}

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
		return 0, ErrToxInit
	}

	ret := C.tox_self_get_name_size(t.tox)

	return int64(ret), nil
}

/* SelfGetName returns your nickname. */
func (t *Tox) SelfGetName() (string, error) {
	if t.tox == nil {
		return "", ErrToxInit
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
		return ErrToxInit
	}

	var cStatus (*C.uint8_t)

	if len(status) == 0 {
		cStatus = nil
	} else {
		cStatus = (*C.uint8_t)(&[]byte(status)[0])
	}

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
		return 0, ErrToxInit
	}

	ret := C.tox_self_get_status_message_size(t.tox)

	return int64(ret), nil
}

/* SelfGetStatusMessage returns your status message. */
func (t *Tox) SelfGetStatusMessage() (string, error) {
	if t.tox == nil {
		return "", ErrToxInit
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
		return ErrToxInit
	}

	C.tox_self_set_status(t.tox, (C.TOX_USER_STATUS)(userstatus))

	return nil
}

/* SelfGetStatus returns your status. */
func (t *Tox) SelfGetStatus() (ToxUserStatus, error) {
	if t.tox == nil {
		return TOX_USERSTATUS_NONE, ErrToxInit
	}

	n := C.tox_self_get_status(t.tox)

	return ToxUserStatus(n), nil
}

/* FriendAdd adds a friend by sending a friend request containing the given
 * message.
 * Returns the friend number on success, or a ToxErrFriendAdd on failure.
 */
func (t *Tox) FriendAdd(address []byte, message string) (uint32, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	if len(address) != TOX_ADDRESS_SIZE || len(message) == 0 {
		return 0, ErrArgs
	}

	caddr := (*C.uint8_t)(&address[0])
	cmessage := (*C.uint8_t)(&[]byte(message)[0])

	var toxErrFriendAdd C.TOX_ERR_FRIEND_ADD
	ret := C.tox_friend_add(t.tox, caddr, cmessage, (C.size_t)(len(message)), &toxErrFriendAdd)

	switch ToxErrFriendAdd(toxErrFriendAdd) {
	case TOX_ERR_FRIEND_ADD_OK:
		return uint32(ret), nil
	case TOX_ERR_FRIEND_ADD_NULL:
		return uint32(ret), ErrArgs
	case TOX_ERR_FRIEND_ADD_TOO_LONG:
		return uint32(ret), ErrFriendAddTooLong
	case TOX_ERR_FRIEND_ADD_NO_MESSAGE:
		return uint32(ret), ErrFriendAddNoMessage
	case TOX_ERR_FRIEND_ADD_OWN_KEY:
		return uint32(ret), ErrFriendAddOwnKey
	case TOX_ERR_FRIEND_ADD_ALREADY_SENT:
		return uint32(ret), ErrFriendAddAlreadySent
	case TOX_ERR_FRIEND_ADD_BAD_CHECKSUM:
		return uint32(ret), ErrFriendAddBadChecksum
	case TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM:
		return uint32(ret), ErrFriendAddSetNewNospam
	case TOX_ERR_FRIEND_ADD_MALLOC:
		return uint32(ret), ErrFriendAddNoMem
	default:
		return uint32(ret), ErrFuncFail
	}

	return uint32(ret), ErrUnknown
}

/* FriendAddNorequest adds a friend without sending a friend request.
 * Returns the friend number on success.
 */
func (t *Tox) FriendAddNorequest(publickey []byte) (uint32, error) {
	if t.tox == nil {
		return C.UINT32_MAX, ErrToxInit
	}

	if len(publickey) != TOX_PUBLIC_KEY_SIZE {
		return C.UINT32_MAX, ErrArgs
	}

	var toxErrFriendAdd C.TOX_ERR_FRIEND_ADD
	ret := C.tox_friend_add_norequest(t.tox, (*C.uint8_t)(&publickey[0]), &toxErrFriendAdd)
	if ret == C.UINT32_MAX {
		return C.UINT32_MAX, ErrFuncFail
	}

	switch ToxErrFriendAdd(toxErrFriendAdd) {
	case TOX_ERR_FRIEND_ADD_OK:
		return uint32(ret), nil
	case TOX_ERR_FRIEND_ADD_NULL:
		return uint32(ret), ErrArgs
	case TOX_ERR_FRIEND_ADD_TOO_LONG:
		return uint32(ret), ErrFriendAddTooLong
	case TOX_ERR_FRIEND_ADD_NO_MESSAGE:
		return uint32(ret), ErrFriendAddNoMessage
	case TOX_ERR_FRIEND_ADD_OWN_KEY:
		return uint32(ret), ErrFriendAddOwnKey
	case TOX_ERR_FRIEND_ADD_ALREADY_SENT:
		return uint32(ret), ErrFriendAddAlreadySent
	case TOX_ERR_FRIEND_ADD_BAD_CHECKSUM:
		return uint32(ret), ErrFriendAddBadChecksum
	case TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM:
		return uint32(ret), ErrFriendAddSetNewNospam
	case TOX_ERR_FRIEND_ADD_MALLOC:
		return uint32(ret), ErrFriendAddNoMem
	default:
		return uint32(ret), ErrFuncFail
	}

	return uint32(ret), ErrUnknown
}

/* FriendDelete removes a friend. */
func (t *Tox) FriendDelete(friendNumber uint32) error {
	if t.tox == nil {
		return ErrToxInit
	}

	var toxErrFriendDelete C.TOX_ERR_FRIEND_DELETE = C.TOX_ERR_FRIEND_DELETE_OK
	C.tox_friend_delete(t.tox, (C.uint32_t)(friendNumber), &toxErrFriendDelete)

	switch ToxErrFriendDelete(toxErrFriendDelete) {
	case TOX_ERR_FRIEND_DELETE_OK:
		return nil
	case TOX_ERR_FRIEND_DELETE_FRIEND_NOT_FOUND:
		return ErrArgs
	default:
		return ErrFuncFail
	}

	return ErrUnknown
}

/* FriendByPublicKey returns the friend number associated to a given publickey. */
func (t *Tox) FriendByPublicKey(publickey []byte) (uint32, error) {
	if t.tox == nil {
		return C.UINT32_MAX, ErrToxInit
	}

	if len(publickey) != TOX_PUBLIC_KEY_SIZE {
		return C.UINT32_MAX, ErrArgs
	}

	var toxErrFriendByPublicKey C.TOX_ERR_FRIEND_BY_PUBLIC_KEY
	n := C.tox_friend_by_public_key(t.tox, (*C.uint8_t)(&publickey[0]), &toxErrFriendByPublicKey)

	switch ToxErrFriendByPublicKey(toxErrFriendByPublicKey) {
	case TOX_ERR_FRIEND_BY_PUBLIC_KEY_OK:
		return uint32(n), nil
	case TOX_ERR_FRIEND_BY_PUBLIC_KEY_NULL:
		return uint32(n), ErrArgs
	default:
		return uint32(n), ErrFuncFail
	}

	return uint32(n), ErrUnknown
}

/* FriendExists returns true if a friend exists with given friendNumber. */
func (t *Tox) FriendExists(friendNumber uint32) (bool, error) {
	if t.tox == nil {
		return false, ErrToxInit
	}

	success := C.tox_friend_exists(t.tox, (C.uint32_t)(friendNumber))

	return bool(success), nil
}

/* SelfGetFriendlistSize returns the number of friends on the friendlist. */
func (t *Tox) SelfGetFriendlistSize() (int64, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}
	n := C.tox_self_get_friend_list_size(t.tox)

	return int64(n), nil
}

/* SelfGetFriendlist returns a slice of uint32 containing the friendNumbers. */
func (t *Tox) SelfGetFriendlist() ([]uint32, error) {
	if t.tox == nil {
		return nil, ErrToxInit
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

/* FriendGetPublickey returns the publickey associated to that friendNumber. */
func (t *Tox) FriendGetPublickey(friendNumber uint32) ([]byte, error) {
	if t.tox == nil {
		return nil, ErrToxInit
	}
	publickey := make([]byte, TOX_PUBLIC_KEY_SIZE)
	var toxErrFriendGetPublicKey C.TOX_ERR_FRIEND_GET_PUBLIC_KEY = C.TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK
	C.tox_friend_get_public_key(t.tox, (C.uint32_t)(friendNumber), (*C.uint8_t)(&publickey[0]), &toxErrFriendGetPublicKey)

	switch ToxErrFriendGetPublicKey(toxErrFriendGetPublicKey) {
	case TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK:
		return publickey, nil
	case TOX_ERR_FRIEND_GET_PUBLIC_KEY_FRIEND_NOT_FOUND:
		return nil, ErrArgs
	default:
		return nil, ErrFuncFail
	}

	return nil, ErrUnknown
}

/* FriendGetLastOnline returns the timestamp of the last time the friend with
 * the given friendNumber was seen online. */
func (t *Tox) FriendGetLastOnline(friendNumber uint32) (time.Time, error) {
	if t.tox == nil {
		return time.Time{}, ErrToxInit
	}

	var toxErrFriendGetLastOnline C.TOX_ERR_FRIEND_GET_LAST_ONLINE = C.TOX_ERR_FRIEND_GET_LAST_ONLINE_OK
	ret := C.tox_friend_get_last_online(t.tox, (C.uint32_t)(friendNumber), &toxErrFriendGetLastOnline)

	if ret == C.INT64_MAX || ToxErrFriendGetLastOnline(toxErrFriendGetLastOnline) != TOX_ERR_FRIEND_GET_LAST_ONLINE_OK {
		return time.Time{}, ErrFuncFail
	}

	last := time.Unix(int64(ret), 0)

	return last, nil
}

/* FriendGetNameSize returns the length of the name of friendNumber. */
func (t *Tox) FriendGetNameSize(friendNumber uint32) (int64, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	ret := C.tox_friend_get_name_size(t.tox, (C.uint32_t)(friendNumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return 0, ErrFuncFail
	}

	return int64(ret), nil
}

/* FriendGetName returns the name of friendNumber. */
func (t *Tox) FriendGetName(friendNumber uint32) (string, error) {
	if t.tox == nil {
		return "", ErrToxInit
	}

	length, err := t.FriendGetNameSize(friendNumber)
	if err != nil {
		return "", ErrFuncFail
	}

	name := make([]byte, length)

	if length > 0 {
		var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
		success := C.tox_friend_get_name(t.tox, (C.uint32_t)(friendNumber), (*C.uint8_t)(&name[0]), &toxErrFriendQuery)

		if success != true || ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
			return "", ErrFuncFail
		}
	}

	return string(name), nil
}

/* FriendGetStatusMessageSize returns the size of the status of a friend with
 * the given friendNumber.
 */
func (t *Tox) FriendGetStatusMessageSize(friendNumber uint32) (int64, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	ret := C.tox_friend_get_status_message_size(t.tox, (C.uint32_t)(friendNumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return 0, ErrFuncFail
	}

	return int64(ret), nil
}

/* FriendGetStatusMessage returns the status message of friend with the given
 * friendNumber.
 */
func (t *Tox) FriendGetStatusMessage(friendNumber uint32) (string, error) {
	if t.tox == nil {
		return "", ErrToxInit
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK

	size, error := t.FriendGetStatusMessageSize(friendNumber)
	if error != nil {
		return "", ErrFuncFail
	}

	statusMessage := make([]byte, size)

	if size > 0 {
		toxErrFriendQuery = C.TOX_ERR_FRIEND_QUERY_OK
		n := C.tox_friend_get_status_message(t.tox, (C.uint32_t)(friendNumber), (*C.uint8_t)(&statusMessage[0]), &toxErrFriendQuery)

		if n != true || ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
			return "", ErrFuncFail
		}
	}

	return string(statusMessage), nil
}

/* FriendGetStatus returns the status of friendNumber. */
func (t *Tox) FriendGetStatus(friendNumber uint32) (ToxUserStatus, error) {
	if t.tox == nil {
		return TOX_USERSTATUS_NONE, ErrToxInit
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	status := C.tox_friend_get_status(t.tox, (C.uint32_t)(friendNumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return TOX_USERSTATUS_NONE, ErrFuncFail
	}

	return ToxUserStatus(status), nil
}

/* FriendGetConnectionStatus returns true if the friend is connected. */
func (t *Tox) FriendGetConnectionStatus(friendNumber uint32) (ToxConnection, error) {
	if t.tox == nil {
		return TOX_CONNECTION_NONE, ErrToxInit
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	status := C.tox_friend_get_connection_status(t.tox, (C.uint32_t)(friendNumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return TOX_CONNECTION_NONE, ErrFuncFail
	}

	return ToxConnection(status), nil
}

/* FriendGetTyping returns true if friendNumber is typing. */
func (t *Tox) FriendGetTyping(friendNumber uint32) (bool, error) {
	if t.tox == nil {
		return false, ErrToxInit
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	istyping := C.tox_friend_get_typing(t.tox, (C.uint32_t)(friendNumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return false, ErrFuncFail
	}

	return bool(istyping), nil
}

/* SelfSetTyping sets your typing status to a friend. */
func (t *Tox) SelfSetTyping(friendNumber uint32, typing bool) error {
	if t.tox == nil {
		return ErrToxInit
	}

	var toxErrSetTyping C.TOX_ERR_SET_TYPING = C.TOX_ERR_SET_TYPING_OK
	success := C.tox_self_set_typing(t.tox, (C.uint32_t)(friendNumber), (C._Bool)(typing), &toxErrSetTyping)

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
func (t *Tox) FriendSendMessage(friendNumber uint32, messagetype ToxMessageType, message string) (uint32, error) {
	if t.tox == nil {
		return 0, ErrToxInit
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
	n := C.tox_friend_send_message(t.tox, (C.uint32_t)(friendNumber), cMessageType, cMessage, (C.size_t)(len(message)), &toxFriendSendMessageError)

	if ToxErrFriendSendMessage(toxFriendSendMessageError) != TOX_ERR_FRIEND_SEND_MESSAGE_OK {
		return 0, ErrFuncFail
	}

	return uint32(n), nil
}

/* Hash generates a cryptographic hash of the given data (can be used to cache
 * avatars). */
func (t *Tox) Hash(data []byte) ([]byte, error) {
	if t.tox == nil {
		return nil, ErrToxInit
	}

	var cData *C.uint8_t

	if len(data) == 0 {
		cData = nil
	} else {
		cData = (*C.uint8_t)(&data[0])
	}

	hash := make([]byte, TOX_HASH_LENGTH)

	success := C.tox_hash((*C.uint8_t)(&hash[0]), cData, C.size_t(len(data)))
	if !bool(success) {
		return nil, ErrFuncFail
	}

	return hash, nil
}

/* FileControl sends a FileControl to a friend with the given friendNumber. */
func (t *Tox) FileControl(friendNumber uint32, fileNumber uint32, fileControl ToxFileControl) error {
	if t.tox == nil {
		return ErrToxInit
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
	success := C.tox_file_control(t.tox, (C.uint32_t)(friendNumber), (C.uint32_t)(fileNumber), cFileControl, &toxErrFileControl)

	if !bool(success) || ToxErrFileControl(toxErrFileControl) != TOX_ERR_FILE_CONTROL_OK {
		return ErrFuncFail
	}

	return nil
}

/* FileSeek sends a file seek control command to a friend for a given file
 * transfer. */
func (t *Tox) FileSeek(friendNumber uint32, fileNumber uint32, position uint64) error {
	if t.tox == nil {
		return ErrToxInit
	}

	var toxErrFileSeek C.TOX_ERR_FILE_SEEK
	success := C.tox_file_seek(t.tox, C.uint32_t(friendNumber), C.uint32_t(fileNumber), C.uint64_t(position), &toxErrFileSeek)

	if !bool(success) || ToxErrFileSeek(toxErrFileSeek) != TOX_ERR_FILE_SEEK_OK {
		return ErrFuncFail
	}

	return nil
}

/* FileGetFileId returns the file id associated to the file transfer. */
func (t *Tox) FileGetFileId(friendNumber uint32, fileNumber uint32) ([]byte, error) {
	if t.tox == nil {
		return nil, ErrToxInit
	}

	fileId := make([]byte, TOX_FILE_ID_LENGTH)

	var toxErrFileGet C.TOX_ERR_FILE_GET
	success := C.tox_file_get_file_id(t.tox, C.uint32_t(friendNumber), C.uint32_t(fileNumber), (*C.uint8_t)(&fileId[0]), &toxErrFileGet)
	if !bool(success) || ToxErrFileGet(toxErrFileGet) != TOX_ERR_FILE_GET_OK {
		return nil, ErrFuncFail
	}

	return fileId, nil
}

/* FileSend sends a file transmission request. */
func (t *Tox) FileSend(friendNumber uint32, fileKind ToxFileKind, fileLength uint64, fileID []byte, fileName string) (uint32, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	var cFileKind = C.TOX_FILE_KIND_DATA
	switch ToxFileKind(fileKind) {
	case TOX_FILE_KIND_AVATAR:
		cFileKind = C.TOX_FILE_KIND_AVATAR
	case TOX_FILE_KIND_DATA:
		cFileKind = C.TOX_FILE_KIND_DATA
	}

	var cFileID *C.uint8_t

	if fileID == nil {
		cFileID = nil
	} else {
		if len(fileID) != TOX_FILE_ID_LENGTH {
			return 0, ErrFileSendInvalidFileID
		}

		cFileID = (*C.uint8_t)(&[]byte(fileID)[0])
	}

	if len(fileName) == 0 {
		return 0, ErrArgs
	}

	cFileName := (*C.uint8_t)(&[]byte(fileName)[0])

	var toxErrFileSend C.TOX_ERR_FILE_SEND
	n := C.tox_file_send(t.tox, (C.uint32_t)(friendNumber), (C.uint32_t)(cFileKind), (C.uint64_t)(fileLength), cFileID, cFileName, (C.size_t)(len(fileName)), &toxErrFileSend)

	if n == C.UINT32_MAX || ToxErrFileSend(toxErrFileSend) != TOX_ERR_FILE_SEND_OK {
		return 0, ErrFuncFail
	}
	return uint32(n), nil
}

/* FileSendChunk sends a chunk of file data to a friend. */
func (t *Tox) FileSendChunk(friendNumber uint32, fileNumber uint32, position uint64, data []byte) error {
	if t.tox == nil {
		return ErrToxInit
	}

	var cData *C.uint8_t

	if len(data) == 0 {
		cData = nil
	} else {
		cData = (*C.uint8_t)(&data[0])
	}

	var toxErrFileSendChunk C.TOX_ERR_FILE_SEND_CHUNK
	success := C.tox_file_send_chunk(t.tox, (C.uint32_t)(friendNumber), (C.uint32_t)(fileNumber), (C.uint64_t)(position), cData, (C.size_t)(len(data)), &toxErrFileSendChunk)

	if !bool(success) || ToxErrFileSendChunk(toxErrFileSendChunk) != TOX_ERR_FILE_SEND_CHUNK_OK {
		return ErrFuncFail
	}
	return nil
}

/* FriendSendLossyPacket sends a custom lossy packet to a friend.
 * The first byte of data must be in the range 200-254. Maximum length of a
 * custom packet is TOX_MAX_CUSTOM_PACKET_SIZE. */
func (t *Tox) FriendSendLossyPacket(friendNumber uint32, data []byte) error {
	if t.tox == nil {
		return ErrToxInit
	}

	var cData *C.uint8_t

	if len(data) == 0 {
		cData = nil
	} else {
		cData = (*C.uint8_t)(&data[0])
	}

	var toxErrFriendCustomPacket C.TOX_ERR_FRIEND_CUSTOM_PACKET
	C.tox_friend_send_lossy_packet(t.tox, C.uint32_t(friendNumber), cData, C.size_t(len(data)), &toxErrFriendCustomPacket)

	switch ToxErrFriendCustomPacket(toxErrFriendCustomPacket) {
	case TOX_ERR_FRIEND_CUSTOM_PACKET_OK:
		return nil
	case TOX_ERR_FRIEND_CUSTOM_PACKET_NULL:
		return ErrArgs
	default:
		return ErrFuncFail
	}

	return ErrUnknown
}

/* FriendSendLosslessPacket sends a custom lossless packet to a friend.
 * The first byte of data must be in the range 160-191. Maximum length of a
 * custom packet is TOX_MAX_CUSTOM_PACKET_SIZE. */
func (t *Tox) FriendSendLosslessPacket(friendNumber uint32, data []byte) error {
	if t.tox == nil {
		return ErrToxInit
	}

	var cData *C.uint8_t

	if len(data) == 0 {
		cData = nil
	} else {
		cData = (*C.uint8_t)(&data[0])
	}

	var toxErrFriendCustomPacket C.TOX_ERR_FRIEND_CUSTOM_PACKET
	C.tox_friend_send_lossless_packet(t.tox, C.uint32_t(friendNumber), cData, C.size_t(len(data)), &toxErrFriendCustomPacket)

	switch ToxErrFriendCustomPacket(toxErrFriendCustomPacket) {
	case TOX_ERR_FRIEND_CUSTOM_PACKET_OK:
		return nil
	case TOX_ERR_FRIEND_CUSTOM_PACKET_NULL:
		return ErrArgs
	default:
		return ErrFuncFail
	}

	return ErrUnknown
}

/* SelfGetDhtId returns the temporary DHT public key of this instance. */
func (t *Tox) SelfGetDhtId() ([]byte, error) {
	if t.tox == nil {
		return nil, ErrToxInit
	}

	publickey := make([]byte, TOX_PUBLIC_KEY_SIZE)

	C.tox_self_get_dht_id(t.tox, (*C.uint8_t)(&publickey[0]))
	return publickey, nil
}

/* SelfGetUDPPort returns the UDP port the Tox instance is bound to. */
func (t *Tox) SelfGetUDPPort() (uint16, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	var toxErrGetPort C.TOX_ERR_GET_PORT
	port := C.tox_self_get_udp_port(t.tox, &toxErrGetPort)

	if ToxErrGetPort(toxErrGetPort) != TOX_ERR_GET_PORT_OK {
		return 0, ErrFuncFail
	}

	return uint16(port), nil
}

/* SelfGetTCPPort returns the TCP port the Tox instance is bound to. This is
 * only relevant if the instance is acting as a TCP relay. */
func (t *Tox) SelfGetTCPPort() (uint16, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	var toxErrGetPort C.TOX_ERR_GET_PORT
	port := C.tox_self_get_tcp_port(t.tox, &toxErrGetPort)

	if ToxErrGetPort(toxErrGetPort) != TOX_ERR_GET_PORT_OK {
		return 0, ErrFuncFail
	}

	return uint16(port), nil
}
