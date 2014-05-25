package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
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

	err = tox.SetUserStatus(golibtox.USERSTATUS_NONE)

	tox.CallbackFriendRequest(func(pubkey []byte, data []byte, length uint16) {
		fmt.Printf("New friend request from %v\n", pubkey)
		fmt.Printf("With message: %v\n", string(data))

		// Auto-accept friend request
		clientId := pubkey[:golibtox.CLIENT_ID_SIZE]
		fmt.Println(tox.AddFriendNorequest(clientId))
	})

	tox.CallbackFriendMessage(func(friendId int, message []byte, length uint16) {
		fmt.Printf("New message from %d : %s\n", friendId, string(message))
		fmt.Println(tox.SendMessage((int32)(friendId), message, (uint32)(length)))
	})

	saveData(tox)

	err = tox.BootstrapFromAddress(server)
	if err != nil {
		panic(err)
	}

	isRunning := true

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for {
			select {
			case <-c:
				fmt.Println("Saving...")
				if err := saveData(tox); err != nil {
					fmt.Println(err)
				}
				fmt.Println("Killing...")
				isRunning = false
				tox.Kill()
				break
			case <-time.After(time.Second * 10):
				connected, _ := tox.IsConnected()
				fmt.Println("IsConnected() =>", connected)
			}
		}
	}()

	for isRunning {
		tox.Do()
		time.Sleep(25 * time.Millisecond)
	}
}
func saveData(t *golibtox.Tox) error {
	var err error
	var data []byte
	if len(dataPath) > 0 {
		data, err = t.Save()
		err = ioutil.WriteFile(dataPath, data, 0644)
	}
	return err
}
