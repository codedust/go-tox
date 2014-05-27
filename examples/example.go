package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/organ/golibtox"
)

func main() {
	var filepath string

	flag.StringVar(&filepath, "save", "", "path to save file")
	flag.Parse()

	tox, err := golibtox.New()
	if err != nil {
		panic(err)
	}

	// If no data could be loaded, we should set the name
	if err := loadData(tox, filepath); err != nil {
		tox.SetName("GolibtoxBot")
	}

	tox.SetStatusMessage([]byte("golibtox is cool!"))

	addr, _ := tox.GetAddress()
	fmt.Println("ID: ", hex.EncodeToString(addr))

	err = tox.SetUserStatus(golibtox.USERSTATUS_NONE)

	tox.CallbackFriendRequest(func(pubkey []byte, data []byte, length uint16) {
		fmt.Printf("New friend request from %s\n", hex.EncodeToString(pubkey))
		fmt.Printf("With message: %v\n", string(data))

		// Auto-accept friend request
		clientId := pubkey[:golibtox.CLIENT_ID_SIZE]
		tox.AddFriendNorequest(clientId)
	})

	tox.CallbackFriendMessage(func(friendNumber int32, message []byte, length uint16) {
		fmt.Printf("New message from %d : %s\n", friendNumber, string(message))
		tox.SendMessage(friendNumber, message)
		friendName, _ := tox.GetName(friendNumber)
		greetings := fmt.Sprintf("thinks %s is cool.", friendName)
		tox.SendAction(friendNumber, []byte(greetings))
	})

	server := &golibtox.Server{"37.187.46.132", 33445, "A9D98212B3F972BD11DA52BEB0658C326FCCC1BFD49F347F9C2D3D8B61E1B927"}
	//server := &golibtox.Server{"192.254.75.98", 33445, "951C88B7E75C867418ACDB5D273821372BB5BD652740BCDF623A4FA293E75D2F"}

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
				if err := saveData(tox, filepath); err != nil {
					fmt.Println(err)
				}
				fmt.Println("Killing")
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

func loadData(t *golibtox.Tox, filepath string) error {
	if len(filepath) == 0 {
		return errors.New("Empty path")
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = t.Load(data)

	return err
}

func saveData(t *golibtox.Tox, filepath string) error {
	if len(filepath) == 0 {
		return errors.New("Empty path")
	}

	data, err := t.Save()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath, data, 0644)
	return err
}
