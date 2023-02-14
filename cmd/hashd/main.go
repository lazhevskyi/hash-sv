package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"hash-sv/internal/hash"
)

func main() {
	cfg := mustParseConfig()

	sigCh := make(chan os.Signal)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	logger := NewLogger(cfg)
	service := NewHashService(hash.NewMemoryStorage(), hash.NewUUID4Generator(), logger, cfg)

	errG, errCtx := errgroup.WithContext(ctx)

	errG.Go(func() error {
		return service.Run(errCtx)
	})

	httpServer := NewHttpServer(NewHttpRouter(service, logger), cfg)

	errG.Go(func() error {
		return httpServer.ListenAndServe()
	})

	errCh := make(chan error)

	go func() {
		if err := errG.Wait(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-sigCh:
		cancel()
		logger.Info("graceful shutdown")
	case err := <-errCh:
		logger.Error("finishing with error", zap.Error(err))
		os.Exit(1)
	}
}
