package localfilestore

import (
	"context"
	"os"
	"sync"

	"bitbucket.org/junglee_games/getsetgo/filestore"
	"bitbucket.org/junglee_games/getsetgo/logger"
	"bitbucket.org/junglee_games/getsetgo/utils/files"
)

type Config interface {
	GetDirectoryPath() string
}
type LocalFileStore struct {
	config   Config
	jobQueue chan *filestore.FileData
	//AckChan will receive a *filestore.FileData once its upload is completed by SaveAsync method
	ackChan chan *filestore.FileData
	wg      *sync.WaitGroup
}

func NewLocalFileStore(config Config) *LocalFileStore {
	fs := LocalFileStore{config: config}

	fs.jobQueue = make(chan *filestore.FileData, 100)

	fs.wg = new(sync.WaitGroup)
	fs.wg.Add(1)
	go fs.startUploader()

	return &fs

}

func (fs *LocalFileStore) Save(filesData *filestore.FileData) (string, error) {
	err := files.SaveFile(fs.config.GetDirectoryPath(), filesData.Name, filesData.Data)
	if err != nil {
		return "", err
	}

	return fs.config.GetDirectoryPath() + filesData.Name, nil
}

// it will upload data to upload queue and ack will receive on ack chan
//it push the object into queue
func (fs *LocalFileStore) SaveAsync(filesData *filestore.FileData) {
	fs.jobQueue <- filesData
}

// it will close the uploading queue Note : don't call SaveAsync method after calling this method
func (fs *LocalFileStore) StopUploading() {
	close(fs.jobQueue)
}

// it will wait until all the workers finish its uploading
func (fs *LocalFileStore) WaitForFinishingUpload() {
	fs.wg.Wait()
	if fs.ackChan != nil {
		close(fs.ackChan)
	}
}

func (fs *LocalFileStore) ack(filesData *filestore.FileData) {
	if fs.ackChan != nil {
		fs.ackChan <- filesData
	}
}

// it will take filesData from jobQueue and upload to s3 and ack by using ackqueue
func (fs *LocalFileStore) startUploader() {
	for filesData := range fs.jobQueue {
		var path string
		var err error
		path, err = fs.Save(filesData)

		if err != nil {
			filesData.UploadInfo.Error = err
			fs.ack(filesData)
			continue
		}
		filesData.UploadInfo.Path = path
		fs.ack(filesData)
		continue
	}
	fs.wg.Done()
	logger.Info(context.Background(), "jobQueue of localfilestore uploader closed , exiting startUploader function ")
}

// download file
func (fs *LocalFileStore) DownloadFile(filename string) ([]byte, error) {

	return os.ReadFile(fs.config.GetDirectoryPath() + filename)

}

// it will return ack chan
func (fs *LocalFileStore) GetAckChan(size int) chan *filestore.FileData {
	if fs.ackChan == nil {
		fs.ackChan = make(chan *filestore.FileData, size)
	}
	return fs.ackChan
}

//it will temporarily save file and returns signed url
func (fs *LocalFileStore) TemporarySave(filesData *filestore.FileData, expriryMinutes int) (string, error) {
	err := files.SaveFile(fs.config.GetDirectoryPath(), filesData.Name, filesData.Data)
	if err != nil {
		return "", err
	}

	return fs.config.GetDirectoryPath() + filesData.Name, nil
}

// it will return signed url for already uploaded file
func (fs *LocalFileStore) GetSignedURL(filepath string, expriryMinutes int) (string, error) {
	return filepath, nil
}
