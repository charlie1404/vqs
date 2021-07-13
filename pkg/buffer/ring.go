package buffer

import (
	"errors"
	"os"
	"sync"
	"syscall"

	"github.com/charlie1404/vqueue/pkg/constants"
	"github.com/charlie1404/vqueue/pkg/utils"
)

type Buffer struct {
	sync.Mutex
	data []byte
}

type fileHeader struct {
	magic_num      uint32
	total_capacity uint32
	count          uint32
	read_off       uint32
	write_off      uint32
}

func (buf *Buffer) read_header() fileHeader {
	magic_number := utils.GetLittleEndianUint32(buf.data, constants.HEADER_MAGIC_NUMBER_OFFSET)
	if magic_number != constants.MAGIC_NUMBER {
		// it should not reach here ever, but if came, then it seems system is set on fire, lets extinguish it or may be we have a bug in system
		panic("Something went terribly wrong. Abort")
	}
	header := fileHeader{
		magic_num:      magic_number,
		total_capacity: utils.GetLittleEndianUint32(buf.data, constants.HEADER_TOTAL_CAPACITY_OFFSET),
		count:          utils.GetLittleEndianUint32(buf.data, constants.HEADER_COUNT_OFFSET),
		read_off:       utils.GetLittleEndianUint32(buf.data, constants.HEADER_READ_OFFSET),
		write_off:      utils.GetLittleEndianUint32(buf.data, constants.HEADER_WRITE_OFFSET),
	}
	return header
}

func (buf *Buffer) write_header(header_data fileHeader) {
	utils.PutLittleEndianUint32(buf.data, constants.HEADER_MAGIC_NUMBER_OFFSET, header_data.magic_num)
	utils.PutLittleEndianUint32(buf.data, constants.HEADER_TOTAL_CAPACITY_OFFSET, header_data.total_capacity)
	utils.PutLittleEndianUint32(buf.data, constants.HEADER_COUNT_OFFSET, header_data.count)
	utils.PutLittleEndianUint32(buf.data, constants.HEADER_READ_OFFSET, header_data.read_off)
	utils.PutLittleEndianUint32(buf.data, constants.HEADER_WRITE_OFFSET, header_data.write_off)
}

func (b *Buffer) Insert(data []byte) error {
	if payload_len := len(data); payload_len != constants.DATA_SIZE {
		return errors.New("Invalid Message")
	}

	b.Lock()
	defer b.Unlock()
	header := b.read_header()

	if header.count+1 > header.total_capacity>>5 {
		return errors.New("Overflow")
	}

	copy(b.data[header.write_off:], data)
	header.count += 1

	// TODO: create abstraction for it
	next_off := header.write_off + constants.DATA_SIZE
	if next_off >= header.total_capacity {
		next_off = constants.DATA_OFFSET
	}
	header.write_off = next_off

	b.write_header(header)
	return nil
}

func (b *Buffer) Pop() []byte {
	b.Lock()
	defer b.Unlock()

	header := b.read_header()

	if header.count == 0 {
		return nil
	}

	ret := make([]byte, constants.DATA_SIZE)
	copy(ret, b.data[header.read_off:])
	header.count -= 1

	next_off := header.read_off + constants.DATA_SIZE
	if next_off >= header.total_capacity {
		next_off = constants.DATA_OFFSET
	}
	header.read_off = next_off

	b.write_header(header)
	return ret
}

func New(queue_name string, capacity int) (*Buffer, error) {
	var created bool
	var file *os.File

	file_size := utils.ConstrainFileSizeLimit(capacity)

	if _, err := os.Stat(queue_name); os.IsNotExist(err) {
		if file, err = os.Create(queue_name); err != nil {
			return nil, err
		}
		created = true
	} else {
		if file, err = os.OpenFile(queue_name, os.O_RDWR, 0644); err != nil {
			return nil, err
		}
	}

	if err := syscall.Truncate(queue_name, int64(file_size)); err != nil {
		return nil, err
	}

	data, err := syscall.Mmap(int(file.Fd()), 0, file_size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	file.Close()
	buf := &Buffer{data: data}

	if created {
		// no lock since things are bootstraping and no threads in system
		buf.write_header(fileHeader{
			magic_num:      constants.MAGIC_NUMBER,
			total_capacity: uint32(file_size - constants.DATA_OFFSET),
			count:          0,
			read_off:       constants.DATA_OFFSET,
			write_off:      constants.DATA_OFFSET,
		})
	} else if magic_number := utils.GetLittleEndianUint32(buf.data, constants.HEADER_MAGIC_NUMBER_OFFSET); magic_number != constants.MAGIC_NUMBER {
		return nil, errors.New("Invalid File.")
	}

	return buf, nil
}
