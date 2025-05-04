package objstore

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
)

type GCStorage struct {
	bucketName string
	projectId  string
	client     *storage.Client
}

func NewGCStorage(projectId, bucketName string) (Storage, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to setup the storage client: %v", err)
	}

	return &GCStorage{
		bucketName,
		projectId,
		client,
	}, nil
}

func (g *GCStorage) CreateSignedURL(unsignedURL string) (string, error) {
	objectName, err := objectNameFromURL(unsignedURL)
	if err != nil {
		return "", err
	}

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(60 * time.Minute), // valid for 1 hour
	}
	return storage.SignedURL(g.bucketName, objectName, opts)
}

func (g *GCStorage) Download(ctx context.Context, objName string) ([]byte, error) {
	objectHandle := g.client.Bucket(g.bucketName).Object(objName)
	reader, err := objectHandle.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %w", objName, err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to read from the object handle reader : %v", err)
	}
	return data, nil
}

func (g *GCStorage) Upload(ctx context.Context, objName string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("unable to open the file:%v", err)
	}
	defer file.Close()

	writer := g.client.Bucket(g.bucketName).Object(objName).NewWriter(ctx)
	defer writer.Close()
	// set cache control so profile image will be served fresh by browsers
	// To do this with object handle, you'd first have to upload, then update
	writer.ObjectAttrs.CacheControl = "Cache-Control:no-cache, max-age=0"
	writer.ContentType = fileHeader.Header.Get("Content-Type")
	_, err = io.Copy(writer, file)
	if err != nil {
		return "", fmt.Errorf("unable to copy the file to storage:%v", err)
	}
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, objName), nil
}

func (g *GCStorage) Delete(ctx context.Context, objName string) error {
	return g.client.Bucket(g.bucketName).Object(objName).Delete(ctx)
}
