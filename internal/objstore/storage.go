package objstore

import (
	"context"
	"mime/multipart"
)

type Storage interface {
	Upload(ctx context.Context, objName string, fileHeader *multipart.FileHeader) (string, error)
	Download(ctx context.Context, objName string) ([]byte, error)
	Delete(ctx context.Context, objName string) error
}

// type UploadResponse struct {
// 	ObjectName string
// 	StorageUrl string
// }
