package schemas

import "github.com/google/uuid"

type DownloadBatchFileSchema struct {
	FileIds uuid.UUIDs `json:"file_ids"`
}
