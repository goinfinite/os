package infra

import (
	"errors"
	"log"
	"os"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraFactory "github.com/speedianet/os/src/infra/factory"
)

type FilesQueryRepo struct{}

func (repo FilesQueryRepo) Get(
	unixFilePath valueObject.UnixFilePath,
) ([]entity.UnixFile, error) {
	unixFileList := []entity.UnixFile{}

	filePathInfo, err := os.Stat(unixFilePath.String())
	if err != nil {
		return unixFileList, errors.New("UnableToGetPathInfo")
	}

	if !filePathInfo.IsDir() {
		unixFile, err := infraFactory.UnixFileFactory(unixFilePath, filePathInfo)
		if err == nil {
			unixFileList = append(unixFileList, unixFile)
		}

		return unixFileList, nil
	}

	dirEntriesToAnalyzeList, err := os.ReadDir(unixFilePath.String())
	if err != nil {
		log.Printf("UnableToGetDirInfo: %s", err)
		return unixFileList, errors.New("UnableToGetDirInfo")
	}

	for _, dirEntry := range dirEntriesToAnalyzeList {
		inodeInfo, _ := dirEntry.Info()
		unixFile, err := infraFactory.UnixFileFactory(unixFilePath, inodeInfo)
		if err != nil {
			continue
		}

		unixFileList = append(unixFileList, unixFile)
	}

	return unixFileList, nil
}
