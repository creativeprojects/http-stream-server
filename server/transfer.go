package server

import (
	"net/http"
	"sync"
)

type Transfer struct {
	hasReceiver          bool
	hasReceiverMutex     *sync.Mutex
	receiverChan         chan http.ResponseWriter
	senderFinishedChan   chan interface{}
	receiverFinishedChan chan interface{}
}

var (
	pathToTransfer map[string]*Transfer
	mapSync        sync.Mutex
)

func init() {
	// Initialize map
	pathToTransfer = map[string]*Transfer{}
	mapSync = sync.Mutex{}
}

func newTransfer() *Transfer {
	return &Transfer{
		hasReceiverMutex:     &sync.Mutex{},
		receiverChan:         make(chan http.ResponseWriter),
		senderFinishedChan:   make(chan interface{}),
		receiverFinishedChan: make(chan interface{}),
	}
}

// getTransfer returns the Transfer object registered at that path.
// If the path does not exist, it creates a new Transfer and registers it at that path.
func getTransfer(path string) *Transfer {
	mapSync.Lock()
	defer mapSync.Unlock()

	if _, ok := pathToTransfer[path]; !ok {
		pathToTransfer[path] = newTransfer()
	}
	return pathToTransfer[path]
}

func deleteTransfer(path string) {
	mapSync.Lock()
	defer mapSync.Unlock()

	delete(pathToTransfer, path)
}

func (t *Transfer) lockReceiver() bool {
	t.hasReceiverMutex.Lock()
	defer t.hasReceiverMutex.Unlock()

	if t.hasReceiver {
		return false
	}
	t.hasReceiver = true
	return true
}
