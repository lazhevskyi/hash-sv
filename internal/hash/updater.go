package hash

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Updater interface {
	Run(ctx context.Context) error
}

type updater struct {
	storage   StorageUpdater
	generator UUIDGenerator
	ttl       time.Duration
	logger    *zap.Logger
}

func NewUpdater(
	storage StorageUpdater,
	generator UUIDGenerator,
	ttl time.Duration,
	logger *zap.Logger,
) Updater {
	return &updater{
		storage:   storage,
		generator: generator,
		ttl:       ttl,
		logger:    logger,
	}
}

func (u *updater) Run(ctx context.Context) error {
	if err := u.update(); err != nil {
		return fmt.Errorf("couln`t update hash: %w", err)
	}

	ticker := time.NewTicker(u.ttl)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := u.update(); err != nil {
				return fmt.Errorf("couln`t update hash: %w", err)
			}
		}
	}
}

func (u *updater) update() error {
	u.logger.Debug("going update hash")

	row := NewRow(u.generator.Generate())

	if err := u.storage.Upsert(row); err != nil {
		return fmt.Errorf("couln`t upsert hash: %w", err)
	}

	return nil
}
