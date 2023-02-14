package hash

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var testCases = []struct{ hash string }{
	{
		hash: "f47ac10b-58cc-4372-0567-0e02b2c3d479",
	},
	{
		hash: "f47ac10b-58cc-f372-8567-0e02b2c3d471",
	},
}

func TestUpdater_Run(t *testing.T) {
	generatorMock := &uuidGeneratorMock{}
	generatorMock.Set(mustParseHashes(t, testCases))

	var updatedCh = make(chan struct{})

	storage := newStorageWrapper(NewMemoryStorage(), func() { updatedCh <- struct{}{} })

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

type uuidGeneratorMock struct {
	sync.Mutex
	hashes  []UUID
	counter int
}

func (g *uuidGeneratorMock) Set(hashes []UUID) {
	g.Lock()
	defer g.Unlock()
	g.hashes = hashes
	g.counter = 0
}

func (g *uuidGeneratorMock) Generate() UUID {
	if g.counter >= len(g.hashes) {
		panic("out of hashes")
	}

	g.counter++

	return g.hashes[g.counter-1]
}

type storageWrapper struct {
	storage Storage
	call    func()
}

func newStorageWrapper(storage Storage, call func()) Storage {
	return &storageWrapper{
		storage: storage,
		call:    call,
	}
}

func (s *storageWrapper) Upsert(row Row) error {
	if err := s.storage.Upsert(row); err != nil {
		return err
	}

	s.call()

	return nil
}

func (s *storageWrapper) Get() Row {
	return s.storage.Get()
}
