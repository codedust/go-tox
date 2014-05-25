package golibtox

// WIP - organ 2014

/*
#cgo LDFLAGS: -ltoxcore

#include <tox/tox.h>
#include <stdlib.h>

void hook_CallbackFriendRequest(Tox*, uint8_t*, uint8_t*, uint16_t, void*);

static void set_callbackfriendrequest(Tox * t) {
	tox_callback_friend_request(t, hook_CallbackFriendRequest, NULL);
}
*/
import "C"

import (
	"encoding/hex"
	"errors"
	"sync"
	"unsafe"
)

var friendRequestFunc FriendRequestFunc

type Tox struct {
	tox *C.struct_Tox
	mtx sync.Mutex
}

type FriendRequestFunc func(publicKey []byte, data []byte, length uint16)

type Server struct {
	Address string
	Port    uint16
	Key     string
}

type UserStatus C.uint8_t

const (
	USERSTATUS_NONE    UserStatus = C.TOX_USERSTATUS_NONE
	USERSTATUS_AWAY    UserStatus = C.TOX_USERSTATUS_AWAY
	USERSTATUS_BUSY    UserStatus = C.TOX_USERSTATUS_BUSY
	USERSTATUS_INVALID UserStatus = C.TOX_USERSTATUS_INVALID
)

// void tox_callback_friend_request(Tox *tox, void (*function)(Tox *tox, uint8_t *, uint8_t *, uint16_t, void *), void *userdata);

func (t *Tox) CallbackFriendRequest(f FriendRequestFunc) {
	if t.tox != nil {
		friendRequestFunc = f
		C.set_callbackfriendrequest(t.tox)
	}
	return
}

func New() (*Tox, error) {
	ctox := C.tox_new(C.TOX_ENABLE_IPV6_DEFAULT)
	if ctox == nil {
		return nil, errors.New("Error initializing Tox")
	}

	t := &Tox{tox: ctox}

	return t, nil
}

func (t *Tox) Kill() {
	C.tox_kill(t.tox)
}

func (t *Tox) GetAddress() ([]byte, error) {
	if t.tox == nil {
		return nil, errors.New("Error getting address, tox not initialized")
	}

	address := make([]byte, C.TOX_FRIEND_ADDRESS_SIZE)
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

	cname := make([]byte, C.TOX_MAX_NAME_LENGTH)

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
	//int tox_bootstrap_from_address(Tox *tox, const char *address, uint8_t ipv6enabled,
	//                              uint16_t port, uint8_t *public_key);
	if t.tox == nil {
		return errors.New("Tox not initialized")
	}

	caddr := C.CString(s.Address)
	defer C.free(unsafe.Pointer(caddr))

	pubkey, err := s.GetPubKey()

	if err != nil {
		return err
	}

	C.tox_bootstrap_from_address(t.tox, caddr, C.TOX_ENABLE_IPV6_DEFAULT, C.htons((C.uint16_t)(s.Port)), (*C.uint8_t)(&pubkey[0]))

	return nil

}
func (s *Server) GetPubKey() ([]byte, error) {
	pubkey, err := hex.DecodeString(s.Key)
	if err != nil {
		return nil, errors.New("Error decoding server key")
	}
	return pubkey, nil
}

func (t *Tox) IsConnected() (bool, error) {
	if t.tox == nil {
		return false, errors.New("Error getting address, tox not initialized")
	}

	return (C.tox_isconnected(t.tox) == 1), nil
}

/*  return size of messenger data (for saving). */
//uint32_t tox_size(Tox *tox);

func (t *Tox) Size() (uint32, error) {
	if t.tox == nil {
		return 0, errors.New("tox not initialized")
	}

	return (uint32)(C.tox_size(t.tox)), nil
}

/* Save the messenger in data (must be allocated memory of size Messenger_size()). */
//void tox_save(Tox *tox, uint8_t *data);

func (t *Tox) Save() ([]byte, error) {
	if t.tox == nil {
		return nil, errors.New("tox not initialized")
	}
	size, _ := t.Size()

	data := make([]byte, size)
	C.tox_save(t.tox, (*C.uint8_t)(&data[0]))

	return data, nil

}

/* Load the messenger from data of size length.
*
 *  returns 0 on success
  *  returns -1 on failure
*/
//int tox_load(Tox *tox, uint8_t *data, uint32_t length);
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
