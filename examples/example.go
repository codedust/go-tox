package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/codedust/go-tox"
	"io/ioutil"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	Address   string
	Port      uint16
	PublicKey []byte
}

// Map of active file transfers
var transfers = make(map[uint32]*os.File)
var transfersFilesizes = make(map[uint32]uint64)

func main() {
	var newToxInstance bool = false
	var filepath string
	var options *gotox.Options

	flag.StringVar(&filepath, "save", "", "path to save file")
	flag.Parse()

	savedata, err := loadData(filepath)
	if err == nil {
		options = &gotox.Options{
			true, true,
			gotox.TOX_PROXY_TYPE_NONE, "127.0.0.1", 5555, 0, 0,
			0, // local TCP server is disabled. Only enable it if your client provides
			   // an option to disable it.
			gotox.TOX_SAVEDATA_TYPE_TOX_SAVE, savedata}
	} else {
		options = &gotox.Options{
			true, true,
			gotox.TOX_PROXY_TYPE_NONE, "127.0.0.1", 5555, 0, 0,
			0,
			gotox.TOX_SAVEDATA_TYPE_NONE, nil}
		newToxInstance = true
	}

	tox, err := gotox.New(options)
	if err != nil {
		panic(err)
	}

	if newToxInstance {
		tox.SelfSetName("gotoxBot")
		tox.SelfSetStatusMessage("gotox is cool!")
	}

	addr, _ := tox.SelfGetAddress()
	fmt.Println("ID: ", hex.EncodeToString(addr))

	err = tox.SelfSetStatus(gotox.TOX_USERSTATUS_NONE)

	// Register our callbacks
	tox.CallbackFriendRequest(onFriendRequest)
	tox.CallbackFriendMessage(onFriendMessage)
	tox.CallbackFileRecv(onFileRecv)
	tox.CallbackFileRecvControl(onFileRecvControl)
	tox.CallbackFileRecvChunk(onFileRecvChunk)

	/* Connect to the network
	 * Use more than one node in a real world szenario. This example relies one
	 * the following node to be up.
	 */
	pubkey, _ := hex.DecodeString("04119E835DF3E78BACF0F84235B300546AF8B936F035185E2A8E9E0A67C8924F")
	server := &Server{"144.76.60.215", 33445, pubkey}

	err = tox.Bootstrap(server.Address, server.Port, server.PublicKey)
	if err != nil {
		panic(err)
	}

	isRunning := true

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ticker := time.NewTicker(25 * time.Millisecond)

	for isRunning {
		select {
		case <-c:
			fmt.Println("Saving...")
			if err := saveData(tox, filepath); err != nil {
				fmt.Println(err)
			}
			fmt.Println("Killing")
			isRunning = false
			tox.Kill()
		case <-ticker.C:
			tox.Iterate()
		}
	}
}

func onFriendRequest(t *gotox.Tox, publicKey []byte, message string) {
	fmt.Printf("New friend request from %s\n", hex.EncodeToString(publicKey))
	fmt.Printf("With message: %v\n", message)
	// Auto-accept friend request
	t.FriendAddNorequest(publicKey)
}

func onFriendMessage(t *gotox.Tox, friendnumber uint32, messagetype gotox.ToxMessageType, message string) {
	if messagetype == gotox.TOX_MESSAGE_TYPE_NORMAL {
		fmt.Printf("New message from %d : %s\n", friendnumber, message)
	} else {
		fmt.Printf("New action from %d : %s\n", friendnumber, message)
	}

	// Echo back
	t.FriendSendMessage(friendnumber, messagetype, message)
}

func onFileRecv(t *gotox.Tox, friendnumber uint32, filenumber uint32, kind gotox.ToxFileKind, filesize uint64, filename string) {
	// Accept any file send request
	t.FileControl(friendnumber, true, filenumber, gotox.TOX_FILE_CONTROL_RESUME, nil)
	// Init *File handle
	f, _ := os.Create("example_" + filename)
	// Append f to the map[uint8]*os.File
	transfers[filenumber] = f
	transfersFilesizes[filenumber] = filesize
}

func onFileRecvControl(t *gotox.Tox, friendnumber uint32, filenumber uint32, fileControl gotox.ToxFileControl) {
	// Do something useful
}

func onFileRecvChunk(t *gotox.Tox, friendnumber uint32, filenumber uint32, position uint64, data []byte) {
	// Write data to the hopefully valid *File handle
	if f, exists := transfers[filenumber]; exists {
		f.WriteAt(data, (int64)(position))
	}

	// Finished receiving file
	if position == transfersFilesizes[filenumber] {
		f := transfers[filenumber]
		f.Sync()
		f.Close()
		delete(transfers, filenumber)
		fmt.Println("Written file", filenumber)
		t.FriendSendMessage(friendnumber, gotox.TOX_MESSAGE_TYPE_NORMAL, "Thanks!")
	}
}

func loadData(filepath string) ([]byte, error) {
	if len(filepath) == 0 {
		return nil, errors.New("Empty path")
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return data, err
}

func saveData(t *gotox.Tox, filepath string) error {
	if len(filepath) == 0 {
		return errors.New("Empty path")
	}

	data, err := t.GetSavedata()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath, data, 0644)
	return err
}
