// MinioConfig is the configuration for Minio storage.
// It includes the endpoint, access key, secret key, and bucket name.

package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	Client *minio.Client
	Bucket string
}

func New(url, key, secret, bucket string) (*Minio, error) {
	cli, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(key, secret, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &Minio{Client: cli, Bucket: bucket}, nil
}

func (m *Minio) Upload(ctx context.Context, name string, r io.Reader, size int64) error {
	_, err := m.Client.PutObject(ctx, m.Bucket, name, r, size, minio.PutObjectOptions{})
	return err
}
