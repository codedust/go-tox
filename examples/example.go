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
	tox.CallbackFileRecvControl(onFileRecvControl)
	tox.CallbackFileChunkRequest(onFileChunkRequest)
	tox.CallbackFileRecv(onFileRecv)
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

func onFriendMessage(t *gotox.Tox, friendNumber uint32, messagetype gotox.ToxMessageType, message string) {
	if messagetype == gotox.TOX_MESSAGE_TYPE_NORMAL {
		fmt.Printf("New message from %d : %s\n", friendNumber, message)
	} else {
		fmt.Printf("New action from %d : %s\n", friendNumber, message)
	}

	switch message {
	case "/help":
		t.FriendSendMessage(friendNumber, gotox.TOX_MESSAGE_TYPE_NORMAL, "Type '/file' to receive a file.")
	case "/file":
		file, err := os.Open("sample_image.jpg")
		if err != nil {
			t.FriendSendMessage(friendNumber, gotox.TOX_MESSAGE_TYPE_NORMAL, "File not found. Please 'cd' into go-tox/examples and type 'go run example.go'")
			file.Close()
			return
		}

		// get the file size
		stat, err := file.Stat()
		if err != nil {
			t.FriendSendMessage(friendNumber, gotox.TOX_MESSAGE_TYPE_NORMAL, "Could not read file stats.")
			file.Close()
			return
		}

		fmt.Println("File size is ", stat.Size())

		fileNumber, err := t.FileSend(friendNumber, gotox.TOX_FILE_KIND_DATA, uint64(stat.Size()), nil, "fileName.jpg")
		if err != nil {
			t.FriendSendMessage(friendNumber, gotox.TOX_MESSAGE_TYPE_NORMAL, "t.FileSend() failed.")
			file.Close()
			return
		}

		transfers[fileNumber] = file
		transfersFilesizes[fileNumber] = uint64(stat.Size())
	default:
		t.FriendSendMessage(friendNumber, gotox.TOX_MESSAGE_TYPE_NORMAL, "Type '/help' for available commands.")
	}
}

func onFileRecv(t *gotox.Tox, friendNumber uint32, fileNumber uint32, kind gotox.ToxFileKind, filesize uint64, filename string) {
	// Accept any file send request
	t.FileControl(friendNumber, true, fileNumber, gotox.TOX_FILE_CONTROL_RESUME, nil)
	// Init *File handle
	f, _ := os.Create("example_" + filename)
	// Append f to the map[uint8]*os.File
	transfers[fileNumber] = f
	transfersFilesizes[fileNumber] = filesize
}

func onFileRecvControl(t *gotox.Tox, friendNumber uint32, fileNumber uint32, fileControl gotox.ToxFileControl) {
	if fileControl == gotox.TOX_FILE_CONTROL_CANCEL {
		// delete (hopefully existing) filehandle
		transfers[fileNumber].Close()
		delete(transfers, fileNumber)
		delete(transfersFilesizes, fileNumber)
	}
}

func onFileChunkRequest(t *gotox.Tox, friendNumber uint32, fileNumber uint32, position uint64, length uint64) {
	// read from the (hopefully existing) filehandle
	if length+position > transfersFilesizes[fileNumber] {
		length = transfersFilesizes[fileNumber] - position
	}

	data := make([]byte, length)
	_, err := transfers[fileNumber].ReadAt(data, int64(position))
	if err != nil {
		fmt.Println("Error reading file", err)
	}

	t.FileSendChunk(friendNumber, fileNumber, position, data)
}

func onFileRecvChunk(t *gotox.Tox, friendNumber uint32, fileNumber uint32, position uint64, data []byte) {
	// write data to the hopefully existing filehandle
	if f, exists := transfers[fileNumber]; exists {
		f.WriteAt(data, (int64)(position))
	}

	// Finished receiving file
	if position == transfersFilesizes[fileNumber] {
		f := transfers[fileNumber]
		f.Sync()
		f.Close()
		delete(transfers, fileNumber)
		delete(transfersFilesizes, fileNumber)
		fmt.Println("Written file", fileNumber)
		t.FriendSendMessage(friendNumber, gotox.TOX_MESSAGE_TYPE_NORMAL, "Thanks!")
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
