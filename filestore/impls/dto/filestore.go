package dto

import "time"

type ListResponse struct {
	FileInfo []FileInfo
}

type FileInfo struct {
	Name         string
	UploadedTime time.Time
}
