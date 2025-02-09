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

func (g *GCStorage) Upload(ctx context.Context, fileHeader *multipart.FileHeader) (*UploadResponse, error) {
	// open the associated File
	srcFile, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open the file:%v", err)
	}

	defer srcFile.Close()
	// create a unique filename
	fileName := fmt.Sprintf("%s_%d", fileHeader.Filename, time.Now().UnixNano())

	// get the bucket handle
	bucket := g.client.Bucket(g.bucketName)
	objectHandle := bucket.Object(fileName)

	writer := objectHandle.NewWriter(ctx)
	writer.ContentType = fileHeader.Header.Get("Content-Type")

	// Copy the file to the Object
	_, err = io.Copy(writer, srcFile)
	if err != nil {
		return nil, fmt.Errorf("unable to copy to storage:%v", err)
	}
	defer writer.Close()
	storageUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, fileName)
	return &UploadResponse{
		StorageUrl: storageUrl,
		ObjectName: fileName,
	}, nil
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

func (g *GCStorage) Update(ctx context.Context, objName string, fileHeader *multipart.FileHeader) (string, error) {
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
	_, err = io.Copy(writer, file)
	if err != nil {
		return "", fmt.Errorf("unable to copy the file to storage:%v", err)
	}
	return "object updated", nil
}

func (g *GCStorage) Delete(ctx context.Context, objName string) error {
	return g.client.Bucket(g.bucketName).Object(objName).Delete(ctx)
}
