package golibtox

/*
#cgo LDFLAGS: -ltoxcore

#include <tox/tox.h>
#include <stdlib.h>
*/
import "C"

import "errors"

const BOOTSTRAP_ADDRESS string = "37.187.46.132"
const BOOTSTRAP_PORT int = 33445
const BOOTSTRAP_KEY string = "A9D98212B3F972BD11DA52BEB0658C326FCCC1BFD49F347F9C2D3D8B61E1B927"

type Tox struct {
	tox *C.struct_Tox
}

func New() (*Tox, error) {
	ctox := C.tox_new(C.TOX_ENABLE_IPV6_DEFAULT)
	if ctox == nil {
		return nil, errors.New("Error initializing Tox")
	}

	t := &Tox{tox: ctox}

	return t, nil
}

func (t *Tox) GetAddress() ([]byte, error) {
	if t.tox == nil {
		return nil, errors.New("Error getting address, tox not initialized")
	}

	var caddress [C.TOX_FRIEND_ADDRESS_SIZE]C.uint8_t
	C.tox_get_address(t.tox, &caddress[0])

	address := make([]byte, C.TOX_FRIEND_ADDRESS_SIZE)
	for i, v := range caddress {
		address[i] = (byte)(v)
	}

	return address, nil
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
