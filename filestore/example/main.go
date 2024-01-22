package main

import (
	"context"

	"bitbucket.org/junglee_games/getsetgo/configs"
	"bitbucket.org/junglee_games/getsetgo/filestore/fsfactory"
	"bitbucket.org/junglee_games/getsetgo/logger"
)

func main() {
	ctx := context.TODO()
	fs, err := fsfactory.GetFileStore(&configs.DefaultFilestoreConfig{Name: "LOCAL_FILE_STORE", Local: configs.DefaultLocalFileStoreConfig{DirectoryPath: "./testdocs"}})
	if err != nil {
		logger.Panic(ctx, "%v", err.Error())
	}
	files, err := fs.ListFiles(ctx, "pending", 10)
	if err != nil {
		logger.Panic(ctx, "%v", err.Error())
	}
	logger.Info(ctx, "%v", files)
	err = fs.RenameFile(ctx, "pending/text.txt", "text.txt")
	if err != nil {
		logger.Panic(ctx, "%v", err.Error())
	}
	logger.Info(ctx, "rename success")
}
