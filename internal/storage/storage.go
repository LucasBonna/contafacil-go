package storage

import (
	"io"

	"github.com/google/uuid"
)

type StorageManager interface {
  Upload(file io.Reader, fileId uuid.UUID) error
  Download(fileId uuid.UUID) ([]byte, error)
}

func SetStorage(storage StorageManager) StorageManager {
  return storage
}
