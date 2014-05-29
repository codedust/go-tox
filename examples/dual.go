package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/organ/golibtox"
)

type Server struct {
	Address   string
	Port      uint16
	PublicKey string
}

func main() {

	server := &Server{"37.187.46.132", 33445, "A9D98212B3F972BD11DA52BEB0658C326FCCC1BFD49F347F9C2D3D8B61E1B927"}

	alice, err := golibtox.New()
	if err != nil {
		panic(err)
	}
	bob, err := golibtox.New()
	if err != nil {
		panic(err)
	}

	alice.SetName("AliceBot")
	bob.SetName("BobBot")

	aliceAddr, _ := alice.GetAddress()
	fmt.Println("ID alice: ", hex.EncodeToString(aliceAddr))

	bobAddr, _ := bob.GetAddress()
	fmt.Println("ID bob: ", hex.EncodeToString(bobAddr))

	// We can set the same callback function for both *Tox instances
	alice.CallbackFriendRequest(onFriendRequest)
	bob.CallbackFriendRequest(onFriendRequest)

	alice.CallbackFriendMessage(onFriendMessage)
	bob.CallbackFriendMessage(onFriendMessage)

	err = alice.BootstrapFromAddress(server.Address, server.Port, server.PublicKey)
	if err != nil {
		panic(err)
	}
	err = bob.BootstrapFromAddress(server.Address, server.Port, server.PublicKey)
	if err != nil {
		panic(err)
	}

	isRunning := true

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ticker := time.NewTicker(25 * time.Millisecond)

	times := 0
	for isRunning {
		select {
		case <-c:
			// Press ^C to trigger those events
			if times == 0 {
				// First Bob adds Alice
				bob.AddFriend(aliceAddr, []byte("o"))
			} else if times == 1 {
				// Then Bob sends a message to Alice
				bob.SendMessage(0, []byte("HELLO ALICE"))
			} else if times == 2 {
				// Alice responds to Bob
				alice.SendMessage(0, []byte("Hey Bob !"))
			} else {
				// We then put an end to their love
				fmt.Println("Killing")
				isRunning = false
				alice.Kill()
				bob.Kill()
			}
			times += 1
			break
		case <-ticker.C:
			alice.Do()
			bob.Do()
			break
		}
	}
}

func onFriendRequest(t *golibtox.Tox, publicKey []byte, data []byte, length uint16) {
	name, _ := t.GetSelfName()
	fmt.Printf("[%s] New friend request from %s\n", name, hex.EncodeToString(publicKey))

	// Auto-accept friend request
	clientId := publicKey[:golibtox.CLIENT_ID_SIZE]
	t.AddFriendNorequest(clientId)
}

func onFriendMessage(t *golibtox.Tox, friendnumber int32, message []byte, length uint16) {
	name, _ := t.GetSelfName()
	friend, _ := t.GetName(friendnumber)
	fmt.Printf("[%s] New message from %s : %s\n", name, friend, string(message))
}
