package utils

import (
	"github.com/charlie1404/vqueue/pkg/constants"
)

func ConstrainFileSizeLimit(bytes int) int {
	size := bytes
	if bytes > constants.MAX_BUFFER_FILE_SIZE {
		size = constants.MAX_BUFFER_FILE_SIZE
	} else if bytes < constants.MIN_BUFFER_FILE_SIZE {
		size = constants.MIN_BUFFER_FILE_SIZE
		size = 4056
	}

	return constants.DATA_OFFSET + constants.DATA_SIZE + size - (size % constants.DATA_SIZE)
}
