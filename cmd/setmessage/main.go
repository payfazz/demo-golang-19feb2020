package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/payfazz/go-errors"
	"github.com/payfazz/go-errors-ext/errhandlerext"
	"github.com/payfazz/go-errors/errhandler"
	"github.com/payfazz/mainutil"
	"github.com/payfazz/stdlog"

	"api/internal/env"
	"api/internal/storage"
	"api/internal/storage/postgres"
)

var m *mainutil.Env

func main() {
	defer errhandler.With(errhandlerext.LogAndExit)

	if len(os.Args) != 4 {
		stdlog.PrintErr(fmt.Sprintf("Usage: %s <id> <title> <message>", os.Args[0]))
		return
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	m.CancelOnInteruptSignal(cancelCtx)

	if env.PostgresConnection == "" {
		stdlog.PrintErr("PostgresConnection cannot be empty")
		return
	}

	sto, err := postgres.New(env.PostgresConnection)
	errhandler.Check(err)
	defer sto.Close()

	id, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		errhandler.Check(errors.Wrap(err))
	}

	err = sto.StoreMessage(ctx, &storage.Message{
		ID:      storage.MessageID(id),
		Title:   os.Args[2],
		Message: os.Args[3],
	})
	errhandler.Check(err)
}
