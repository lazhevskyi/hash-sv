package hash

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var testCases = []struct{ hash string }{
	{"f47ac10b-58cc-4372-0567-0e02b2c3d479"},
	{"f47ac10b-58cc-f372-8567-0e02b2c3d471"},
}

func TestUpdater_Run(t *testing.T) {
	var (
		ctrl          = gomock.NewController(t)
		generatorMock = NewMockUUIDGenerator(ctrl)
	)

	for _, h := range mustParseHashes(t, testCases) {
		h := h
		generatorMock.EXPECT().Generate().DoAndReturn(func() UUID { return h })
	}

	var updatedCh = make(chan struct{})

	storage := NewStorageFuncWrapper(NewMemoryStorage(), func() { updatedCh <- struct{}{} })

	hashTTL := 500 * time.Millisecond

	updater := NewUpdater(storage, generatorMock, hashTTL, zap.NewNop())

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go func() {
		err := updater.Run(ctx)
		assert.Nil(t, err, "updater error")
	}()

	for _, testCase := range testCases {
		ticker := time.NewTicker(hashTTL * 2)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			assert.Fail(t, "did not obtain hash regeneration signal")
		case <-updatedCh:
			row := storage.Get()
			assert.Equal(t, testCase.hash, row.Hash.String())
		}
	}
}

func mustParseHashes(t *testing.T, testCases []struct{ hash string }) []UUID {
	hashes := make([]UUID, 0, len(testCases))

	for _, testCase := range testCases {
		hash, err := uuid.Parse(testCase.hash)
		assert.Nil(t, err, "could not parse test case hash")

		hashes = append(hashes, hash)
	}

	return hashes
}
