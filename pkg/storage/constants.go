package storage

import (
	"os"
	"unsafe"
)

const MAGIC_NUMBER = 0x01535156 // VQS (0x56 0x51 0x53) Followed By 0x01 (version 1) in reversed order becase we use little endianness

var DATA_BUFFER_SIZE uint32 = 32 << 10              // 32KB
var MMAP_BUFFER_SIZE = uint32(os.Getpagesize()) * 1 // should be a multiple of os page size, will set to 1MB later

var META_FILE_SIZE = uint32(os.Getpagesize())                // can do 1MB later
var META_FILE_META_DATA_SIZE = uint32(unsafe.Sizeof(Meta{})) // 36 as of now
