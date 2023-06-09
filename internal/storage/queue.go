package storage

import (
	"encoding/binary"
	"errors"
	"path"
	"reflect"
	"sync"
	"syscall"
	"unsafe"

	"github.com/charlie1404/vqs/internal/o11y/logs"
)

type Queue struct {
	meta      *Meta
	metaMutex sync.RWMutex // RLock can only be used for peek cases as of now
	messages  []byte
	msgMutex  sync.RWMutex // RLock can only be used for peek cases as of now
	tags      []byte
	// inFlightMessages []byte // for now we will keep these in memory, later put on disk
	// inFlightMutex    sync.Mutex
	// delayedMessages  []byte // for now we will keep these in memory, later put on disk
	// delayedMutex     sync.Mutex
}

func getQueue(queueDirPath string) (*Queue, error) {
	metaFile := path.Join(queueDirPath, "meta.dat")
	dataFile := path.Join(queueDirPath, "data.dat")
	// inFlightFile := path.Join(queueDirPath, "in_flight.dat")
	// delayedFile := path.Join(queueDirPath, "delayed.dat")

	var (
		meta     *Meta
		messages *[]byte
		tags     *[]byte
		err      error
		// inFlightMessages []byte
		// delayedMessages  []byte
	)

	if messages, err = getMessagesMmapedRegion(dataFile); err != nil {
		return nil, err
	}

	if meta, tags, err = getMetaAndTagsRegion(metaFile); err != nil {
		return nil, err
	}

	// if inFlightMessages, err = getMmapedRegion(inFlightFile); err != nil {
	// 	return nil, err
	// }

	// if delayedMessages, err = getMmapedRegion(delayedFile); err != nil {
	// 	return nil, err
	// }

	return &Queue{
		meta:     meta,
		messages: *messages,
		tags:     *tags,
		// inFlightMessages: make(map[string][]byte),
		// delayedMessages:  make(map[string][]byte),
	}, nil
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

	data, n, _ := serializeMessage(msg)

	capacityLeft := q.meta.CapacityLeft
	nextOffset := q.meta.WriteOffset + uint32(n) + 4

	if capacityLeft < uint32(len(data)) {
		return errors.New("queue is full")
	}

	binary.LittleEndian.PutUint32(q.messages[q.meta.WriteOffset:], uint32(len(data)))
	copy(q.messages[q.meta.WriteOffset+4:], data)

	q.meta.Count++
	q.meta.CapacityLeft = capacityLeft - uint32(n) - 4
	q.meta.WriteOffset = nextOffset % DATA_BUFFER_SIZE

	// if nextOffset > DATA_BUFFER_SIZE {
	// 	q.meta.WriteOffset = nextOffset % DATA_BUFFER_SIZE
	// } else {
	// 	q.meta.WriteOffset = nextOffset
	// }

	return nil
}

func (q *Queue) Pop() (*Message, error) {
	q.msgMutex.Lock()
	defer q.msgMutex.Unlock()

	if q.meta.Count == 0 {
		return nil, nil
	}

	messageSize := binary.LittleEndian.Uint32(q.messages[q.meta.ReadOffset:])
	nextOffset := q.meta.ReadOffset + messageSize + 4

	msg, err := deserializeMessage(q.messages[q.meta.ReadOffset+4 : nextOffset])
	if err != nil {
		logs.Logger.Error().Err(err).Msg("error deserializing message")
		return nil, err
	}

	q.meta.Count--
	q.meta.CapacityLeft = q.meta.CapacityLeft + messageSize + 4
	q.meta.ReadOffset = nextOffset % DATA_BUFFER_SIZE

	// if nextOffset > DATA_BUFFER_SIZE {
	// 	q.meta.ReadOffset = nextOffset % DATA_BUFFER_SIZE
	// } else {
	// 	q.meta.ReadOffset = nextOffset
	// }

	return msg, nil
}

func (q *Queue) closeMmap() {
	dh := (*reflect.SliceHeader)(unsafe.Pointer(&q.messages))
	_, _, err := syscall.Syscall(syscall.SYS_MUNMAP, uintptr(dh.Data), uintptr(dh.Len), 0)
	if err != syscall.Errno(0) {
		logs.Logger.Warn().Err(err).Msg("error closing data mmap")
	}

	meta := (*reflect.SliceHeader)(unsafe.Pointer(unsafe.Pointer(&q.meta)))
	_, _, err = syscall.Syscall(syscall.SYS_MUNMAP, uintptr(meta.Data), uintptr(META_FILE_SIZE), 0)
	if err != syscall.Errno(0) {
		logs.Logger.Warn().Err(err).Msg("error closing meta mmap")
	}
}

//	func (mm mmap) Sync() {
//		rh := *(*reflect.SliceHeader)(unsafe.Pointer(&mm))
//		_, _, err := syscall.Syscall(syscall.SYS_MSYNC, uintptr(rh.Data), uintptr(rh.Len), uintptr(syscall.MS_ASYNC))
//		if err != 0 {
//			panic(err)
//		}
