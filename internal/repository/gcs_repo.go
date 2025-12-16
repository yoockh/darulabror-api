package repository

import (
	"context"
	"errors"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
)

var ErrStorageNotConfigured = errors.New("gcs storage is not configured")

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

func (r *gcpStorageRepo) validate() error {
	if r.client == nil || r.bucketName == "" {
		return ErrStorageNotConfigured
	}
	return nil
}

// UploadFile — handle file upload to GCS
func (r *gcpStorageRepo) UploadFile(ctx context.Context, file io.Reader, objectName string) (string, error) {
	if err := r.validate(); err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	obj := r.client.Bucket(r.bucketName).Object(objectName)
	writer := obj.NewWriter(ctx)

	if _, err := io.Copy(writer, file); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"bucket": r.bucketName,
			"object": objectName,
		}).Error("gcs upload failed")
		_ = writer.Close()
		return "", err
	}

	if err := writer.Close(); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"bucket": r.bucketName,
			"object": objectName,
		}).Error("gcs writer close failed")
		return "", err
	}

	// PUBLIC bucket → return URL
	if r.isPublic {
		url := "https://storage.googleapis.com/" + r.bucketName + "/" + objectName
		logrus.WithFields(logrus.Fields{
			"bucket": r.bucketName,
			"object": objectName,
			"url":    url,
		}).Info("public file uploaded to gcs")
		return url, nil
	}

	// PRIVATE bucket → return objectName
	logrus.WithFields(logrus.Fields{
		"bucket": r.bucketName,
		"object": objectName,
	}).Info("private file uploaded to gcs")
	return objectName, nil
}

// GenerateSignedURL — generate signed URL for private objects
func (r *gcpStorageRepo) GenerateSignedURL(ctx context.Context, objectName string, expire time.Duration) (string, error) {
	_ = ctx // reserved (kalau nanti mau pakai IAM SignBlob / per-request trace)

	if err := r.validate(); err != nil {
		return "", err
	}

	if r.isPublic {
		return "https://storage.googleapis.com/" + r.bucketName + "/" + objectName, nil
	}

	// NOTE: SignedURL butuh konfigurasi credentials signing (GoogleAccessID/PrivateKey atau IAM SignBlob).
	url, err := storage.SignedURL(r.bucketName, objectName, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(expire),
	})
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"bucket": r.bucketName,
			"object": objectName,
		}).Error("failed generating signed url")
		return "", err
	}

	return url, nil
}

// NOTE:
// - Kalau bucket public: GenerateSignedURL() cuma return public URL (signed URL tidak diperlukan).
// - Mode private belum dipakai sekarang, tapi disiapkan untuk kebutuhan future (restricted media).
