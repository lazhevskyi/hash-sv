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
