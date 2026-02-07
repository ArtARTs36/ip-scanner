package storage

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/artarts36/ip-scanner/internal/metrics"
)

type AutoUpdater struct {
	coldStorage Storage
	hotStorage  Storage

	cfg     Config
	metrics *metrics.Metrics

	currentDB *DB
}

func NewAutoUpdater(storages map[Type]Storage, cfg Config, metr *metrics.Metrics) *AutoUpdater {
	return &AutoUpdater{
		hotStorage:  NewOrderable(storages, cfg.Order),
		coldStorage: storages[cfg.AutoUpdate.Order],
		cfg:         cfg,
		metrics:     metr,
		currentDB:   &DB{},
	}
}

func (au *AutoUpdater) GetActualTime(ctx context.Context) (time.Time, error) {
	return au.hotStorage.GetActualTime(ctx)
}

func (au *AutoUpdater) Load(ctx context.Context) (*DB, error) {
	db, err := au.hotStorage.Load(ctx)
	if err != nil {
		return nil, err
	}

	au.currentDB.update(db)

	return db, nil
}

func (au *AutoUpdater) Run() {
	t := time.Tick(au.cfg.AutoUpdate.Interval)

	ctx := context.Background()

	slog.InfoContext(ctx, "[auto-update] run in background")

	for range t {
		err := au.run(ctx)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[auto-update] %s", err.Error()))
		}

		au.metrics.MMDBAutoUpdate.IncCycles()
		au.metrics.MMDBAutoUpdate.UpdateLastCycle()
	}
}

func (au *AutoUpdater) run(ctx context.Context) error {
	slog.InfoContext(ctx, "[auto-update] check exists new db")

	actualTime, err := au.coldStorage.GetActualTime(ctx)
	if err != nil {
		return fmt.Errorf("get actual time from storage: %w", err)
	}

	if au.currentDB.ActualTime.Equal(actualTime) || au.currentDB.ActualTime.After(actualTime) {
		slog.InfoContext(ctx, "[auto-update] storage doesnt have new db file")
		return nil
	}

	au.metrics.MMDBAutoUpdate.IncStartedLoads()

	slog.InfoContext(ctx, "[auto-update] loading new db")

	db, err := au.coldStorage.Load(ctx)
	if err != nil {
		au.metrics.MMDBAutoUpdate.IncCompletedLoads(false)
		return fmt.Errorf("load db: %w", err)
	}

	slog.InfoContext(ctx, "[auto-update] new db loaded")

	au.metrics.MMDBAutoUpdate.IncCompletedLoads(true)

	au.currentDB.update(db)
	au.metrics.MMDB.Scrap(db.ActualTime, db.DB.Metadata.DatabaseType, db.DB.Metadata.Languages, db.Size)

	return nil
}
