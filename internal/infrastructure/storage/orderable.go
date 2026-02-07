package storage

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

type Orderable struct {
	storages []Storage
}

func NewOrderable(storages map[Type]Storage, order []Type) *Orderable {
	storage := &Orderable{
		storages: make([]Storage, len(order)),
	}

	for i, typ := range order {
		storage.storages[i] = storages[typ]
	}

	return storage
}

func (s *Orderable) GetActualTime(ctx context.Context) (time.Time, error) {
	errs := []error{}

	for _, storage := range s.storages {
		t, err := storage.GetActualTime(ctx)
		if err != nil {
			errs = append(errs, err)
			slog.ErrorContext(ctx, "[storage-orderable] failed to get actual time", slog.Any("err", err))
			continue
		}

		return t, nil
	}

	return time.Time{}, errors.Join(errs...)
}

func (s *Orderable) Load(ctx context.Context) (*DB, error) {
	errs := []error{}

	for _, storage := range s.storages {
		db, err := storage.Load(ctx)
		if err != nil {
			errs = append(errs, err)
			slog.ErrorContext(ctx, "[storage-orderable] failed to load", slog.Any("err", err))
			continue
		}

		return db, nil
	}

	return nil, errors.Join(errs...)
}
