package utils

import "unsafe"

func GetLittleEndianUint32(b []byte, offset int) uint32 {
	return *(*uint32)(unsafe.Pointer(&b[offset]))
}

func PutLittleEndianUint32(b []byte, offset int, val uint32) {
	*(*uint32)(unsafe.Pointer(&b[offset])) = val
}
