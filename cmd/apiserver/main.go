package main

import (
	"context"
	"io"

	"github.com/payfazz/go-errors-ext/errhandlerext"
	"github.com/payfazz/go-errors/errhandler"
	"github.com/payfazz/go-middleware"
	"github.com/payfazz/mainutil"

	"api/internal/env"
)

var m *mainutil.Env

func main() {
	defer errhandler.With(errhandlerext.LogAndExit)
	var err error

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	m.CancelOnInteruptSignal(cancelCtx)

	sto := getStorage()
	if closer, ok := sto.(io.Closer); ok {
		defer closer.Close()
	}

	server := m.DefaultHTTPServer(env.Addr, middleware.Compile(
		m.CommonHTTPMiddlware(true),
		createHandler(sto, m.ErrLogger()),
	))

	err = m.RunHTTPServerOn(ctx, server, nil, 0)
	if err != nil && ctx.Err() == nil {
		errhandler.Check(err)
	}
}
