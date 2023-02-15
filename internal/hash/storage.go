package hash

import "sync"

type Storage interface {
	StorageGetter
	StorageUpdater
}

type StorageGetter interface {
	Get() Row
}

type StorageUpdater interface {
	Upsert(row Row) error
}

type memoryStorage struct {
	sync.RWMutex
	row Row
}

func NewMemoryStorage() Storage {
	return &memoryStorage{}
}

func (s *memoryStorage) Upsert(row Row) error {
	s.Lock()
	defer s.Unlock()
	s.row = row

	return nil
}

func (s *memoryStorage) Get() Row {
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

func (s *storageFuncWrapper) Upsert(row Row) error {
	if err := s.storage.Upsert(row); err != nil {
		return err
	}

	s.call()

	return nil
}

func (s *storageFuncWrapper) Get() Row {
	return s.storage.Get()
}
