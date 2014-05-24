package golibtox

/*
#cgo LDFLAGS: -ltoxcore

#include <tox/tox.h>
#include <stdlib.h>
*/
import "C"

import "fmt"

const BOOTSTRAP_ADDRESS string = "37.187.46.132"
const BOOTSTRAP_PORT int = 33445
const BOOTSTRAP_KEY string = "A9D98212B3F972BD11DA52BEB0658C326FCCC1BFD49F347F9C2D3D8B61E1B927"

func Do() {
	// void tox_get_address(Tox *tox, uint8_t *address);
	fmt.Printf("Tox address size : %d\n", C.TOX_CLIENT_ID_SIZE)

	// Init a Tox
	var tox *C.struct_Tox
	tox = C.tox_new(0)

	fmt.Printf("Tox connected : %d\n", C.tox_isconnected(tox))

	//	cbuffer := (*C.uint8_t)(C.malloc(C.TOX_FRIEND_ADDRESS_SIZE))
	//	defer C.free(unsafe.Pointer(cbuffer))

	//var myid = make([]C.uint8_t, C.TOX_FRIEND_ADDRESS_SIZE)
	var myid [C.TOX_FRIEND_ADDRESS_SIZE]C.uint8_t

	C.tox_get_address(tox, &myid[0])

	fmt.Println("ID:")
	for _, v := range myid {
		fmt.Printf("%02X", v)
	}
	fmt.Println()
}
