package postgres

import (
	"context"
	"io"

	"github.com/payfazz/go-errors"

	"api/internal/storage"
)

type Postgres struct{}

// static check for "implement interface"
var (
	_ storage.Storage = (*Postgres)(nil)
	_ io.Closer       = (*Postgres)(nil)
)

func New(connection string) (*Postgres, error) {
	return nil, errors.New("Not Yet Implemented")
}

func (p *Postgres) Close() error {
	return nil
}

func (p *Postgres) StoreMessage(ctx context.Context, msg *storage.Message) error {
	return errors.New("Not Yet Implemented")
}

func (p *Postgres) GetMessage(ctx context.Context, ID storage.MessageID) (*storage.Message, error) {
	return nil, errors.New("Not Yet Implemented")
}

func (p *Postgres) DeleteMessage(ctx context.Context, ID storage.MessageID) error {
	return errors.New("Not Yet Implemented")
}
