package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/organ/golibtox"
)

var dataPath string
var dataLoaded bool

func main() {
	flag.StringVar(&dataPath, "save", "", "path to save file")
	flag.Parse()

	tox, err := golibtox.New()
	if err != nil {
		panic(err)
	}

	if len(dataPath) > 0 {
		data, err := ioutil.ReadFile(dataPath)
		if err != nil {
			fmt.Println(err)
		} else {
			err = tox.Load(data, (uint32)(len(data)))
			if err != nil {
				fmt.Println(err)
			}
			dataLoaded = true
		}
	}

	server := &golibtox.Server{"37.187.46.132", 33445, "A9D98212B3F972BD11DA52BEB0658C326FCCC1BFD49F347F9C2D3D8B61E1B927"}
	//server := &golibtox.Server{"192.254.75.98", 33445, "951C88B7E75C867418ACDB5D273821372BB5BD652740BCDF623A4FA293E75D2F"}

	if !dataLoaded {
		tox.SetName("GolibtoxBot")
	}

	badr, _ := tox.GetAddress()
	fmt.Printf("ID: ")
	for _, v := range badr {
		fmt.Printf("%02X", v)
	}
	fmt.Println()

	err = tox.SetUserStatus(golibtox.USERSTATUS_BUSY)

	if len(dataPath) > 0 {
		data, err := tox.Save()
		err = ioutil.WriteFile(dataPath, data, 0644)
		if err != nil {
			panic(err)
		}
	}

	tox.CallbackFriendRequest(func(pubkey []byte, data []byte, length uint16) {
		fmt.Println("SHBLAH")
		fmt.Println("%v", pubkey)
		fmt.Println("%v", data)
		fmt.Println("%d", length)
	})

	err = tox.BootstrapFromAddress(server)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			connected, _ := tox.IsConnected()
			fmt.Println("IsConnected() =>", connected)
			time.Sleep(2 * time.Second)
		}
	}()

	for {
		tox.Do()
		time.Sleep(25 * time.Millisecond)
	}
	tox.Kill()
}
