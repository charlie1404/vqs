package api

import "github.com/charlie1404/vqueue/pkg/storage"

type AppContext struct {
	queues    *storage.Queues
	validator ApiValidator
}
