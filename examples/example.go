package main

import (
	"fmt"

	"github.com/organ/golibtox"
)

func main() {
	tox, err := golibtox.New()
	if err != nil {
		panic(err)
	}

	adr, err := tox.GetAddress()
	fmt.Println(adr)

	connected, err := tox.IsConnected()
	fmt.Println(connected)

	server := golibtox.Server{"37.187.46.132", 334455, "A9D98212B3F972BD11DA52BEB0658C326FCCC1BFD49F347F9C2D3D8B61E1B927"}

	pubkey, err := server.GetPubKey()
	fmt.Println(pubkey)

	tox.SetName("Coucou")
	fmt.Println(tox.GetSelfName())

	err = tox.Connect(server)
	if err != nil {
		panic(err)
	}

	connected, err = tox.IsConnected()
	fmt.Println(connected)

	tox.Kill()
}
