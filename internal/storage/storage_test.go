package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestMinioStorage(t *testing.T) *MinioStorage {
	t.Helper()
	s3Endpoint := "localhost:9000"
	s3AccessKeyID := "salty-spitoon"
	s3SecretAccessKey := "@salty-spitoon"
	return New(s3Endpoint, s3AccessKeyID, s3SecretAccessKey)
}

func TestStorage(t *testing.T) {
	minioStorage := setupTestMinioStorage(t)

	t.Run("Upload_FromFile", func(t *testing.T) {
		bucket := "images"
		localPath := "testdata/sample.jpg"
		remotePath := "sample.jpg"

		location, err := minioStorage.UploadFile(context.TODO(), bucket, localPath, remotePath)
		assert.Nil(t, err)
		assert.NotEmpty(t, location)
	})
}
