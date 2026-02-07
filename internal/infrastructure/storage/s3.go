package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/oschwald/maxminddb-golang/v2"
)

type S3 struct {
	client *minio.Client
	cfg    S3Config
}

type S3Config struct {
	Endpoint        string `env:"ENDPOINT,required"`
	AccessKeyID     string `env:"ACCESS_KEY_ID,required,file,unset"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY,required,file,unset"`
	UseSSL          bool   `env:"USE_SSL,required"`
	Region          string `env:"REGION,required"`
	Bucket          string `env:"BUCKET,required"`
}

func NewS3(cfg S3Config) (*S3, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Region: cfg.Region,
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("init minio client: %w", err)
	}

	slog.Info("[storage] connected to s3", slog.String("endpoint", cfg.Endpoint))

	return &S3{
		client: minioClient,
		cfg:    cfg,
	}, nil
}

func (s *S3) GetActualTime(ctx context.Context) (time.Time, error) {
	t, _, err := s.getActual(ctx)
	if err != nil {
		return time.Time{}, err
	}

	return t, err
}

func (s *S3) Load(ctx context.Context) (*DB, error) {
	actualTime, key, err := s.getActual(ctx)
	if err != nil {
		return nil, fmt.Errorf("get actual key: %w", err)
	}

	obj, err := s.client.GetObject(ctx, s.cfg.Bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}

	stat, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("get object stat: %w", err)
	}

	slog.InfoContext(ctx, "[storage-s3] loading object", slog.String("key", key), slog.Int64("size", stat.Size))

	content, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("read object: %w", err)
	}

	db, err := maxminddb.FromBytes(content)
	if err != nil {
		return nil, err
	}

	return &DB{
		DB:         db,
		ActualTime: actualTime,
		Size:       stat.Size,
	}, nil
}

func (s *S3) getActual(ctx context.Context) (t time.Time, key string, err error) {
	actualTime := time.Time{}
	actualKey := ""

	for obj := range s.client.ListObjects(ctx, s.cfg.Bucket, minio.ListObjectsOptions{}) {
		if obj.Err != nil {
			return time.Time{}, "", obj.Err
		}

		var unix int64

		unix, err = strconv.ParseInt(obj.Key, 10, 64)
		if err != nil {
			slog.WarnContext(
				ctx,
				"[storage-s3] skipped file, because failed to parse filename",
				slog.Any("err", err),
				slog.String("object_key", obj.Key),
			)

			continue
		}

		t = time.Unix(unix, 0)
		if t.After(actualTime) {
			actualTime = t
			actualKey = obj.Key
		}
	}

	if actualKey == "" {
		return time.Time{}, "", errors.New("bucket is empty")
	}

	return actualTime, actualKey, nil
}
