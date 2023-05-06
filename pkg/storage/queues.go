package storage

import (
	"os"
	"path"
	"sync"

	"github.com/charlie1404/vqs/pkg/app_errors"
	"github.com/charlie1404/vqs/pkg/o11y/logs"
)

type Queues struct {
	queues map[string]*Queue
	sync.Mutex
}

func (qs *Queues) CreateQueue(name string, delaySeconds uint16, maxMsgSize uint32, messageRetentionPeriod uint32, receiveMessageWaitTime uint16, defaultVisiblityTimeout uint16, tags *[][2]string) (*Queue, error) {
	// This creates file and mmap and init meta, and returns queue
	qs.Lock()
	defer qs.Unlock()

	// for now every create will be blocked, but in future we can use mutex per queue name
	if queue, ok := qs.queues[name]; ok {
		return queue, app_errors.CreateQueueQueueExists
	}

	queue, err := NewQueue(name)
	if err != nil {
		return nil, err
	}

	queue.initMeta(delaySeconds, maxMsgSize, messageRetentionPeriod, receiveMessageWaitTime, defaultVisiblityTimeout)
	queue.initTags(tags)

	qs.queues[name] = queue

	return queue, nil
}

func (qs *Queues) GetQueue(queueName string) (*Queue, error) {
	if queue, ok := qs.queues[queueName]; ok {
		return queue, nil
	}

	queueDirPath := path.Join("data", queueName)
	_, err := os.Stat(queueDirPath)

	if os.IsNotExist(err) {
		return nil, app_errors.QueueNotExists
	}

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	queue, err := getQueue(queueDirPath)
	if err != nil {
		return nil, err
	}
	return queue, nil
}

func LoadQueues() *Queues {
	qs := &Queues{
		queues: make(map[string]*Queue),
	}

	return qs
}

func (qs *Queues) CloseQueues() {
	qs.Lock()
	defer qs.Unlock()

	for name, queue := range qs.queues {
		logs.Logger.Warn().Str("name", name).Msg("Closing queue")
		queue.closeMmap()
	}
}
