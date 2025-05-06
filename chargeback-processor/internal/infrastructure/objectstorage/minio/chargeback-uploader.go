package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"processor/internal/infrastructure/objectstorage"
)

var _ objectstorage.Uploader = (*ChargebackUploader)(nil)

type ChargebackUploader struct {
	client     *minio.Client
	bucketName string
}

func NewChargebackUploader(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (objectstorage.Uploader, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Check and create bucket if not exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		log.Printf("Bucket %s does not exist. Creating...", bucketName)
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		log.Printf("Bucket %s created successfully.", bucketName)
	}

	log.Printf("Bucket %s connect successfully.", bucketName)
	return &ChargebackUploader{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (cbu *ChargebackUploader) UploadFile(localPath string, objectName string) error {
	ctx := context.Background()

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", localPath, err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not stat file %s: %w", localPath, err)
	}

	_, err = cbu.client.PutObject(ctx, cbu.bucketName, objectName, file, stat.Size(), minio.PutObjectOptions{
		ContentType: "application/x-ndjson",
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to minio: %w", err)
	}

	return nil
}
