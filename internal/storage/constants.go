package storage

import (
	"unsafe"
)

const MAGIC_NUMBER = 0x01535156 // VQS (0x56 0x51 0x53) Followed By 0x01 (version 1) in reversed order becase we use little endianness

var (
	DATA_BUFFER_SIZE         uint32 = 128 << 20
	MMAP_BUFFER_SIZE         uint32 = 1 << 20
	META_FILE_SIZE           uint32 = 1 << 20
	META_FILE_META_DATA_SIZE uint32 = uint32(unsafe.Sizeof(Meta{})) // 36 as of now
)
