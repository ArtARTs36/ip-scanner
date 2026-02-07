package storage

import "fmt"

func Create(cfg Config) (map[Type]Storage, error) {
	storages := map[Type]Storage{}

	if cfg.S3 != nil {
		var err error
		storages[TypeS3], err = NewS3(*cfg.S3)
		if err != nil {
			return nil, fmt.Errorf("create s3: %w", err)
		}
	}

	if cfg.Local != nil {
		var err error
		storages[TypeLocal], err = NewLocal(*cfg.Local)
		if err != nil {
			return nil, fmt.Errorf("create local: %w", err)
		}
	}

	return storages, nil
}
