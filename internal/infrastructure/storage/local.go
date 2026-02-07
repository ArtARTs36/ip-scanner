package storage

import (
	"context"
	"os"
	"time"

	"github.com/oschwald/maxminddb-golang/v2"
)

type Local struct {
	path string

	actualTime time.Time
	stat       os.FileInfo
}

type LocalConfig struct {
	Path string `env:"PATH,required"`
}

func NewLocal(cfg LocalConfig) (*Local, error) {
	l := &Local{
		path: cfg.Path,
	}

	stat, err := os.Stat(cfg.Path)
	if err != nil {
		return nil, err
	}
	l.stat = stat

	return l, nil
}

func (l *Local) GetActualTime(_ context.Context) (time.Time, error) {
	return l.actualTime, nil
}

func (l *Local) Load(_ context.Context) (*DB, error) {
	db, err := maxminddb.Open(l.path)
	if err != nil {
		return nil, err
	}

	l.actualTime = time.Unix(int64(db.Metadata.BuildEpoch), 0) //nolint:gosec // not need

	return &DB{
		DB:         db,
		ActualTime: l.actualTime,
		Size:       l.stat.Size(),
	}, nil
}
