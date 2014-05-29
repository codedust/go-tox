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

type Server struct {
	Address   string
	Port      uint16
	PublicKey string
}

// Map of active file transfers
var transfers = make(map[uint8]*os.File)

func main() {
	var filepath string

	flag.StringVar(&filepath, "save", "", "path to save file")
	flag.Parse()

	server := &Server{"37.187.46.132", 33445, "A9D98212B3F972BD11DA52BEB0658C326FCCC1BFD49F347F9C2D3D8B61E1B927"}

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

	// Register our callbacks
	tox.CallbackFriendRequest(onFriendRequest)
	tox.CallbackFriendMessage(onFriendMessage)
	tox.CallbackFileSendRequest(onFileSendRequest)
	tox.CallbackFileControl(onFileControl)
	tox.CallbackFileData(onFileData)

	// Connect to the network
	err = tox.BootstrapFromAddress(server.Address, server.Port, server.PublicKey)
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
			break
		case <-ticker.C:
			tox.Do()
			break
		}
	}
}

func onFriendRequest(t *golibtox.Tox, publicKey []byte, data []byte, length uint16) {
	fmt.Printf("New friend request from %s\n", hex.EncodeToString(publicKey))
	fmt.Printf("With message: %v\n", string(data))
	// Auto-accept friend request
	t.AddFriendNorequest(publicKey)
}

func onFriendMessage(t *golibtox.Tox, friendnumber int32, message []byte, length uint16) {
	fmt.Printf("New message from %d : %s\n", friendnumber, string(message))
	// Echo back
	t.SendMessage(friendnumber, message)
}

func onFileSendRequest(t *golibtox.Tox, friendnumber int32, filenumber uint8, filesize uint64, filename []byte, filenameLength uint16) {
	// Accept any file send request
	t.FileSendControl(friendnumber, true, filenumber, golibtox.FILECONTROL_ACCEPT, nil)
	// Init *File handle
	f, _ := os.Create("example_" + string(filename))
	// Append f to the map[uint8]*os.File
	transfers[filenumber] = f
}

func onFileControl(t *golibtox.Tox, friendnumber int32, sending bool, filenumber uint8, fileControl golibtox.FileControl, data []byte, length uint16) {
	// Finished receiving file
	if fileControl == golibtox.FILECONTROL_FINISHED {
		f := transfers[filenumber]
		f.Sync()
		f.Close()
		delete(transfers, filenumber)
		fmt.Println("Written file", filenumber)
		t.SendMessage(friendnumber, []byte("Thanks!"))
	}
}

func onFileData(t *golibtox.Tox, friendnumber int32, filenumber uint8, data []byte, length uint16) {
	// Write data to the hopefully valid *File handle
	if f, exists := transfers[filenumber]; exists {
		f.Write(data)
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
