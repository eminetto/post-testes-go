package api

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	logger "github.com/PicPay/lib-go-logger"
)

const TIMEOUT = 30 * time.Second

//@todo esse pacote poderia ser uma lib compartilhada
func Start(l *logger.Logger, port string, handler http.Handler) error {
	srv := &http.Server{
		ReadTimeout:  TIMEOUT,
		WriteTimeout: TIMEOUT,
		Addr:         ":" + port,
		Handler:      handler,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer stop()
	errShutdown := make(chan error, 1)
	go shutdown(srv, ctx, errShutdown)

	l.Info(fmt.Sprintf("Current service listening on port %s\n", port))
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	err = <-errShutdown
	if err != nil {
		return err
	}
	return nil
}

func shutdown(server *http.Server, ctxShutdown context.Context, errShutdown chan error) {
	<-ctxShutdown.Done()

	ctxTimeout, stop := context.WithTimeout(context.Background(), TIMEOUT)
	defer stop()

	err := server.Shutdown(ctxTimeout)
	switch err {
	case nil:
		errShutdown <- nil
	case context.DeadlineExceeded:
		errShutdown <- fmt.Errorf("Forcing closing the server")
	default:
		errShutdown <- fmt.Errorf("Forcing closing the server")
	}
}
