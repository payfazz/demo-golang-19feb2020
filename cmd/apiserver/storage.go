package main

import (
	"github.com/payfazz/go-errors/errhandler"

	"api/internal/env"
	"api/internal/storage"
	"api/internal/storage/postgres"
	"api/internal/storage/ram"
)

func getStorage() storage.Storage {
	if env.PostgresConnection == "" {
		return ram.New()
	}

	p, err := postgres.New(env.PostgresConnection)
	errhandler.Check(err)
	return p
}
