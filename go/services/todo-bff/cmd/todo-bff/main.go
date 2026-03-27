package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"

	app_validator "github.com/phamquanandpad/training-project/go/services/todo-bff/internal/handler/graph/validator"

	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/config"
	"github.com/phamquanandpad/training-project/go/services/todo-bff/internal/registry"
)

func main() {
	conf := config.Load()

	// init validator
	v := app_validator.New()

	// launch Http server
	listener, closer, errCh := runHTTPServer(conf, v)
	defer closer()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-sigCh:
			_ = listener.Close()
		case _ = <-errCh:
			_ = listener.Close()
		}
		cancel()
	}()
	<-ctx.Done()
}

func runHTTPServer(c *config.Config, v *validator.Validate) (listener net.Listener, closer func(), ch chan error) {
	ctx := context.Background()
	router, cleanup, err := registry.InitHttpServer(ctx, c, v)
	if err != nil {
		panic(err)
	}

	var lc net.ListenConfig
	addr := fmt.Sprintf(":%s", c.Port)
	listener, err2 := lc.Listen(ctx, "tcp", addr)
	if err2 != nil {
		panic(err2)
	}

	ch = make(chan error)
	go func() {
		srv := &http.Server{
			Handler:           router,
			ReadTimeout:       15 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       30 * time.Second,
		}
		ch <- srv.Serve(listener)
	}()

	return listener, cleanup, ch
}
