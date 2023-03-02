package storage

import "unsafe"

func GetMaxHeaderTagsSize() int {
	var m Meta
	headerSize := int(unsafe.Sizeof(m))

	return headerSize
}
