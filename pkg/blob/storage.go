package blob

import (
	"context"
	"io"

	"doublequote/pkg/domain"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
)

var _ domain.StorageService = (*StorageService)(nil)

type StorageService struct {
	svc        *blob.Bucket
	bucketName string
}

func NewStorageService(cfg domain.Config) (*StorageService, func() error, error) {
	s := &StorageService{}

	bucket, err := blob.OpenBucket(context.Background(), "file://"+cfg.App.DataFolder)
	if err != nil {
		return nil, nil, err
	}

	s.svc = bucket

	return s, bucket.Close, nil
}

func (s *StorageService) Get(ctx context.Context, key string) (io.Reader, error) {
	obj, err := s.svc.NewReader(ctx, key, nil)
	return obj, err
}

func (s *StorageService) Set(ctx context.Context, key string, value []byte) error {
	obj, err := s.svc.NewWriter(ctx, key, nil)
	if err != nil {
		return err
	}
	_, err = obj.Write(value)
	if err != nil {
		return err
	}
	return obj.Close()
}
