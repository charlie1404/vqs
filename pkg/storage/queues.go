package storage

import (
	"sync"
)

type Queues struct {
	queues map[string]*Queue
	sync.Mutex
}

func (qs *Queues) CreateQueue(createQueueInput *CreateQueueInput) error {
	qs.Lock()
	defer qs.Unlock()

	file, err := getFileHandle(createQueueInput.QueueName)

	if err != nil {
		return err
	}

	queue, err := newQueue(file)
	defer file.Close()

	if err != nil {
		return err
	}

	queue.initMeta(&createQueueInput.Attributes, &createQueueInput.Tags)

	qs.queues[createQueueInput.QueueName] = queue

	return nil
}

func NewQueues() *Queues {
	return &Queues{
		queues: make(map[string]*Queue),
	}
}
