package gotox

/* gotox - a go wrapper for toxcore
 *
 * This work is based on the great work by organ
 * (https://github.com/organ/gotox/).
 *
 * Pull requests, bug reporting and feature request (via github issues) are
 * always welcome. :)
 *
 * For a list of supported toxcore features see PROGRESS.md.
 */

//#cgo LDFLAGS: -ltoxcore
//#include <tox/tox.h>
import "C"
import "sync"

/* Tox is the main struct. */
type Tox struct {
	cOptions *C.struct_Tox_Options
	tox      *C.Tox
	mtx      sync.Mutex

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
	/* The type of socket to create.
	 * If IPv6Enabled is true, both IPv6 and IPv4 connections are allowed.
	 */
	IPv6Enabled bool

	/* Enable the use of UDP communication when available.
	 *
	 * Setting this to false will force Tox to use TCP only. Communications will
	 * need to be relayed through a TCP relay node, potentially slowing them down.
	 * Disabling UDP support is necessary when using anonymous proxies or Tor.
	 */
	UDPEnabled bool

	/* The type of the proxy (PROXY_TYPE_NONE, PROXY_TYPE_HTTP or PROXY_TYPE_SOCKS5). */
	ProxyType ToxProxyType

	/* The IP address or DNS name of the proxy to be used. */
	ProxyHost string

	/* The port to use to connect to the proxy server. */
	ProxyPort uint16

	/* The start port of the inclusive port range to attempt to use. */
	StartPort uint16

	/* The end port of the inclusive port range to attempt to use. */
	EndPort uint16

	/* The port to use for the TCP server. If 0, the tcp server is disabled. */
	TcpPort uint16

	/* The type of savedata to load from. */
	SaveDataType ToxSaveDataType

	/* The savedata. */
	SaveData []byte
}
