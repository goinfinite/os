package infraHelper

import (
	"errors"
	"os"
	"path/filepath"
)

func IsSymlink(linkPath string) (bool, error) {
	fileExists := FileExists(linkPath)
	if !fileExists {
		return false, errors.New("FileDoesNotExists")
	}

	linkInfo, err := os.Lstat(linkPath)
	if err != nil {
		return false, err
	}

	isSymlink := linkInfo.Mode()&os.ModeSymlink == os.ModeSymlink
	if !isSymlink {
		return false, nil
	}

	return true, nil
}

func IsSymlinkTo(linkPath string, targetPath string) (bool, error) {
	isSymlink, err := IsSymlink(linkPath)
	if err != nil {
		return false, err
	}

	fileExists := FileExists(targetPath)
	if !fileExists {
		return false, errors.New("FileDoesNotExists")
	}

	if !isSymlink {
		return false, nil
	}

	linkTarget, err := os.Readlink(linkPath)
	if err != nil {
		return false, err
	}

	absTargetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return false, err
	}

	absLinkTarget, err := filepath.Abs(linkTarget)
	if err != nil {
		return false, err
	}

	return absLinkTarget == absTargetPath, nil
}
