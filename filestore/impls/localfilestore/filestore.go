package localfilestore

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
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
// it push the object into queue
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

// it will temporarily save file and returns signed url
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

// it will list files in current directory
func (fs *LocalFileStore) ListFiles(ctx context.Context, folder string, limit int64) ([]string, error) {
	files, err := os.ReadDir(fs.config.GetDirectoryPath() + "/" + folder)
	if err != nil {
		logger.Error(ctx, "failed to read directory: %v", err.Error())
		return nil, fmt.Errorf("failed to list files: %v", err)
	}

	var result []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		result = append(result, file.Name())
		if int64(len(result)) >= limit {
			break
		}
	}

	return result, nil
}

// it will rename a file in current  directory
func (fs *LocalFileStore) RenameFile(ctx context.Context, oldname string, newname string) error {
	err := os.Rename(fs.config.GetDirectoryPath()+"/"+oldname, fs.config.GetDirectoryPath()+"/"+newname)
	if err != nil {
		logger.Error(ctx, "failed to rename file %v", err.Error())
		return fmt.Errorf("failed to rename file: %v", err)
	}
	return nil
}

func (fs *LocalFileStore) GetFileStream(filename string) (io.ReadCloser, error) {
	file, err := os.Open(fs.config.GetDirectoryPath() + "/" + filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (fs *LocalFileStore) DownloadFileToLocal(filename string, localPath string) error {
	// Open the source file in the local file store
	srcFile, err := os.Open(path.Join(fs.config.GetDirectoryPath(), filename))
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer srcFile.Close()

	// Create the destination file on the local system
	destFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer destFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("error copying file contents: %v", err)
	}
	return nil
}
