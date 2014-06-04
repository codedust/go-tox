package golibtox

/*
#include <tox/tox.h>
*/
import "C"
import "errors"

// General errors
var (
	ErrInit     = errors.New("Error initializing Tox")
	ErrBadTox   = errors.New("Tox not initialized")
	ErrFuncFail = errors.New("Function failed")
	ErrArgs     = errors.New("Nil arguments or wrong size")
)

// Errors returned by AddFriend()
var (
	FaerrTooLong      = errors.New("Message too long")
	FaerrNoMessage    = errors.New("Empty message")
	FaerrOwnKey       = errors.New("Own key")
	FaerrAlreadySent  = errors.New("Already sent")
	FaerrUnkown       = errors.New("Unknown error")
	FaerrBadChecksum  = errors.New("Bad checksum in address")
	FaerrSetNewNospam = errors.New("Different nospam")
	FaerrNoMem        = errors.New("Failed increasing friend list")
)

const (
	MAX_NAME_LENGTH          = C.TOX_MAX_NAME_LENGTH
	MAX_MESSAGE_LENGTH       = C.TOX_MAX_MESSAGE_LENGTH
	MAX_STATUSMESSAGE_LENGTH = C.TOX_MAX_STATUSMESSAGE_LENGTH
	CLIENT_ID_SIZE           = C.TOX_CLIENT_ID_SIZE
	FRIEND_ADDRESS_SIZE      = C.TOX_FRIEND_ADDRESS_SIZE
)

const (
	ENABLE_IPV6_DEFAULT = C.TOX_ENABLE_IPV6_DEFAULT
)

type FriendAddError C.int32_t

const (
	FAERR_TOOLONG      FriendAddError = C.TOX_FAERR_TOOLONG
	FAERR_NOMESSAGE    FriendAddError = C.TOX_FAERR_NOMESSAGE
	FAERR_OWNKEY       FriendAddError = C.TOX_FAERR_OWNKEY
	FAERR_ALREADYSENT  FriendAddError = C.TOX_FAERR_ALREADYSENT
	FAERR_UNKNOWN      FriendAddError = C.TOX_FAERR_UNKNOWN
	FAERR_BADCHECKSUM  FriendAddError = C.TOX_FAERR_BADCHECKSUM
	FAERR_SETNEWNOSPAM FriendAddError = C.TOX_FAERR_SETNEWNOSPAM
	FAERR_NOMEM        FriendAddError = C.TOX_FAERR_NOMEM
)

type UserStatus C.uint8_t

const (
	USERSTATUS_NONE    UserStatus = C.TOX_USERSTATUS_NONE
	USERSTATUS_AWAY    UserStatus = C.TOX_USERSTATUS_AWAY
	USERSTATUS_BUSY    UserStatus = C.TOX_USERSTATUS_BUSY
	USERSTATUS_INVALID UserStatus = C.TOX_USERSTATUS_INVALID
)

type ChatChange C.uint8_t

const (
	CHAT_CHANGE_PEER_ADD  ChatChange = C.TOX_CHAT_CHANGE_PEER_ADD
	CHAT_CHANGE_PEER_DEL  ChatChange = C.TOX_CHAT_CHANGE_PEER_DEL
	CHAT_CHANGE_PEER_NAME ChatChange = C.TOX_CHAT_CHANGE_PEER_NAME
)

type FileControl C.uint8_t

const (
	FILECONTROL_ACCEPT        FileControl = C.TOX_FILECONTROL_ACCEPT
	FILECONTROL_PAUSE         FileControl = C.TOX_FILECONTROL_PAUSE
	FILECONTROL_KILL          FileControl = C.TOX_FILECONTROL_KILL
	FILECONTROL_FINISHED      FileControl = C.TOX_FILECONTROL_FINISHED
	FILECONTROL_RESUME_BROKEN FileControl = C.TOX_FILECONTROL_RESUME_BROKEN
)
