package imgstore

import (
	"context"
	"mime/multipart"
)

type Storage interface {
	Upload(ctx context.Context, fileHeader *multipart.FileHeader) (string, error)
	Download(ctx context.Context, objName string) ([]byte, error)
	Update(ctx context.Context, objName string, fileHeader *multipart.FileHeader) (string, error)
	Delete(ctx context.Context, objName string) error
}
