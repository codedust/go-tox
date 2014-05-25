package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/organ/golibtox"
)

func main() {
	tox, err := golibtox.New()
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile("/home/organ/.config/tox/data")

	tox.Load(data, (uint32)(len(data)))

	adr, err := tox.GetAddress()
	fmt.Println(adr)

	connected, err := tox.IsConnected()
	fmt.Println(connected)

	//server := golibtox.Server{"37.187.46.132", 33445, "A9D98212B3F972BD11DA52BEB0658C326FCCC1BFD49F347F9C2D3D8B61E1B927"}
	server := &golibtox.Server{"192.254.75.98", 33445, "951C88B7E75C867418ACDB5D273821372BB5BD652740BCDF623A4FA293E75D2F"}

	pubkey, err := server.GetPubKey()
	fmt.Println(pubkey)

	//tox.SetName("Coucou")
	fmt.Println(tox.GetSelfName())

	err = tox.BootstrapFromAddress(server)
	if err != nil {
		panic(err)
	}

	fmt.Println(tox.Size())

	go func() {
		for {
			connected, err = tox.IsConnected()
			fmt.Println(connected)
			time.Sleep(3 * time.Second)
		}
	}()

	for {
		tox.Do()
		time.Sleep(10 * time.Millisecond)
	}
	tox.Kill()
}
