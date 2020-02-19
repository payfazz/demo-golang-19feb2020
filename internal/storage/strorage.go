package storage

import "context"

type Storage interface {
	StoreMessage(ctx context.Context, msg *Message) error
	GetMessage(ctx context.Context, ID MessageID) (*Message, error)
	DeleteMessage(ctx context.Context, ID MessageID) error
}
