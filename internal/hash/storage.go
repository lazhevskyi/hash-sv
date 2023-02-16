package hash

import (
	"context"
	"sync"
)

type Storage interface {
	StorageGetter
	StorageUpdater
}

type StorageGetter interface {
	Get(ctx context.Context) Row
}

type StorageUpdater interface {
	Upsert(ctx context.Context, row Row) error
}

type memoryStorage struct {
	sync.RWMutex
	row Row
}

func NewMemoryStorage() Storage {
	return &memoryStorage{}
}

func (s *memoryStorage) Upsert(_ context.Context, row Row) error {
	s.Lock()
	defer s.Unlock()
	s.row = row

	return nil
}

func (s *memoryStorage) Get(_ context.Context) Row {
	s.RLock()
	row := s.row
	s.RUnlock()

	return row
}

type storageFuncWrapper struct {
	storage Storage
	call    func()
}

func NewStorageFuncWrapper(storage Storage, call func()) Storage {
	return &storageFuncWrapper{
		storage: storage,
		call:    call,
	}
}

func (s *storageFuncWrapper) Upsert(ctx context.Context, row Row) error {
	if err := s.storage.Upsert(ctx, row); err != nil {
		return err
	}

	s.call()

	return nil
}

func (s *storageFuncWrapper) Get(ctx context.Context) Row {
	return s.storage.Get(ctx)
}
