package logger

import (
	"fmt"
	"sync"
)

type InMemoryHandler struct {
	logs           []Log
	lock           sync.Mutex
	bufferOverflow uint
}

func NewInMemoryHandler(bufferOverflow uint) *InMemoryHandler {
	return &InMemoryHandler{bufferOverflow: bufferOverflow}
}

func (h *InMemoryHandler) Handle(log Log) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.bufferOverflow != 0 && uint(len(h.logs)) >= h.bufferOverflow {
		return fmt.Errorf("InMemoryHandler - number of logs exceeded buffer limit (%d records)", h.bufferOverflow)
	}

	h.logs = append(h.logs, log)

	return nil
}

func (h *InMemoryHandler) IsEmpty() bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	return len(h.logs) == 0
}

func (h *InMemoryHandler) Pop() Log {
	h.lock.Lock()
	defer h.lock.Unlock()

	if len(h.logs) < 1 {
		panic("no logs in handler")
	}

	var log Log
	log, h.logs = h.logs[len(h.logs)-1], h.logs[:len(h.logs)-1]

	return log
}

func (h *InMemoryHandler) Clear() {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.logs = nil
}
