package storage

// const (
// 	MAGIC_NUMBER            = 0x01535156 // VQS (0x56 0x51 0x53) Followed By 0x01 (version 1) in reversed order becase we use little endianness
// 	HEADER_SIZE             = 16 << 10   // 16KB
// 	INITIAL_QUEUE_FILE_SIZE = 128 << 20  // 128MB
// 	MMAP_BUFFER_SIZE        = 1 << 20    // 1MB
// )

const (
	MAGIC_NUMBER            = 0x01535156 // VQS (0x56 0x51 0x53) Followed By 0x01 (version 1) in reversed order becase we use little endianness
	INITIAL_QUEUE_FILE_SIZE = 64 << 10   // 64kB
	HEADER_SIZE             = 16 << 10   // 16KB
	MMAP_BUFFER_SIZE        = 16 << 10   // 16KB
)

const DATA_BUFFER_SIZE = INITIAL_QUEUE_FILE_SIZE - HEADER_SIZE
