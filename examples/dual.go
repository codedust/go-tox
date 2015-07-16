package main

import (
	"encoding/hex"
	"fmt"
	"github.com/codedust/go-tox"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	Address   string
	Port      uint16
	PublicKey []byte
}

func main() {
	alice, err := gotox.New(nil)
	if err != nil {
		panic(err)
	}
	bob, err := gotox.New(nil)
	if err != nil {
		panic(err)
	}

	alice.SelfSetName("AliceBot")
	bob.SelfSetName("BobBot")

	aliceAddr, _ := alice.SelfGetAddress()
	fmt.Println("ID alice: ", hex.EncodeToString(aliceAddr))

	bobAddr, _ := bob.SelfGetAddress()
	fmt.Println("ID bob: ", hex.EncodeToString(bobAddr))

	// We can set the same callback function for both *Tox instances
	bob.CallbackFriendRequest(onFriendRequest)
	bob.CallbackFriendMessage(onFriendMessage)
	bob.CallbackFriendConnectionStatusChanges(onFriendConnectionStatusChanges)
	bob.CallbackSelfConnectionStatusChanges(onSelfConnectionStatusChanges)

	alice.CallbackFriendRequest(onFriendRequest)
	alice.CallbackFriendMessage(onFriendMessage)
	alice.CallbackFriendConnectionStatusChanges(onFriendConnectionStatusChanges)
	alice.CallbackSelfConnectionStatusChanges(onSelfConnectionStatusChanges)

	/* Connect to the network
	 * Use more than one node in a real world szenario. This example relies one
	 * the following node to be up.
	 */
	pubkey, _ := hex.DecodeString("04119E835DF3E78BACF0F84235B300546AF8B936F035185E2A8E9E0A67C8924F")
	server := &Server{"144.76.60.215", 33445, pubkey}

	err = alice.Bootstrap(server.Address, server.Port, server.PublicKey)
	if err != nil {
		panic(err)
	}
	err = bob.Bootstrap(server.Address, server.Port, server.PublicKey)
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
				bob.FriendAdd(aliceAddr, "Hey Alice, wanna be my friend. ;)")
				fmt.Printf("[BobBot] Friend request send. Waiting for Alice to response.\n")
			} else if times == 1 {
				// Then Bob sends a message to Alice
				friendnumbers, _ := bob.SelfGetFriendlist()
				_, err := bob.FriendSendMessage(friendnumbers[0], gotox.TOX_MESSAGE_TYPE_NORMAL, "HELLO ALICE")
				fmt.Printf("[BobBot] Sending message to Alice (friendnumber: %d, error: %v)\n", friendnumbers[0], err)
			} else if times == 2 {
				// Alice responds to Bob
				friendnumbers, _ := alice.SelfGetFriendlist()
				_, err := alice.FriendSendMessage(friendnumbers[0], gotox.TOX_MESSAGE_TYPE_NORMAL, "Hey Bob!")
				fmt.Printf("[AliceBot] Sending message to Bob (friendnumber: %d, error: %v)\n", friendnumbers[0], err)
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
			alice.Iterate()
			bob.Iterate()
			break
		}
	}
}

func onFriendRequest(t *gotox.Tox, publicKey []byte, message string) {
	name, _ := t.SelfGetName()
	fmt.Printf("[%s] New friend request from %s\n", name, hex.EncodeToString(publicKey))

	// Auto-accept friend request
	friendnumber, err := t.FriendAddNorequest(publicKey)
	fmt.Printf("[%s] Friend added (friendnumber: %d, error: %v)\n", name, friendnumber, err)
}

func onFriendMessage(t *gotox.Tox, friendnumber uint32, messageType gotox.ToxMessageType, message string) {
	name, _ := t.SelfGetName()
	friend, _ := t.FriendGetName(friendnumber)
	fmt.Printf("[%s] New message from %s : %s\n", name, friend, message)
}

func onFriendConnectionStatusChanges(t *gotox.Tox, friendnumber uint32, connectionstatus gotox.ToxConnection) {
	name, _ := t.SelfGetName()
	fmt.Printf("[%s] Connection status of friend changed to %v\n", name, connectionstatus)
}

func onSelfConnectionStatusChanges(t *gotox.Tox, connectionstatus gotox.ToxConnection) {
	name, _ := t.SelfGetName()
	fmt.Printf("[%s] Connection status changed to %v\n", name, connectionstatus)
}
