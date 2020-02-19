package storage

type MessageID int

type Message struct {
	ID      MessageID
	Title   string
	Message string
}
