package gotox

/* gotox - a go wrapper for toxcore
 *
 * This work is based on the great work by organ
 * (https://github.com/organ/gotox/).
 *
 * Pull requests, bug reporting and feature request (via github issues) are
 * always welcome. :)
 *
 * TODO:
 * - groupchats
 * - sending files
 */

/*
#cgo LDFLAGS: -ltoxcore

#include <tox/tox.h>
#include <stdlib.h>
#include "hooks.c"
*/
import "C"
import "sync"
import "unsafe"

/* This event is triggered whenever there is a change in the DHT connection
 * state. When disconnected, a client may choose to call tox_bootstrap again, to
 * reconnect to the DHT. Note that this state may frequently change for short
 * amounts of time. Clients should therefore not immediately bootstrap on
 * receiving a disconnect.
 */
type OnSelfConnectionStatusChanges func(tox *Tox, status ConnectionStatus)

/* This event is triggered when a friend changes their name. */
type OnFriendNameChanges func(tox *Tox, friendnumber uint32, name string)

/* This event is triggered when a friend changes their status message. */
type OnFriendStatusMessageChanges func(tox *Tox, friendnumber uint32, message string)

/* This event is triggered when a friend changes their user status. */
type OnFriendStatusChanges func(tox *Tox, friendnumber uint32, userstatus UserStatus)

/* This event is triggered when a friend goes offline after having been online,
 * or when a friend goes online.
 *
 * This callback is not called when adding friends. It is assumed that when
 * adding friends, their connection status is initially offline.
 */
type OnFriendConnectionStatusChanges func(tox *Tox, friendnumber uint32, connectionstatus ConnectionStatus)

/* This event is triggered when a friend starts or stops typing. */
type OnFriendTypingChanges func(tox *Tox, friendnumber uint32, istyping bool)

/* This event is triggered when the friend receives the message with the
 * corresponding message ID.
 */
type OnFriendReadReceipt func(tox *Tox, friendnumber uint32, messageid uint32)

/* This event is triggered when a friend request is received. */
type OnFriendRequest func(tox *Tox, publickey []byte, message string)

/* This event is triggered when a message from a friend is received. */
type OnFriendMessage func(tox *Tox, friendnumber uint32, messagetype MessageType, message string)

/* This event is triggered when a file control command is received from a
 * friend.
 */
type OnFileRecvControl func(tox *Tox, friendnumber uint32, filenumber uint32, filecontrol FileControl)

/* This event is triggered when Core is ready to send more file data. */
type OnFileChunkRequest func(tox *Tox, friendnumber uint32, filenumber uint32, position uint64, length uint64)

/* This event is triggered when a file transfer request is received. */
type OnFileRecv func(tox *Tox, friendnumber uint32, filenumber uint32, kind uint32, filesize uint64, filename string)

/* This event is first triggered when a file transfer request is received, and
 * subsequently when a chunk of file data for an accepted request was received.
 */
type OnFileRecvChunk func(tox *Tox, friendnumber uint32, filenumber uint32, position uint64, data []byte)

/* This event is triggered when a lossy packet is received from a friend. */
type OnFriendLossyPacket func(tox *Tox, friendnumber uint32, data []byte)

/* This event is triggered when a lossless packet is received from a friend. */
type OnFriendLosslessPacket func(tox *Tox, friendnumber uint32, data []byte)

/* Tox is the main struct. */
type Tox struct {
	tox *C.Tox
	mtx sync.Mutex

	// Callbacks
	onSelfConnectionStatusChanges   OnSelfConnectionStatusChanges
	onFriendNameChanges             OnFriendNameChanges
	onFriendStatusMessageChanges    OnFriendStatusMessageChanges
	onFriendStatusChanges           OnFriendStatusChanges
	onFriendConnectionStatusChanges OnFriendConnectionStatusChanges
	onFriendTypingChanges           OnFriendTypingChanges
	onFriendReadReceipt             OnFriendReadReceipt
	onFriendRequest                 OnFriendRequest
	onFriendMessage                 OnFriendMessage
	onFileRecvControl               OnFileRecvControl
	onFileChunkRequest              OnFileChunkRequest
	onFileRecv                      OnFileRecv
	onFileRecvChunk                 OnFileRecvChunk
	onFriendLossyPacket             OnFriendLossyPacket
	onFriendLosslessPacket          OnFriendLosslessPacket
}

type Options struct {
	/**
	 * If IPv6Enabled is true, both IPv6 and IPv4 connections are allowed.
	 */
	IPv6Enabled bool

	/**
	 * Enable the use of UDP communication when available.
	 *
	 * Setting this to false will force Tox to use TCP only. Communications will
	 * need to be relayed through a TCP relay node, potentially slowing them down.
	 * Disabling UDP support is necessary when using anonymous proxies or Tor.
	 */
	UDPEnabled bool

	/* The type of the proxy (PROXY_TYPE_NONE, PROXY_TYPE_HTTP or PROXY_TYPE_SOCKS5). */
	TCPProxyType ProxyType

	/* The IP address or DNS name of the proxy to be used. */
	ProxyHost string

	/* The port to use to connect to the proxy server. */
	ProxyPort uint16

	/* The start port of the inclusive port range to attempt to use. */
	StartPort uint16

	/* The end port of the inclusive port range to attempt to use. */
	EndPort uint16
}

// New returns a new Tox instance.
func New(o *Options, data []byte) (*Tox, error) {
	var ctox *C.Tox
	var cToxNewError C.TOX_ERR_NEW

	var cData *C.uint8_t
	if len(data) > 0 {
		cData = (*C.uint8_t)(&data[0])
	} else {
		cData = nil
	}

	if o != nil {
		// Let's map o from Options to C.Tox_Options
		cIPv6Enabled := (C._Bool)(o.IPv6Enabled)
		cUDPEnabled := (C._Bool)(o.UDPEnabled)

		var cProxyType C.TOX_PROXY_TYPE
		if o.TCPProxyType == PROXY_TYPE_HTTP {
			cProxyType = C.TOX_PROXY_TYPE_HTTP
		} else if o.TCPProxyType == PROXY_TYPE_SOCKS5 {
			cProxyType = C.TOX_PROXY_TYPE_SOCKS5
		}

		// max ProxyHost length is 255
		if len(o.ProxyHost) > 255 {
			return nil, ErrArgs
		}
		cProxyHost := C.CString(o.ProxyHost)
		defer C.free(unsafe.Pointer(cProxyHost))

		cProxyPort := (C.uint16_t)(o.ProxyPort)
		cStartPort := (C.uint16_t)(o.StartPort)
		cEndPort := (C.uint16_t)(o.EndPort)

		co := &C.struct_Tox_Options{
			ipv6_enabled: cIPv6Enabled,
			udp_enabled:  cUDPEnabled,
			proxy_type:   cProxyType,
			proxy_host:   cProxyHost,
			proxy_port:   cProxyPort,
			start_port:   cStartPort,
			end_port:     cEndPort}

		ctox = C.tox_new(co, cData, (C.size_t)(len(data)), &cToxNewError)
	} else {
		ctox = C.tox_new(nil, cData, (C.size_t)(len(data)), &cToxNewError)
	}

	if ctox == nil || ToxNewError(cToxNewError) != ERR_NEW_OK {
		return nil, ErrInit
	}

	t := &Tox{tox: ctox}

	return t, nil
}

// GetSaveDataSize returns the size of the save data returned by GetSavedata.
func (t *Tox) GetSaveDataSize() (uint32, error) {
	if t.tox == nil {
		return 0, ErrBadTox
	}

	return uint32(C.tox_get_savedata_size(t.tox)), nil
}

// GetSavedata returns a byte slice of the save data.
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

// Kill stops a Tox instance.
func (t *Tox) Kill() error {
	if t.tox == nil {
		return ErrBadTox
	}
	C.tox_kill(t.tox)

	return nil
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

/* Iterate is the main loop needs to be called every IterationInterval()
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

/* BootstrapFromAddress resolves address into an IP address. If successful, it
 * sends a request to the given node to setup connection. */
func (t *Tox) BootstrapFromAddress(address string, port uint16, publickey []byte) error {
	if t.tox == nil {
		return ErrBadTox
	}

	if len(publickey) != PUBLIC_KEY_SIZE {
		return ErrArgs
	}

	caddr := C.CString(address)
	defer C.free(unsafe.Pointer(caddr))

	var cBootstrapError C.TOX_ERR_BOOTSTRAP
	C.tox_bootstrap(t.tox, caddr, (C.uint16_t)(port), (*C.uint8_t)(&publickey[0]), &cBootstrapError)

	var bootstrapError error

	switch BootstrapError(cBootstrapError) {
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
