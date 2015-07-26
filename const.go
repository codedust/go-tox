package gotox

//#include <tox/tox.h>
import "C"
import "errors"

const (
	TOX_PUBLIC_KEY_SIZE           = C.TOX_PUBLIC_KEY_SIZE
	TOX_SECRET_KEY_SIZE           = C.TOX_SECRET_KEY_SIZE
	TOX_ADDRESS_SIZE              = C.TOX_ADDRESS_SIZE
	TOX_MAX_NAME_LENGTH           = C.TOX_MAX_NAME_LENGTH
	TOX_MAX_STATUS_MESSAGE_LENGTH = C.TOX_MAX_STATUS_MESSAGE_LENGTH
	TOX_MAX_FRIEND_REQUEST_LENGTH = C.TOX_MAX_FRIEND_REQUEST_LENGTH
	TOX_MAX_MESSAGE_LENGTH        = C.TOX_MAX_MESSAGE_LENGTH
	TOX_MAX_CUSTOM_PACKET_SIZE    = C.TOX_MAX_CUSTOM_PACKET_SIZE
	TOX_HASH_LENGTH               = C.TOX_HASH_LENGTH
	TOX_FILE_ID_LENGTH            = C.TOX_FILE_ID_LENGTH
	TOX_MAX_FILENAME_LENGTH       = C.TOX_MAX_FILENAME_LENGTH
)

type ToxUserStatus C.enum_TOX_USER_STATUS

const (
	TOX_USERSTATUS_NONE ToxUserStatus = C.TOX_USER_STATUS_NONE
	TOX_USERSTATUS_AWAY ToxUserStatus = C.TOX_USER_STATUS_AWAY
	TOX_USERSTATUS_BUSY ToxUserStatus = C.TOX_USER_STATUS_BUSY
)

type ToxMessageType C.enum_TOX_MESSAGE_TYPE

const (
	TOX_MESSAGE_TYPE_NORMAL ToxMessageType = C.TOX_MESSAGE_TYPE_NORMAL
	TOX_MESSAGE_TYPE_ACTION ToxMessageType = C.TOX_MESSAGE_TYPE_ACTION
)

type ToxProxyType C.enum_TOX_PROXY_TYPE

const (
	TOX_PROXY_TYPE_NONE   ToxProxyType = C.TOX_PROXY_TYPE_NONE
	TOX_PROXY_TYPE_HTTP   ToxProxyType = C.TOX_PROXY_TYPE_HTTP
	TOX_PROXY_TYPE_SOCKS5 ToxProxyType = C.TOX_PROXY_TYPE_SOCKS5
)

type ToxSaveDataType C.enum_TOX_SAVEDATA_TYPE

const (
	TOX_SAVEDATA_TYPE_NONE       ToxSaveDataType = C.TOX_SAVEDATA_TYPE_NONE
	TOX_SAVEDATA_TYPE_TOX_SAVE   ToxSaveDataType = C.TOX_SAVEDATA_TYPE_TOX_SAVE
	TOX_SAVEDATA_TYPE_SECRET_KEY ToxSaveDataType = C.TOX_SAVEDATA_TYPE_SECRET_KEY
)

type ToxErrOptionsNew C.enum_TOX_ERR_OPTIONS_NEW

const (
	TOX_ERR_OPTIONS_NEW_OK     ToxErrOptionsNew = C.TOX_ERR_OPTIONS_NEW_OK
	TOX_ERR_OPTIONS_NEW_MALLOC ToxErrOptionsNew = C.TOX_ERR_OPTIONS_NEW_MALLOC
)

type ToxConnection C.enum_TOX_CONNECTION

const (
	TOX_CONNECTION_NONE ToxConnection = C.TOX_CONNECTION_NONE
	TOX_CONNECTION_TCP  ToxConnection = C.TOX_CONNECTION_TCP
	TOX_CONNECTION_UDP  ToxConnection = C.TOX_CONNECTION_UDP
)

type ToxFileKind C.enum_TOX_FILE_KIND

const (
	TOX_FILE_KIND_DATA   ToxFileKind = C.TOX_FILE_KIND_DATA
	TOX_FILE_KIND_AVATAR ToxFileKind = C.TOX_FILE_KIND_AVATAR
)

type ToxFileControl C.enum_TOX_FILE_CONTROL

const (
	TOX_FILE_CONTROL_RESUME ToxFileControl = C.TOX_FILE_CONTROL_RESUME
	TOX_FILE_CONTROL_PAUSE  ToxFileControl = C.TOX_FILE_CONTROL_PAUSE
	TOX_FILE_CONTROL_CANCEL ToxFileControl = C.TOX_FILE_CONTROL_CANCEL
)

/* === Errors === */
// General errors
var (
	ErrToxNew   = errors.New("Error initializing Tox")
	ErrToxInit  = errors.New("Tox not initialized")
	ErrArgs     = errors.New("Nil arguments or wrong size")
	ErrFuncFail = errors.New("Function failed")
	ErrUnknown  = errors.New("An unknown error occoured")
)

type ToxErrNew C.enum_TOX_ERR_NEW

const (
	TOX_ERR_NEW_OK              ToxErrNew = C.TOX_ERR_NEW_OK
	TOX_ERR_NEW_NULL            ToxErrNew = C.TOX_ERR_NEW_NULL
	TOX_ERR_NEW_MALLOC          ToxErrNew = C.TOX_ERR_NEW_MALLOC
	TOX_ERR_NEW_PORT_ALLOC      ToxErrNew = C.TOX_ERR_NEW_PORT_ALLOC
	TOX_ERR_NEW_PROXY_BAD_TYPE  ToxErrNew = C.TOX_ERR_NEW_PROXY_BAD_TYPE
	TOX_ERR_NEW_PROXY_BAD_HOST  ToxErrNew = C.TOX_ERR_NEW_PROXY_BAD_HOST
	TOX_ERR_NEW_PROXY_BAD_PORT  ToxErrNew = C.TOX_ERR_NEW_PROXY_BAD_PORT
	TOX_ERR_NEW_PROXY_NOT_FOUND ToxErrNew = C.TOX_ERR_NEW_PROXY_NOT_FOUND
	TOX_ERR_NEW_LOAD_ENCRYPTED  ToxErrNew = C.TOX_ERR_NEW_LOAD_ENCRYPTED
	TOX_ERR_NEW_LOAD_BAD_FORMAT ToxErrNew = C.TOX_ERR_NEW_LOAD_BAD_FORMAT
)

type ToxErrBootstrap C.enum_TOX_ERR_BOOTSTRAP

const (
	TOX_ERR_BOOTSTRAP_OK       ToxErrBootstrap = C.TOX_ERR_BOOTSTRAP_OK
	TOX_ERR_BOOTSTRAP_NULL     ToxErrBootstrap = C.TOX_ERR_BOOTSTRAP_NULL
	TOX_ERR_BOOTSTRAP_BAD_HOST ToxErrBootstrap = C.TOX_ERR_BOOTSTRAP_BAD_HOST
	TOX_ERR_BOOTSTRAP_BAD_PORT ToxErrBootstrap = C.TOX_ERR_BOOTSTRAP_BAD_PORT
)

var (
	ErrFriendAddNull         = errors.New("One of the arguments was NULL")
	ErrFriendAddTooLong      = errors.New("Message too long")
	ErrFriendAddNoMessage    = errors.New("Empty message")
	ErrFriendAddOwnKey       = errors.New("Own key")
	ErrFriendAddAlreadySent  = errors.New("Already sent")
	ErrFriendAddUnkown       = errors.New("Unknown error")
	ErrFriendAddBadChecksum  = errors.New("Bad checksum in address")
	ErrFriendAddSetNewNospam = errors.New("Different nospam")
	ErrFriendAddNoMem        = errors.New("Failed increasing friend list")
)

type ToxErrFriendAdd C.enum_TOX_ERR_FRIEND_ADD

const (
	TOX_ERR_FRIEND_ADD_OK             ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_OK
	TOX_ERR_FRIEND_ADD_NULL           ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_NULL
	TOX_ERR_FRIEND_ADD_TOO_LONG       ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_TOO_LONG
	TOX_ERR_FRIEND_ADD_NO_MESSAGE     ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_NO_MESSAGE
	TOX_ERR_FRIEND_ADD_OWN_KEY        ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_OWN_KEY
	TOX_ERR_FRIEND_ADD_ALREADY_SENT   ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_ALREADY_SENT
	TOX_ERR_FRIEND_ADD_BAD_CHECKSUM   ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_BAD_CHECKSUM
	TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM
	TOX_ERR_FRIEND_ADD_MALLOC         ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_MALLOC
)

type ToxErrFriendByPublicKey C.enum_TOX_ERR_FRIEND_BY_PUBLIC_KEY

const (
	TOX_ERR_FRIEND_BY_PUBLIC_KEY_OK        ToxErrFriendByPublicKey = C.TOX_ERR_FRIEND_BY_PUBLIC_KEY_OK
	TOX_ERR_FRIEND_BY_PUBLIC_KEY_NULL      ToxErrFriendByPublicKey = C.TOX_ERR_FRIEND_BY_PUBLIC_KEY_NULL
	TOX_ERR_FRIEND_BY_PUBLIC_KEY_NOT_FOUND ToxErrFriendByPublicKey = C.TOX_ERR_FRIEND_BY_PUBLIC_KEY_NOT_FOUND
)

type ToxErrFriendGetPublicKey C.enum_TOX_ERR_FRIEND_GET_PUBLIC_KEY

const (
	TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK               ToxErrFriendGetPublicKey = C.TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK
	TOX_ERR_FRIEND_GET_PUBLIC_KEY_FRIEND_NOT_FOUND ToxErrFriendGetPublicKey = C.TOX_ERR_FRIEND_GET_PUBLIC_KEY_FRIEND_NOT_FOUND
)

type ToxErrFriendDelete C.enum_TOX_ERR_FRIEND_DELETE

const (
	TOX_ERR_FRIEND_DELETE_OK               ToxErrFriendDelete = C.TOX_ERR_FRIEND_DELETE_OK
	TOX_ERR_FRIEND_DELETE_FRIEND_NOT_FOUND ToxErrFriendDelete = C.TOX_ERR_FRIEND_DELETE_FRIEND_NOT_FOUND
)

type ToxErrFriendQuery C.enum_TOX_ERR_FRIEND_QUERY

const (
	TOX_ERR_FRIEND_QUERY_OK               ToxErrFriendQuery = C.TOX_ERR_FRIEND_QUERY_OK
	TOX_ERR_FRIEND_QUERY_NULL             ToxErrFriendQuery = C.TOX_ERR_FRIEND_QUERY_NULL
	TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND ToxErrFriendQuery = C.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND
)

type ToxErrSetInfo C.enum_TOX_ERR_SET_INFO

const (
	TOX_ERR_SET_INFO_OK       ToxErrSetInfo = C.TOX_ERR_SET_INFO_OK
	TOX_ERR_SET_INFO_NULL     ToxErrSetInfo = C.TOX_ERR_SET_INFO_NULL
	TOX_ERR_SET_INFO_TOO_LONG ToxErrSetInfo = C.TOX_ERR_SET_INFO_TOO_LONG
)

type ToxErrSetTyping C.enum_TOX_ERR_SET_TYPING

const (
	TOX_ERR_SET_TYPING_OK               ToxErrSetTyping = C.TOX_ERR_SET_TYPING_OK
	TOX_ERR_SET_TYPING_FRIEND_NOT_FOUND ToxErrSetTyping = C.TOX_ERR_SET_TYPING_FRIEND_NOT_FOUND
)

type ToxErrFriendSendMessage C.enum_TOX_ERR_FRIEND_SEND_MESSAGE

const (
	TOX_ERR_FRIEND_SEND_MESSAGE_OK                   ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_OK
	TOX_ERR_FRIEND_SEND_MESSAGE_NULL                 ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_NULL
	TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_FOUND     ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_FOUND
	TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_CONNECTED ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_CONNECTED
	TOX_ERR_FRIEND_SEND_MESSAGE_SENDQ                ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_SENDQ
	TOX_ERR_FRIEND_SEND_MESSAGE_TOO_LONG             ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_TOO_LONG
	TOX_ERR_FRIEND_SEND_MESSAGE_EMPTY                ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_EMPTY
)

type ToxErrFriendGetLastOnline C.enum_TOX_ERR_FRIEND_GET_LAST_ONLINE

const (
	TOX_ERR_FRIEND_GET_LAST_ONLINE_OK               ToxErrFriendGetLastOnline = C.TOX_ERR_FRIEND_GET_LAST_ONLINE_OK
	TOX_ERR_FRIEND_GET_LAST_ONLINE_FRIEND_NOT_FOUND ToxErrFriendGetLastOnline = C.TOX_ERR_FRIEND_GET_LAST_ONLINE_FRIEND_NOT_FOUND
)

type ToxErrFileControl C.enum_TOX_ERR_FILE_CONTROL

const (
	TOX_ERR_FILE_CONTROL_OK                   ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_OK
	TOX_ERR_FILE_CONTROL_FRIEND_NOT_FOUND     ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_FRIEND_NOT_FOUND
	TOX_ERR_FILE_CONTROL_FRIEND_NOT_CONNECTED ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_FRIEND_NOT_CONNECTED
	TOX_ERR_FILE_CONTROL_NOT_FOUND            ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_NOT_FOUND
	TOX_ERR_FILE_CONTROL_NOT_PAUSED           ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_NOT_PAUSED
	TOX_ERR_FILE_CONTROL_DENIED               ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_DENIED
	TOX_ERR_FILE_CONTROL_ALREADY_PAUSED       ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_ALREADY_PAUSED
	TOX_ERR_FILE_CONTROL_SENDQ                ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_SENDQ
)

type ToxErrFileSeek C.enum_TOX_ERR_FILE_SEEK

const (
	TOX_ERR_FILE_SEEK_OK                   ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_OK
	TOX_ERR_FILE_SEEK_FRIEND_NOT_FOUND     ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_FRIEND_NOT_FOUND
	TOX_ERR_FILE_SEEK_FRIEND_NOT_CONNECTED ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_FRIEND_NOT_CONNECTED
	TOX_ERR_FILE_SEEK_NOT_FOUND            ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_NOT_FOUND
	TOX_ERR_FILE_SEEK_DENIED               ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_DENIED
	TOX_ERR_FILE_SEEK_INVALID_POSITION     ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_INVALID_POSITION
	TOX_ERR_FILE_SEEK_SENDQ                ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_SENDQ
)

type ToxErrFileGet C.enum_TOX_ERR_FILE_GET

const (
	TOX_ERR_FILE_GET_OK               ToxErrFileGet = C.TOX_ERR_FILE_GET_OK
	TOX_ERR_FILE_GET_NULL             ToxErrFileGet = C.TOX_ERR_FILE_GET_NULL
	TOX_ERR_FILE_GET_FRIEND_NOT_FOUND ToxErrFileGet = C.TOX_ERR_FILE_GET_FRIEND_NOT_FOUND
	TOX_ERR_FILE_GET_NOT_FOUND        ToxErrFileGet = C.TOX_ERR_FILE_GET_NOT_FOUND
)

var (
	ErrFileSendInvalidFileID = errors.New("The size of the given FileID is invalid.")
)

type ToxErrFileSend C.enum_TOX_ERR_FILE_SEND

const (
	TOX_ERR_FILE_SEND_OK                   ToxErrFileSend = C.TOX_ERR_FILE_SEND_OK
	TOX_ERR_FILE_SEND_NULL                 ToxErrFileSend = C.TOX_ERR_FILE_SEND_NULL
	TOX_ERR_FILE_SEND_FRIEND_NOT_FOUND     ToxErrFileSend = C.TOX_ERR_FILE_SEND_FRIEND_NOT_FOUND
	TOX_ERR_FILE_SEND_FRIEND_NOT_CONNECTED ToxErrFileSend = C.TOX_ERR_FILE_SEND_FRIEND_NOT_CONNECTED
	TOX_ERR_FILE_SEND_NAME_TOO_LONG        ToxErrFileSend = C.TOX_ERR_FILE_SEND_NAME_TOO_LONG
	TOX_ERR_FILE_SEND_TOO_MANY             ToxErrFileSend = C.TOX_ERR_FILE_SEND_TOO_MANY
)

type ToxErrFileSendChunk C.enum_TOX_ERR_FILE_SEND_CHUNK

const (
	TOX_ERR_FILE_SEND_CHUNK_OK                   ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_OK
	TOX_ERR_FILE_SEND_CHUNK_NULL                 ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_NULL
	TOX_ERR_FILE_SEND_CHUNK_FRIEND_NOT_FOUND     ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_FRIEND_NOT_FOUND
	TOX_ERR_FILE_SEND_CHUNK_FRIEND_NOT_CONNECTED ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_FRIEND_NOT_CONNECTED
	TOX_ERR_FILE_SEND_CHUNK_NOT_FOUND            ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_NOT_FOUND
	TOX_ERR_FILE_SEND_CHUNK_NOT_TRANSFERRING     ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_NOT_TRANSFERRING
	TOX_ERR_FILE_SEND_CHUNK_INVALID_LENGTH       ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_INVALID_LENGTH
	TOX_ERR_FILE_SEND_CHUNK_SENDQ                ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_SENDQ
	TOX_ERR_FILE_SEND_CHUNK_WRONG_POSITION       ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_WRONG_POSITION
)

type ToxErrGetPort C.enum_TOX_ERR_GET_PORT

const (
	TOX_ERR_GET_PORT_OK        ToxErrGetPort = C.TOX_ERR_GET_PORT_OK
	TOX_ERR_GET_PORT_NOT_BOUND ToxErrGetPort = C.TOX_ERR_GET_PORT_NOT_BOUND
)
