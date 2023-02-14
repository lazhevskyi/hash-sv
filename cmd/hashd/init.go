package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"hash-sv/internal/hash"
	route "hash-sv/internal/http"
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
	router.Handle("/hash", route.NewHashHandler(service, logger))

	return router
}

func NewHttpServer(handler http.Handler, cfg config) *http.Server {
	return &http.Server{
		Addr:    cfg.Addr + ":" + cfg.Port,
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
