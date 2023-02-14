package hash

import "context"

type Service struct {
	updater Updater
	storage Storage
}

func NewService(
	updater Updater,
	storage Storage,
) *Service {
	return &Service{
		updater: updater,
		storage: storage,
	}
}

func (s *Service) Get() Row {
	return s.storage.Get()
}

func (s *Service) Run(ctx context.Context) error {
	return s.updater.Run(ctx)
}
