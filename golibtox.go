package golibtox

/*
#cgo LDFLAGS: -ltoxcore

#include <tox/tox.h>
#include <stdlib.h>
*/
import "C"

import (
	"encoding/hex"
	"errors"
	"unsafe"
)

const BOOTSTRAP_ADDRESS string = "37.187.46.132"
const BOOTSTRAP_PORT int = 33445
const BOOTSTRAP_KEY string = "A9D98212B3F972BD11DA52BEB0658C326FCCC1BFD49F347F9C2D3D8B61E1B927"

type Tox struct {
	tox *C.struct_Tox
}

type Server struct {
	Address string
	Port    int
	Key     string
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

	//var caddress [C.TOX_FRIEND_ADDRESS_SIZE]C.uint8_t
	var caddress [C.TOX_FRIEND_ADDRESS_SIZE]byte
	C.tox_get_address(t.tox, (*C.uint8_t)(&caddress[0]))

	address := make([]byte, C.TOX_FRIEND_ADDRESS_SIZE)
	for i, v := range caddress {
		address[i] = (byte)(v)
	}

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

func (t *Tox) Connect(s Server) error {
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

	C.tox_bootstrap_from_address(t.tox, caddr, C.TOX_ENABLE_IPV6_DEFAULT, (C.uint16_t)(s.Port), (*C.uint8_t)(&pubkey[0]))

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

func (t *Tox) Do() {
	//	cbuffer := (*C.uint8_t)(C.malloc(C.TOX_FRIEND_ADDRESS_SIZE))
	//	defer C.free(unsafe.Pointer(cbuffer))

}
