package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	grpcapi "hash-sv/internal/grpc"
	"hash-sv/internal/hash"
	http2 "hash-sv/internal/http"
)

var testCases = []struct {
	hash string
}{
	{"f47ac10b-58cc-f372-8567-0e02b2c3d471"},
	{"f47ac10b-58cc-4372-0567-0e02b2c3d479"},
}

func TestGrpcHashServer_Get(t *testing.T) {
	var (
		cfg           = config{HashTTL: 500 * time.Millisecond}
		logger        = zap.NewNop()
		ctrl          = gomock.NewController(t)
		hashGenerator = hash.NewMockUUIDGenerator(ctrl)
	)

	defer ctrl.Finish()

	var (
		updatedCh = make(chan struct{})
		storage   = hash.NewStorageFuncWrapper(hash.NewMemoryStorage(), func() { updatedCh <- struct{}{} })
		service   = NewHashService(storage, hashGenerator, logger, cfg)
	)

	for _, h := range mustParseHashes(t, testCases) {
		h := h
		hashGenerator.EXPECT().Generate().DoAndReturn(func() hash.UUID { return h })
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go func() {
		err := service.Run(ctx)
		assert.Nil(t, err, "hash updater")
	}()

	listener := bufconn.Listen(5 * 1024)
	grpcServer := NewGrpcServer(service)

	go func() {
		err := grpcServer.Serve(listener)
		assert.Nil(t, err, "grpc server")
	}()

	defer func() {
		err := listener.Close()
		assert.Nil(t, err, "closing listener")

		grpcServer.Stop()
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.Nil(t, err, "client dial")

	client := grpcapi.NewHashClient(conn)
	for _, tc := range testCases {
		<-updatedCh
		response, err := client.Get(ctx, &grpcapi.Empty{})
		assert.Nil(t, err, "client get request")
		assert.Equal(t, tc.hash, response.Hash)
	}
}

func TestHttpHashServer_Get(t *testing.T) {
	var (
		cfg           = config{HashTTL: 500 * time.Millisecond}
		logger        = zap.NewNop()
		ctrl          = gomock.NewController(t)
		hashGenerator = hash.NewMockUUIDGenerator(ctrl)
	)

	defer ctrl.Finish()

	var (
		updatedCh = make(chan struct{})
		storage   = hash.NewStorageFuncWrapper(hash.NewMemoryStorage(), func() { updatedCh <- struct{}{} })
		service   = NewHashService(storage, hashGenerator, logger, cfg)
	)

	for _, h := range mustParseHashes(t, testCases) {
		h := h
		hashGenerator.EXPECT().Generate().DoAndReturn(func() hash.UUID { return h })
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go func() {
		err := service.Run(ctx)
		assert.Nil(t, err, "hash updater")
	}()

	server := httptest.NewServer(NewHttpRouter(service, logger))

	for _, tc := range testCases {
		<-updatedCh
		response, err := http.Get(server.URL + "/hash")
		assert.Nil(t, err, "http client get request")
		assert.Equal(t, http.StatusOK, response.StatusCode)
		var resp http2.HashRowResponse
		err = json.NewDecoder(response.Body).Decode(&resp)
		assert.Nil(t, err, "decoding hash row")
		assert.Equal(t, tc.hash, resp.Hash.String())
	}
}

func mustParseHashes(t *testing.T, testCases []struct{ hash string }) []hash.UUID {
	hashes := make([]hash.UUID, 0, len(testCases))

	for _, testCase := range testCases {
		h, err := uuid.Parse(testCase.hash)
		assert.Nil(t, err, "could not parse test case hash")

		hashes = append(hashes, h)
	}

	return hashes
}
