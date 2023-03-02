package storage

import (
	"encoding/binary"
	"errors"
	"os"
	"path"
	"reflect"
	"sync"
	"syscall"
	"unsafe"

	"github.com/charlie1404/vqs/pkg/o11y/logs"
)

type Meta struct {
	MagicNum               uint32
	CapacityLeft           uint32
	Count                  uint32
	ReadOffset             uint32
	WriteOffset            uint32
	MaxMessageSize         uint32
	MessageRetentionPeriod uint32
	MessageWaitSeconds     uint16
	DelaySeconds           uint16
	VisibilityTimeout      uint16
	HeaderMetaMarker       uint16 // TODO: this is not used yet
}

type Queue struct {
	meta     *Meta
	messages []byte
	msgMutex sync.RWMutex // RLock can only be used for peek cases as of now
	// inFlightMessages []byte // for now we will keep these in memory, later put on disk
	// inFlightMutex    sync.Mutex
	// delayedMessages  []byte // for now we will keep these in memory, later put on disk
	// delayedMutex     sync.Mutex
}

func createQueueFolderExists(queueDir string) error {
	err := os.RemoveAll(queueDir)
	if err != nil {
		logs.Logger.Warn().Err(err).Msg("error removing queue dir")
	}

	if err := os.MkdirAll(queueDir, 0700); err != nil {
		logs.Logger.Error().Err(err).Msg("error creating queue dir")
		return err
	}

	return nil
}

func createQueueDataFilesExists(queueDirPath string) error {
	dataFile := path.Join(queueDirPath, "data.dat")
	inFlightFile := path.Join(queueDirPath, "in_flight.dat")
	delayedFile := path.Join(queueDirPath, "delayed.dat")

	_ = os.Remove(dataFile)
	_ = os.Remove(inFlightFile)
	_ = os.Remove(delayedFile)

	file, err := os.OpenFile(dataFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := syscall.Truncate(dataFile, INITIAL_QUEUE_FILE_SIZE); err != nil {
		os.Remove(dataFile)
		return err
	}

	// ===============================

	file, err = os.OpenFile(inFlightFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := syscall.Truncate(inFlightFile, INITIAL_QUEUE_FILE_SIZE); err != nil {
		os.Remove(dataFile)
		return err
	}

	// ===============================

	file, err = os.OpenFile(delayedFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := syscall.Truncate(delayedFile, INITIAL_QUEUE_FILE_SIZE); err != nil {
		os.Remove(dataFile)
		return err
	}

	return nil
}

// Black magic happens here, follow very closely to understand
func getMmapedRegion(path string) ([]byte, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// should not occur, add custom error
		}
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fd := -1
	memoryMapSize := INITIAL_QUEUE_FILE_SIZE + MMAP_BUFFER_SIZE

	mem, _, err := syscall.Syscall6(syscall.SYS_MMAP, 0, uintptr(memoryMapSize), syscall.PROT_NONE, syscall.MAP_PRIVATE|syscall.MAP_ANON, uintptr(fd), 0)
	if err != syscall.Errno(0) {
		return nil, err
	}

	_, _, err = syscall.Syscall6(syscall.SYS_MMAP, mem, uintptr(INITIAL_QUEUE_FILE_SIZE), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_FIXED, file.Fd(), 0)
	if err != syscall.Errno(0) {
		return nil, err
	}

	_, _, err = syscall.Syscall6(syscall.SYS_MMAP, mem+INITIAL_QUEUE_FILE_SIZE, uintptr(MMAP_BUFFER_SIZE), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_FIXED, file.Fd(), HEADER_SIZE)
	if err != syscall.Errno(0) {
		return nil, err
	}

	var byteArray []byte
	dh := (*reflect.SliceHeader)(unsafe.Pointer(&byteArray))
	dh.Data = mem
	dh.Len = INITIAL_QUEUE_FILE_SIZE + MMAP_BUFFER_SIZE
	dh.Cap = INITIAL_QUEUE_FILE_SIZE + MMAP_BUFFER_SIZE

	// _, _, e1 := syscall.Syscall(
	// 	syscall.SYS_MADVISE,
	// 	uintptr(unsafe.Pointer(&byteArray[HEADER_SIZE])),
	// 	uintptr(INITIAL_QUEUE_FILE_SIZE-HEADER_SIZE),
	// 	uintptr(syscall.MADV_SEQUENTIAL),
	// )
	// if e1 != 0 {
	// 	err = e1
	// }

	return byteArray, nil
}

func getQueue(queueDirPath string) (*Queue, error) {
	dataFile := path.Join(queueDirPath, "data.dat")
	// inFlightFile := path.Join(queueDirPath, "in_flight.dat")
	// delayedFile := path.Join(queueDirPath, "delayed.dat")

	var (
		messages []byte
		// inFlightMessages []byte
		// delayedMessages  []byte
		err error
	)

	if messages, err = getMmapedRegion(dataFile); err != nil {
		return nil, err
	}

	// if inFlightMessages, err = getMmapedRegion(inFlightFile); err != nil {
	// 	return nil, err
	// }

	// if delayedMessages, err = getMmapedRegion(delayedFile); err != nil {
	// 	return nil, err
	// }

	return &Queue{
		meta:     (*Meta)(unsafe.Pointer(&messages[0])),
		messages: messages,
		// inFlightMessages: make(map[string][]byte),
		// delayedMessages:  make(map[string][]byte),
	}, nil
}

func (q *Queue) initMeta(delaySeconds uint16, maxMsgSize uint32, messageRetentionPeriod uint32, receiveMessageWaitTime uint16, defaultVisiblityTimeout uint16, tags *map[string]string) {
	q.msgMutex.Lock()
	defer q.msgMutex.Unlock()

	meta := (*Meta)(unsafe.Pointer(&q.messages[0]))

	meta.MagicNum = MAGIC_NUMBER
	meta.CapacityLeft = INITIAL_QUEUE_FILE_SIZE - HEADER_SIZE
	meta.Count = 0
	meta.ReadOffset = HEADER_SIZE
	meta.WriteOffset = HEADER_SIZE
	meta.MaxMessageSize = maxMsgSize
	meta.MessageRetentionPeriod = messageRetentionPeriod
	meta.MessageWaitSeconds = receiveMessageWaitTime
	meta.DelaySeconds = delaySeconds
	meta.VisibilityTimeout = defaultVisiblityTimeout
	meta.HeaderMetaMarker = 0xADDE

	headerSize := int(unsafe.Sizeof(*meta))
	bytes := 0
	for key, val := range *tags {
		copy(q.messages[headerSize+bytes:], key)
		bytes = bytes + len(key) + 1
		q.messages[headerSize+bytes-1] = 0x00

		copy(q.messages[headerSize+bytes:], val)
		bytes = bytes + len(val) + 1
		q.messages[headerSize+bytes-1] = 0x00
	}
}

func NewQueue(queueName string) (*Queue, error) {
	// we delete existing folder and create new one
	queueDirPath := path.Join("data", queueName)
	if err := createQueueFolderExists(queueDirPath); err != nil {
		return nil, err
	}

	if err := createQueueDataFilesExists(queueDirPath); err != nil {
		return nil, err
	}

	return getQueue(queueDirPath)
}

func (q *Queue) Push(msg *Message) error {
	q.msgMutex.Lock()
	defer q.msgMutex.Unlock()

	data, _ := serializeMessage(msg)

	off := q.meta.WriteOffset
	capacityLeft := q.meta.CapacityLeft

	if capacityLeft < uint32(len(data)) {
		return errors.New("queue is full")
	}

	binary.LittleEndian.PutUint32(q.messages[off:], uint32(len(data)))
	n := copy(q.messages[off+4:], data)

	q.meta.Count++
	q.meta.CapacityLeft = capacityLeft - uint32(n) - 4

	nextOffset := off + uint32(n) + 4
	if nextOffset > uint32(DATA_BUFFER_SIZE) {
		q.meta.WriteOffset = (nextOffset % uint32(DATA_BUFFER_SIZE)) + HEADER_SIZE // todo check if this is correct
	} else {
		q.meta.WriteOffset = nextOffset
	}

	return nil
}

func (q *Queue) Pop() (*Message, error) {
	q.msgMutex.Lock()
	defer q.msgMutex.Unlock()

	off := q.meta.ReadOffset
	messageSize := binary.LittleEndian.Uint32(q.messages[off:])

	msg, err := deserializeMessage(q.messages[off+4 : off+4+messageSize])
	if err != nil {
		logs.Logger.Error().Err(err).Msg("error deserializing message")
		return nil, err
	}

	q.meta.Count--
	q.meta.CapacityLeft = q.meta.CapacityLeft + messageSize + 4

	nextOffset := off + uint32(messageSize) + 4

	if nextOffset > uint32(DATA_BUFFER_SIZE) {
		q.meta.ReadOffset = (nextOffset % uint32(DATA_BUFFER_SIZE)) + HEADER_SIZE // todo check if this is correct
	} else {
		q.meta.ReadOffset = nextOffset
	}

	return msg, nil
}

func (q *Queue) closeMmap() {
	dh := (*reflect.SliceHeader)(unsafe.Pointer(&q.messages))
	_, _, err := syscall.Syscall(syscall.SYS_MUNMAP, uintptr(dh.Data), uintptr(dh.Len), 0)
	if err != syscall.Errno(0) {
		logs.Logger.Warn().Err(err).Msg("error closing mmap")
	}
}

//	func (mm mmap) Sync() {
//		rh := *(*reflect.SliceHeader)(unsafe.Pointer(&mm))
//		_, _, err := syscall.Syscall(syscall.SYS_MSYNC, uintptr(rh.Data), uintptr(rh.Len), uintptr(syscall.MS_ASYNC))
//		if err != 0 {
//			panic(err)
//		}
