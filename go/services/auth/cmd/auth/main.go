package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/panjf2000/ants/v2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	app_validator "github.com/phamquanandpad/training-project/go/services/auth/internal/handler/validator"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/config"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/registry"
)

const (
	waitingTimeSecForGracefulShutdown = 30 * time.Second
)

func initializeServer(
	conf *config.Config,
	validate *validator.Validate,
) (*grpc.Server, func()) {
	server, cleanup, err := registry.InitializeServer(conf, validate)
	if err != nil {
		log.Panic("failed to wire")
	}

	if conf.GrpcReflectionEnable {
		reflection.Register(server)
	}

	return server, cleanup
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	validate, err := app_validator.InitValidator()
	if err != nil {
		log.Panic("failed to init validator", err)
	}

	server, clean := initializeServer(cfg, validate)
	defer clean()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		var lc net.ListenConfig
		l, err := lc.Listen(
			ctx,
			"tcp",
			fmt.Sprintf(":%d", cfg.ServerPort),
		)
		if err != nil {
			log.Panic("failed to listen")
		}
		log.Printf("server listening on port %d", cfg.ServerPort)
		//nolint: wrapcheck
		return server.Serve(l)
	})

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	log.Print("received shutdown signal")

	cancel()

	// Termination Processing
	_, shutdownCancel := context.WithTimeout(
		context.Background(),
		waitingTimeSecForGracefulShutdown,
	)
	defer shutdownCancel()

	server.GracefulStop()
	if err := eg.Wait(); err != nil {
		log.Printf("server returning an error: %v", err)
	}

	ants.Release()
}
