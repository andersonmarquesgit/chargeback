package minio

import (
	"batch/internal/infrastructure/objectstorage"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"os"
)

type BatchFileDownloader struct {
	client     *minio.Client
	bucketName string
}

func NewBatchFileDownloader(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (objectstorage.Downloader, error) {
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
	return &BatchFileDownloader{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (cbu *BatchFileDownloader) DownloadFile(localPath string, objectName string) error {
	ctx := context.Background()

	object, err := cbu.client.GetObject(ctx, cbu.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get object %s: %w", objectName, err)
	}
	defer object.Close()

	// Cria ou sobrescreve o arquivo local
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("could not create local file %s: %w", localPath, err)
	}
	defer file.Close()

	// Copia o conteúdo do objeto para o arquivo local
	_, err = io.Copy(file, object)
	if err != nil {
		return fmt.Errorf("failed to copy object to file: %w", err)
	}

	return nil
}
