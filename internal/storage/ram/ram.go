package ram

import (
	"context"
	"sync"

	"api/internal/storage"
)

type Ram struct {
	lock sync.RWMutex
	data map[storage.MessageID]*storage.Message
}

// static check for "implement interface"
var (
	_ storage.Storage = (*Ram)(nil)
)

func New() *Ram {
	return &Ram{
		data: make(map[storage.MessageID]*storage.Message),
	}
}

func (r *Ram) StoreMessage(ctx context.Context, msg *storage.Message) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.data[msg.ID] = msg

	return nil
}

func (r *Ram) GetMessage(ctx context.Context, ID storage.MessageID) (*storage.Message, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.data[ID], nil
}

func (r *Ram) DeleteMessage(ctx context.Context, ID storage.MessageID) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.data, ID)

	return nil
}
