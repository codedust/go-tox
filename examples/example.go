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

	tox.Do()
}
