package golibtox

// WIP - organ 2014

/*
#cgo LDFLAGS: -ltoxcore

#include <tox/tox.h>
#include <stdlib.h>

// Convenient macro:
// Creates the C function to directly register a given callback
#define HOOK(x) \
static void set_##x(Tox *t) { \
	tox_##x(t, hook_##x, NULL); \
}

void hook_callback_friend_request(Tox*, uint8_t*, uint8_t*, uint16_t, void*);
void hook_callback_friend_message(Tox*, int, uint8_t*, uint16_t, void*);

HOOK(callback_friend_request)
HOOK(callback_friend_message)
*/
import "C"

import (
	"encoding/hex"
	"errors"
	"sync"
	"unsafe"
)

type FriendRequestFunc func(publicKey []byte, data []byte, length uint16)
type FriendMessageFunc func(friendNumber int, message []byte, length uint16)

var friendRequestFunc FriendRequestFunc
var friendMessageFunc FriendMessageFunc

type Tox struct {
	tox *C.struct_Tox
	mtx sync.Mutex
}

type Server struct {
	Address string
	Port    uint16
	Key     string
}

func New() (*Tox, error) {
	ctox := C.tox_new(ENABLE_IPV6_DEFAULT)
	if ctox == nil {
		return nil, errors.New("Error initializing Tox")
	}

	t := &Tox{tox: ctox}

	return t, nil
}

func (t *Tox) Kill() {
	C.tox_kill(t.tox)
}

func (t *Tox) Do() error {
	if t.tox == nil {
		return errors.New("Tox not initialized")
	}

	t.mtx.Lock()
	C.tox_do(t.tox)
	t.mtx.Unlock()

	return nil
}

func (t *Tox) BootstrapFromAddress(s *Server) error {
	if t.tox == nil {
		return errors.New("Tox not initialized")
	}

	caddr := C.CString(s.Address)
	defer C.free(unsafe.Pointer(caddr))

	pubkey, err := s.GetPubKey()

	if err != nil {
		return err
	}

	C.tox_bootstrap_from_address(t.tox, caddr, ENABLE_IPV6_DEFAULT, C.htons((C.uint16_t)(s.Port)), (*C.uint8_t)(&pubkey[0]))

	return nil

}

func (t *Tox) IsConnected() (bool, error) {
	if t.tox == nil {
		return false, errors.New("Error getting address, tox not initialized")
	}

	return (C.tox_isconnected(t.tox) == 1), nil
}

func (t *Tox) GetAddress() ([]byte, error) {
	if t.tox == nil {
		return nil, errors.New("Error getting address, tox not initialized")
	}

	address := make([]byte, FRIEND_ADDRESS_SIZE)
	C.tox_get_address(t.tox, (*C.uint8_t)(&address[0]))

	return address, nil
}

func (t *Tox) SetName(name string) error {
	if t.tox == nil {
		return errors.New("Tox not initialized")
	}
	name += "\x00"

	ret := C.tox_set_name(t.tox, (*C.uint8_t)(&[]byte(name)[0]), (C.uint16_t)(len(name)))
	if ret != 0 {
		return errors.New("Error setting name")
	}
	return nil
}

func (t *Tox) GetSelfName() (string, error) {
	if t.tox == nil {
		return "", errors.New("Tox not initialized")
	}

	cname := make([]byte, MAX_NAME_LENGTH)

	n := C.tox_get_self_name(t.tox, (*C.uint8_t)(&cname[0]))

	name := string(cname[:n])

	return name, nil
}

func (t *Tox) SetUserStatus(status UserStatus) error {
	if t.tox == nil {
		return errors.New("Tox not initialized")
	}

	ret := C.tox_set_user_status(t.tox, (C.uint8_t)(status))
	if ret != 0 {
		return errors.New("Error setting status")
	}
	return nil
}

func (t *Tox) SendMessage(friendNumber int32, message []byte, length uint32) (int32, error) {
	if t.tox == nil {
		return -1, errors.New("Tox not initialized")
	}

	n := C.tox_send_message(t.tox, (C.int32_t)(friendNumber), (*C.uint8_t)(&message[0]), (C.uint32_t)(length))
	if n == 0 {
		return -1, errors.New("Error sending message")
	}
	return (int32)(n), nil
}

func (t *Tox) AddFriendNorequest(clientId []byte) (int32, error) {
	if t.tox == nil {
		return -1, errors.New("Tox not initialized")
	}

	if len(clientId) != CLIENT_ID_SIZE {
		return -1, errors.New("Incorrect client id")
	}

	n := C.tox_add_friend_norequest(t.tox, (*C.uint8_t)(&clientId[0]))
	if n == -1 {
		return -1, errors.New("Error adding friend")
	}
	return (int32)(n), nil
}

func (t *Tox) Size() (uint32, error) {
	if t.tox == nil {
		return 0, errors.New("tox not initialized")
	}

	return (uint32)(C.tox_size(t.tox)), nil
}

func (t *Tox) Save() ([]byte, error) {
	if t.tox == nil {
		return nil, errors.New("tox not initialized")
	}
	size, _ := t.Size()

	data := make([]byte, size)
	C.tox_save(t.tox, (*C.uint8_t)(&data[0]))

	return data, nil

}

func (t *Tox) Load(data []byte, length uint32) error {
	if t.tox == nil {
		return errors.New("tox not initialized")
	}

	ret := C.tox_load(t.tox, (*C.uint8_t)(&data[0]), (C.uint32_t)(length))

	if ret == -1 {
		return errors.New("Error loading data")
	}
	return nil
}

func (s *Server) GetPubKey() ([]byte, error) {
	pubkey, err := hex.DecodeString(s.Key)
	if err != nil {
		return nil, errors.New("Error decoding server key")
	}
	return pubkey, nil
}

func (t *Tox) CallbackFriendRequest(f FriendRequestFunc) {
	if t.tox != nil {
		friendRequestFunc = f
		C.set_callback_friend_request(t.tox)
	}
}

func (t *Tox) CallbackFriendMessage(f FriendMessageFunc) {
	if t.tox != nil {
		friendMessageFunc = f
		C.set_callback_friend_message(t.tox)
	}
}
