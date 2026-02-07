package storage

import (
	"log/slog"
	"time"

	"github.com/oschwald/maxminddb-golang/v2"
)

type DB struct {
	DB         *maxminddb.Reader
	ActualTime time.Time
	Size       int64
}

func (d *DB) update(newDB *DB) {
	prev := d.DB
	d.DB = newDB.DB
	d.ActualTime = newDB.ActualTime

	if prev != nil {
		time.AfterFunc(time.Minute, func() {
			if err := prev.Close(); err != nil {
				slog.Error("[db] failed to close db", slog.Any("err", err))
			}
		})
	}
}

func (d *DB) Close() error {
	if d.DB == nil {
		return nil
	}

	return d.DB.Close()
}
