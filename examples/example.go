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

const MAX_AVATAR_SIZE = 65536 // see github.com/Tox/Tox-STS/blob/master/STS.md#avatars

type FileTransfer struct {
	fileHandle *os.File
	fileSize   uint64
}

// Map of active file transfers
var transfers = make(map[uint32]FileTransfer)

func main() {
	var newToxInstance bool = false
	var filepath string
	var options *gotox.Options

	flag.StringVar(&filepath, "save", "./example_savedata", "path to save file")
	flag.Parse()

	fmt.Printf("[INFO] Using Tox version %d.%d.%d\n", gotox.VersionMajor(), gotox.VersionMinor(), gotox.VersionPatch())

	if !gotox.VersionIsCompatible(0, 0, 0) {
		fmt.Println("[ERROR] The compiled library (toxcore) is not compatible with this example.")
		fmt.Println("[ERROR] Please update your Tox library. If this error persists, please report it to the gotox developers.")
		fmt.Println("[ERROR] Thanks!")
		return
	}

	savedata, err := loadData(filepath)
	if err == nil {
		fmt.Println("[INFO] Loading Tox profile from savedata...")
		options = &gotox.Options{
			IPv6Enabled:  true,
			UDPEnabled:   true,
			ProxyType:    gotox.TOX_PROXY_TYPE_NONE,
			ProxyHost:    "127.0.0.1",
			ProxyPort:    5555,
			StartPort:    0,
			EndPort:      0,
			TcpPort:      0, // only enable TCP server if your client provides an option to disable it
			SaveDataType: gotox.TOX_SAVEDATA_TYPE_TOX_SAVE,
			SaveData:     savedata}
	} else {
		fmt.Println("[INFO] Creating new Tox profile...")
		options = nil // default options
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
			fmt.Printf("\nSaving...\n")
			if err := saveData(tox, filepath); err != nil {
				fmt.Println("[ERROR]", err)
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

		transfers[fileNumber] = FileTransfer{fileHandle: file, fileSize: uint64(stat.Size())}
	default:
		t.FriendSendMessage(friendNumber, gotox.TOX_MESSAGE_TYPE_NORMAL, "Type '/help' for available commands.")
	}
}

func onFileRecv(t *gotox.Tox, friendNumber uint32, fileNumber uint32, kind gotox.ToxFileKind, filesize uint64, filename string) {
	if kind == gotox.TOX_FILE_KIND_AVATAR {

		if filesize > MAX_AVATAR_SIZE {
			// reject file send request
			t.FileControl(friendNumber, fileNumber, gotox.TOX_FILE_CONTROL_CANCEL)
			return
		}

		publicKey, _ := t.FriendGetPublickey(friendNumber)
		file, err := os.Create("example_" + hex.EncodeToString(publicKey) + ".png")
		if err != nil {
			fmt.Println("[ERROR] Error creating file", "example_"+hex.EncodeToString(publicKey)+".png")
		}

		// append the file to the map of active file transfers
		transfers[fileNumber] = FileTransfer{fileHandle: file, fileSize: filesize}

		// accept the file send request
		t.FileControl(friendNumber, fileNumber, gotox.TOX_FILE_CONTROL_RESUME)

	} else {
		// accept files of any length
		file, err := os.Create("example_" + filename)
		if err != nil {
			fmt.Println("[ERROR] Error creating file", "example_"+filename)
		}

		// append the file to the map of active file transfers
		transfers[fileNumber] = FileTransfer{fileHandle: file, fileSize: filesize}

		// accept the file send request
		t.FileControl(friendNumber, fileNumber, gotox.TOX_FILE_CONTROL_RESUME)
	}
}

func onFileRecvControl(t *gotox.Tox, friendNumber uint32, fileNumber uint32, fileControl gotox.ToxFileControl) {
	transfer, ok := transfers[fileNumber]
	if !ok {
		fmt.Println("Error: File handle does not exist")
		return
	}

	if fileControl == gotox.TOX_FILE_CONTROL_CANCEL {
		// delete file handle
		transfer.fileHandle.Close()
		delete(transfers, fileNumber)
	}
}

func onFileChunkRequest(t *gotox.Tox, friendNumber uint32, fileNumber uint32, position uint64, length uint64) {
	transfer, ok := transfers[fileNumber]
	if !ok {
		fmt.Println("Error: File handle does not exist")
		return
	}

	// read from the file handle
	if length+position > transfer.fileSize {
		length = transfer.fileSize - position
	}

	if length == 0 {
		transfer.fileHandle.Close()
		delete(transfers, fileNumber)
		fmt.Println("File transfer completed (sending)", fileNumber)
		return
	}

	data := make([]byte, length)
	_, err := transfers[fileNumber].fileHandle.ReadAt(data, int64(position))
	if err != nil {
		fmt.Println("Error reading file", err)
	}

	t.FileSendChunk(friendNumber, fileNumber, position, data)
}

func onFileRecvChunk(t *gotox.Tox, friendNumber uint32, fileNumber uint32, position uint64, data []byte) {
	transfer, ok := transfers[fileNumber]
	if !ok {
		if len(data) == 0 {
			// ignore the zero-length chunk that indicates that the transfer is
			// complete (see below)
			return
		}

		fmt.Println("Error: File handle does not exist")
		return
	}

	// write data to the file handle
	transfer.fileHandle.WriteAt(data, (int64)(position))

	// file transfer completed
	if position+uint64(len(data)) >= transfer.fileSize {
		// Some clients will send us another zero-length chunk without data (only
		// required for stream, not necessary for files with a known size) and some
		// will not.
		// We will delete the file handle now (we aleady reveived the whole file)
		// and ignore the file handle error when the empty chunk arrives.

		transfer.fileHandle.Sync()
		transfer.fileHandle.Close()
		delete(transfers, fileNumber)
		fmt.Println("File transfer completed (receiving)", fileNumber)
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
