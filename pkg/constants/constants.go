package constants

// 0x7F Followed By VQS (0x56 0x51 0x53)
// are in reversed order becase we use little endianness
const MAGIC_NUMBER = 0x5351567F

const (
	MAX_BUFFER_FILE_SIZE = 50 << 20 // 50 MB
	MIN_BUFFER_FILE_SIZE = 10 << 20 // ~10 MB
)

const HEADER_SIZE = 20
const DATA_SIZE = 32

// TODO add file version info
const (
	HEADER_MAGIC_NUMBER_OFFSET   = iota * 4
	HEADER_TOTAL_CAPACITY_OFFSET = iota * 4
	HEADER_COUNT_OFFSET          = iota * 4
	HEADER_READ_OFFSET           = iota * 4
	HEADER_WRITE_OFFSET          = iota * 4
	DATA_OFFSET                  = iota * 4
)
