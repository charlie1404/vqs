package storage

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"syscall"
	"unsafe"

	"github.com/charlie1404/vqueue/pkg/constants"
	"github.com/charlie1404/vqueue/pkg/utils"
)

type Queue struct {
	dataMessages []byte
	// inFlightMessages []byte // TODO: need this later
	readMutex  sync.Mutex
	writeMutex sync.Mutex
}

func getFileHandle(fileName string) (*os.File, error) {
	var (
		file          *os.File
		err           error
		queueFileName string
	)

	queueFileName = fmt.Sprintf("data/%s.dat", utils.Hash32(fileName))

	if _, fStatErr := os.Stat(queueFileName); fStatErr != nil {
		if os.IsNotExist(fStatErr) {
			if file, err = os.Create(queueFileName); err != nil {
				return nil, err
			}

			if err := syscall.Truncate(queueFileName, constants.INITIAL_QUEUE_FILE_SIZE); err != nil {
				os.Remove(queueFileName)
				return nil, err
			}
			return file, nil
		}
		return nil, fStatErr
	} else {
		return nil, errors.New("Queue Already exist")
	}
}

func newQueue(file *os.File) (*Queue, error) {
	var memoryMapSize uintptr = constants.INITIAL_QUEUE_FILE_SIZE + constants.MMAP_BUFFER_SIZE

	// void *mmap(void *addr, size_t length, int prot, int flags, int fd, off_t offset);
	fd := -1
	mem, _, err := syscall.Syscall6(syscall.SYS_MMAP, 0, memoryMapSize, syscall.PROT_NONE, syscall.MAP_PRIVATE|syscall.MAP_ANON, uintptr(fd), 0)
	if err != syscall.Errno(0) {
		return nil, err
	}

	_, _, err = syscall.Syscall6(syscall.SYS_MMAP, mem, uintptr(constants.INITIAL_QUEUE_FILE_SIZE), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_FIXED, file.Fd(), 0)
	if err != syscall.Errno(0) {
		return nil, err
	}

	_, _, err = syscall.Syscall6(syscall.SYS_MMAP, mem+constants.INITIAL_QUEUE_FILE_SIZE, uintptr(constants.MMAP_BUFFER_SIZE), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_FIXED, file.Fd(), constants.HEADER_SIZE)
	if err != syscall.Errno(0) {
		return nil, err
	}

	var byteArray []byte
	dh := (*reflect.SliceHeader)(unsafe.Pointer(&byteArray))
	dh.Data = mem
	dh.Len = constants.INITIAL_QUEUE_FILE_SIZE + constants.MMAP_BUFFER_SIZE
	dh.Cap = constants.INITIAL_QUEUE_FILE_SIZE + constants.MMAP_BUFFER_SIZE

	return &Queue{
		dataMessages: byteArray,
	}, nil
}

func (q *Queue) initMeta(attributes *QueueAttributes, tags *QueueTags) {
	q.writeMutex.Lock()
	defer q.writeMutex.Unlock()

	*(*uint32)(unsafe.Pointer(&q.dataMessages[constants.FILE_HEADER_MAGIC_NUMBER_OFFSET])) = constants.MAGIC_NUMBER
	*(*uint32)(unsafe.Pointer(&q.dataMessages[constants.FILE_HEADER_TOTAL_CAPACITY_LEFT_OFFSET])) = constants.INITIAL_QUEUE_FILE_SIZE - constants.HEADER_SIZE
	*(*uint32)(unsafe.Pointer(&q.dataMessages[constants.FILE_HEADER_COUNT_OFFSET])) = 0
	*(*uint32)(unsafe.Pointer(&q.dataMessages[constants.FILE_HEADER_READ_OFFSET])) = constants.HEADER_SIZE
	*(*uint32)(unsafe.Pointer(&q.dataMessages[constants.FILE_HEADER_WRITE_OFFSET])) = constants.HEADER_SIZE
	*(*uint16)(unsafe.Pointer(&q.dataMessages[constants.FILE_HEADER_DELAY_SECONDS_OFFSET])) = uint16(attributes.DelaySeconds)
	*(*uint32)(unsafe.Pointer(&q.dataMessages[constants.FILE_HEADER_MAX_MESSAGE_SIZE_OFFSET])) = uint32(attributes.MaximumMessageSize)
	*(*uint32)(unsafe.Pointer(&q.dataMessages[constants.FILE_HEADERM_MESSAGE_RETENTION_OFFSET])) = uint32(attributes.MessageRetentionPeriod)
	*(*uint16)(unsafe.Pointer(&q.dataMessages[constants.FILE_HEADER_MESSAGE_WAIT_OFFSET])) = uint16(attributes.ReceiveMessageWaitTimeSeconds)
	*(*uint16)(unsafe.Pointer(&q.dataMessages[constants.FILE_HEADER_VISIBLITY_TIMEOUT_OFFSET])) = uint16(attributes.VisibilityTimeout)

	bytes := 0
	for key, val := range *tags {
		copy(q.dataMessages[constants.FILE_HEADER_TAGS_OFFSET+bytes:], key)
		bytes = bytes + len(key) + 1
		q.dataMessages[constants.FILE_HEADER_TAGS_OFFSET+bytes-1] = 0

		copy(q.dataMessages[constants.FILE_HEADER_TAGS_OFFSET+bytes:], val)
		bytes = bytes + len(val) + 1
		q.dataMessages[constants.FILE_HEADER_TAGS_OFFSET+bytes-1] = 0
	}
}
