package filestore

import (
	"context"
	"io"
)

type FileData struct {
	Name       string
	Data       string
	UserId     string
	UploadInfo struct {
		Path  string
		Error error
	}
}

type FileStore interface {
	Save(filesData *FileData) (string, error)

	// it will upload data to upload queue and ack will receive on ack chan
	SaveAsync(filesData *FileData)

	// it will close the uploading queue Note : don't call SaveAsync method after calling this method
	StopUploading()

	// it will wait until all the workers finish its uploading
	WaitForFinishingUpload()

	// download file
	DownloadFile(filename string) ([]byte, error)
	DownloadFileToLocal(filename string, localPath string) error

	GetFileStream(filename string) (io.ReadCloser, error)

	// it will return ack chan
	GetAckChan(size int) chan *FileData

	//it will temporarily save file and returns signed url
	TemporarySave(filesData *FileData, expriryMinutes int) (string, error)

	// it will return signed url for already uploaded file
	GetSignedURL(filepath string, expriryMinutes int) (string, error)

	ListFiles(ctx context.Context, folder string, limit int64) ([]string, error)
	RenameFile(ctx context.Context, oldname string, newname string) error
}
