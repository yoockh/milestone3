package repository

import (
	"context"
	"io"
	"log"
	"time"

	"cloud.google.com/go/storage"
)

type GCPStorageRepo interface {
	UploadFile(ctx context.Context, file io.Reader, objectName string) (string, error)
	GenerateSignedURL(ctx context.Context, objectName string, expire time.Duration) (string, error)
}

type gcpStorageRepo struct {
	client     *storage.Client
	bucketName string
	isPublic   bool
}

func NewGCPStorageRepo(client *storage.Client, bucketName string, isPublic bool) GCPStorageRepo {
	return &gcpStorageRepo{
		client:     client,
		bucketName: bucketName,
		isPublic:   isPublic,
	}
}

// UploadFile — handle file upload to GCS
func (r *gcpStorageRepo) UploadFile(ctx context.Context, file io.Reader, objectName string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	obj := r.client.Bucket(r.bucketName).Object(objectName)
	writer := obj.NewWriter(ctx)

	// Copy file
	if _, err := io.Copy(writer, file); err != nil {
		log.Printf("GCS upload failed: %v", err)
		return "", err
	}

	// Close writer
	if err := writer.Close(); err != nil {
		log.Printf("GCS writer close failed: %v", err)
		return "", err
	}

	// PUBLIC bucket → return URL
	if r.isPublic {
		url := "https://storage.googleapis.com/" + r.bucketName + "/" + objectName
		log.Printf("Public file uploaded to: %s", url)
		return url, nil
	}

	// PRIVATE bucket → return objectName
	log.Printf("Private file uploaded as object: %s", objectName)
	return objectName, nil
}

// GenerateSignedURL — generate signed URL for private objects
func (r *gcpStorageRepo) GenerateSignedURL(ctx context.Context, objectName string, expire time.Duration) (string, error) {
	if r.isPublic {
		// Public bucket doesn't need signed URL
		url := "https://storage.googleapis.com/" + r.bucketName + "/" + objectName
		return url, nil
	}

	// PRIVATE → generate signed URL
	url, err := storage.SignedURL(r.bucketName, objectName, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(expire),
	})

	if err != nil {
		log.Printf("Failed generating signed URL: %v", err)
		return "", err
	}

	return url, nil
}
