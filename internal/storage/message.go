package storage

import (
	"bytes"
	"encoding/gob"
	"os"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/charlie1404/vqs/internal/utils"
)

type Attribute struct {
	DataType string
	Value    []byte
}

type Message struct {
	Id               string
	Body             string
	Attributes       map[string]Attribute
	SystemAttributes map[string]Attribute
	DelaySeconds     uint16
}

// Some black magic happens here, follow very closely to understand
func getMessagesMmapedRegion(path string) (*[]byte, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fd := -1
	memoryMapSize := DATA_BUFFER_SIZE + MMAP_BUFFER_SIZE

	mem, _, err := syscall.Syscall6(
		syscall.SYS_MMAP,
		0,                                    // null address, let kernel choose the address
		uintptr(memoryMapSize),               // size of the memory region, this includes the wrapping buffer too
		syscall.PROT_NONE,                    // will be overlapped with the actual mmap
		syscall.MAP_PRIVATE|syscall.MAP_ANON, // no real file backing this memory region, and private to this process
		uintptr(fd),                          // MAP_ANON requires fd to be -1
		0,                                    // offset
	)
	if err != syscall.Errno(0) {
		return nil, err
	}

	_, _, err = syscall.Syscall6(
		syscall.SYS_MMAP,
		mem,                                  // address of the memory region, main storage
		uintptr(DATA_BUFFER_SIZE),            // size of the memory region, size of the file
		syscall.PROT_READ|syscall.PROT_WRITE, // read and write access
		syscall.MAP_SHARED|syscall.MAP_FIXED, // place the mapping at exactly at start of the previously allocated memory region
		file.Fd(),                            // this is backed by data file
		0,                                    // offset
	)
	if err != syscall.Errno(0) {
		return nil, err
	}

	_, _, err = syscall.Syscall6(
		syscall.SYS_MMAP,
		mem+uintptr(DATA_BUFFER_SIZE),        // address of the memory region, wrapping buffer
		uintptr(MMAP_BUFFER_SIZE),            // size of the memory region, size of the wrapping buffer
		syscall.PROT_READ|syscall.PROT_WRITE, // read and write access
		syscall.MAP_SHARED|syscall.MAP_FIXED, // place the mapping at exactly after the main storage, (but wrap the data from beginning)
		file.Fd(),                            // this is backed by same data file, used above
		0,                                    // offset
	)
	if err != syscall.Errno(0) {
		return nil, err
	}

	var byteArray []byte
	dh := (*reflect.SliceHeader)(unsafe.Pointer(&byteArray))
	dh.Data = mem
	dh.Len = int(memoryMapSize)
	dh.Cap = int(memoryMapSize)

	return &byteArray, nil
}

func serializeMessage(m interface{}) ([]byte, int, error) {
	gobBuffer := new(bytes.Buffer)

	enc := gob.NewEncoder(gobBuffer)
	enc.Encode(m)

	return gobBuffer.Bytes(), gobBuffer.Len(), nil
}

func deserializeMessage(b []byte) (*Message, error) {
	// for debugging purposes use json encoding instead of gob+gzip

	gobBuffer := bytes.NewBuffer(b)

	dec := gob.NewDecoder(gobBuffer)
	var m Message
	err := dec.Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func NewMessage(delaySeconds uint16, body string) *Message {
	return &Message{
		Id:               utils.GenerateUUIDLikeId(),
		DelaySeconds:     delaySeconds,
		Body:             body,
		Attributes:       make(map[string]Attribute),
		SystemAttributes: make(map[string]Attribute),
	}
}

// _, _, e1 := syscall.Syscall(
// 	syscall.SYS_MADVISE,
// 	uintptr(unsafe.Pointer(&byteArray[HEADER_SIZE])),
// 	uintptr(INITIAL_QUEUE_FILE_SIZE-HEADER_SIZE),
// 	uintptr(syscall.MADV_SEQUENTIAL),
// )
// if e1 != 0 {
// 	err = e1
// }
