package objstore

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"path"
	"time"
)

type Storage interface {
	Upload(ctx context.Context, objName string, fileHeader *multipart.FileHeader) (string, error)
	Download(ctx context.Context, objName string) ([]byte, error)
	Delete(ctx context.Context, objName string) error
	CreateSignedURL(unsignedURL string, duration time.Duration) (string, error)
}

func objectNameFromURL(URL string) (string, error) {
	// split off the last part of the url
	urlPath, err := url.Parse(URL)
	if err != nil {
		return "", fmt.Errorf("failed to parse object name from URL")
	}
	return path.Base(urlPath.Path), nil
}

// type UploadResponse struct {
// 	ObjectName string
// 	StorageUrl string
//
