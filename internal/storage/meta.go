package storage

import (
	"encoding/binary"
	"os"
	"reflect"
	"syscall"
	"unsafe"
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

// Some black magic happens here, follow very closely to understand
func getMetaAndTagsRegion(path string) (*Meta, *[]byte, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	mem, _, err := syscall.Syscall6(
		syscall.SYS_MMAP,
		0,                                    // let kernel choose the address
		uintptr(META_FILE_SIZE),              // size of the memory region
		syscall.PROT_READ|syscall.PROT_WRITE, // read and write access
		syscall.MAP_SHARED,                   //
		file.Fd(),                            // this is backed by meta file
		0,                                    // offset
	)
	if err != syscall.Errno(0) {
		return nil, nil, err
	}

	var metaDataBuffer []byte
	dh := (*reflect.SliceHeader)(unsafe.Pointer(&metaDataBuffer))
	dh.Data = mem
	dh.Len = int(META_FILE_META_DATA_SIZE)
	dh.Cap = int(META_FILE_META_DATA_SIZE)

	meta := (*Meta)(unsafe.Pointer(&metaDataBuffer[0]))

	var tags []byte
	dha := (*reflect.SliceHeader)(unsafe.Pointer(&tags))
	dha.Data = mem + uintptr(META_FILE_META_DATA_SIZE)
	dha.Len = int(META_FILE_SIZE - META_FILE_META_DATA_SIZE)
	dha.Cap = int(META_FILE_SIZE - META_FILE_META_DATA_SIZE)

	return meta, &tags, nil
}

func (q *Queue) initMeta(delaySeconds uint16, maxMsgSize uint32, messageRetentionPeriod uint32, receiveMessageWaitTime uint16, defaultVisiblityTimeout uint16) {
	q.metaMutex.Lock()
	defer q.metaMutex.Unlock()

	q.meta.MagicNum = MAGIC_NUMBER
	q.meta.CapacityLeft = DATA_BUFFER_SIZE
	q.meta.Count = 0
	q.meta.ReadOffset = 0
	q.meta.WriteOffset = 0
	q.meta.MaxMessageSize = maxMsgSize
	q.meta.MessageRetentionPeriod = messageRetentionPeriod
	q.meta.MessageWaitSeconds = receiveMessageWaitTime
	q.meta.DelaySeconds = delaySeconds
	q.meta.VisibilityTimeout = defaultVisiblityTimeout
	q.meta.HeaderMetaMarker = 0xADDE
}

func (q *Queue) initTags(tags *[][2]string) {
	q.metaMutex.Lock()
	defer q.metaMutex.Unlock()

	var bytes = 0
	for _, tag := range *tags {
		keyLen := len(tag[0])
		valueLen := len(tag[1])

		binary.LittleEndian.PutUint16(q.tags[bytes:], uint16(keyLen))
		bytes = bytes + 2
		binary.LittleEndian.PutUint16(q.tags[bytes:], uint16(valueLen))
		bytes = bytes + 2

		copy(q.tags[bytes:], tag[0])
		bytes = bytes + keyLen
		copy(q.tags[bytes:], tag[1])
		bytes = bytes + valueLen
	}
}
