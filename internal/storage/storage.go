package storage

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	s3AccessKeyID     = os.Getenv("S3_ACCESS_KEY_ID")
	s3SecretAccessKey = os.Getenv("S3_SECRET_ACCESS_KEY")
	s3Endpoint        = os.Getenv("S3_ENDPOINT")
)

type MinioStorage struct {
	client *minio.Client
}

func New(
	s3Endpoint, s3AccessKeyID, s3SecretAccessKey string,
) *MinioStorage {

	mc, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3AccessKeyID, s3SecretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		log.Fatal(err)
	}

	return &MinioStorage{
		client: mc,
	}
}

func (s *MinioStorage) UploadFile(ctx context.Context, bucket, localPath, remotePath string) (string, error) {
	info, err := s.client.FPutObject(ctx, bucket, remotePath, localPath, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return info.Key, nil
}
