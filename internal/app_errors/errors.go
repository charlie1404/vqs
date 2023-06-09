package app_errors

import (
	"errors"
)

var (
	QueueNotExists         = errors.New("QueueNotExists")
	CreateQueueQueueExists = errors.New("CreateQueueQueueExists")
)
