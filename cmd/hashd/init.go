package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	grpcapi "hash-sv/internal/grpc"
	"hash-sv/internal/hash"
	httphandler "hash-sv/internal/http"
)

func NewHashService(
	storage hash.Storage,
	generator hash.UUIDGenerator,
	logger *zap.Logger,
	cfg config,
) *hash.Service {
	hashUpdater := hash.NewUpdater(storage, generator, cfg.HashTTL, logger)

	return hash.NewService(hashUpdater, storage)
}

func NewHttpRouter(service *hash.Service, logger *zap.Logger) http.Handler {
	router := mux.NewRouter()
	router.Handle("/hash", httphandler.NewHashHandler(service, logger))

	return router
}

func NewGrpcServer(service *hash.Service) *grpc.Server {
	s := grpc.NewServer()
	grpcapi.RegisterHashServer(s, grpcapi.NewServer(service))

	return s
}

func MustNewNetListener(cfg config) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort))
	if err != nil {
		panic(fmt.Errorf("failed to listen: %w", err))
	}

	return lis
}

func NewHttpServer(handler http.Handler, cfg config) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HttpPort),
		Handler: handler,
	}
}

func NewLogger(cfg config) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	var core zapcore.Core
	if cfg.Debug {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, zap.DebugLevel)
	} else {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, zap.InfoLevel)
	}

	return zap.New(core)
}
