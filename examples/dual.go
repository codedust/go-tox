package main

import (
	"encoding/hex"
	"fmt"
	"github.com/codedust/go-tox"
	"time"
)

type Server struct {
	Address   string
	Port      uint16
	PublicKey []byte
}

var counter int = 0

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
	fmt.Println("[ID alice]", hex.EncodeToString(aliceAddr))

	bobAddr, _ := bob.SelfGetAddress()
	fmt.Println("[ID bob]", hex.EncodeToString(bobAddr))

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
	pubkey, _ := hex.DecodeString("B75583B6D967DB8AD7C6D3B6F9318194BCC79B2FEF18F69E2DF275B779E7AA30")
	server := &Server{"maggie.prok.pw", 33445, pubkey}

	err = alice.Bootstrap(server.Address, server.Port, server.PublicKey)
	if err != nil {
		panic(err)
	}
	err = bob.Bootstrap(server.Address, server.Port, server.PublicKey)
	if err != nil {
		panic(err)
	}

	isRunning := true

	ticker := time.NewTicker(25 * time.Millisecond)

	for isRunning {
		select {
		case <-ticker.C:
			if counter == 2 {
				time.Sleep(2 * time.Second)
				// First Bob adds Alice
				bob.FriendAdd(aliceAddr, "Hey Alice, wanna be my friend. ;)")
				fmt.Printf("[BobBot] Friend request send. Waiting for Alice to respond.\n")
				counter++
			} else if counter == 6 {
				time.Sleep(2 * time.Second)
				// Then Bob sends a message to Alice
				friendnumbers, _ := bob.SelfGetFriendlist()
				_, err := bob.FriendSendMessage(friendnumbers[0], gotox.TOX_MESSAGE_TYPE_NORMAL, "HELLO ALICE")
				fmt.Printf("[BobBot] Sending message to Alice (friendnumber: %d, error: %v)\n", friendnumbers[0], err)
				counter++
			} else if counter == 8 {
				time.Sleep(2 * time.Second)
				// Alice responds to Bob
				friendnumbers, _ := alice.SelfGetFriendlist()
				_, err := alice.FriendSendMessage(friendnumbers[0], gotox.TOX_MESSAGE_TYPE_NORMAL, "Hey Bob!")
				fmt.Printf("[AliceBot] Sending message to Bob (friendnumber: %d, error: %v)\n", friendnumbers[0], err)
				counter++
			} else if counter == 10 {
				time.Sleep(2 * time.Second)
				// We then put an end to their love
				fmt.Println("\\o/ It worked! Killing...")
				isRunning = false
				alice.Kill()
				bob.Kill()
				break
			}
			alice.Iterate()
			bob.Iterate()
			break
		}
	}
}

func onFriendRequest(t *gotox.Tox, publicKey []byte, message string) {
	counter++
	name, _ := t.SelfGetName()
	fmt.Printf("[%s] New friend request from %s\n", name, hex.EncodeToString(publicKey))

	// Auto-accept friend request
	friendnumber, err := t.FriendAddNorequest(publicKey)
	fmt.Printf("[%s] Friend added (friendnumber: %d, error: %v)\n", name, friendnumber, err)
}

func onFriendMessage(t *gotox.Tox, friendnumber uint32, messageType gotox.ToxMessageType, message string) {
	counter++
	name, _ := t.SelfGetName()
	friend, _ := t.FriendGetName(friendnumber)
	fmt.Printf("[%s] New message from %s : %s\n", name, friend, message)
}

func onFriendConnectionStatusChanges(t *gotox.Tox, friendnumber uint32, connectionstatus gotox.ToxConnection) {
	counter++
	name, _ := t.SelfGetName()
	fmt.Printf("[%s] Connection status of friend changed to %v\n", name, connectionstatus)
}

func onSelfConnectionStatusChanges(t *gotox.Tox, connectionstatus gotox.ToxConnection) {
	counter++
	name, _ := t.SelfGetName()
	fmt.Printf("[%s] Connection status changed to %v\n", name, connectionstatus)
}
