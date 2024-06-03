package infraHelper

import (
	"errors"
	"os"
	"path/filepath"
)

func IsSymlink(sourcePath string) bool {
	linkInfo, err := os.Lstat(sourcePath)
	if err != nil {
		return false
	}

	isSymlink := linkInfo.Mode()&os.ModeSymlink == os.ModeSymlink
	return isSymlink
}

func IsSymlinkTo(sourcePath string, targetPath string) bool {
	isSymlink := IsSymlink(sourcePath)
	if !isSymlink {
		return false
	}

	linkTarget, err := os.Readlink(sourcePath)
	if err != nil {
		return false
	}

	absTargetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return false
	}

	absLinkTarget, err := filepath.Abs(linkTarget)
	if err != nil {
		return false
	}

	return absLinkTarget == absTargetPath
}

func CreateSymlink(sourcePath string, targetPath string, shouldOverwrite bool) error {
	if !shouldOverwrite && !FileExists(sourcePath) {
		return errors.New("FileNotFound")
	}

	if shouldOverwrite {
		err := os.Remove(targetPath)
		if err != nil {
			return err
		}
	}

	return os.Symlink(sourcePath, targetPath)
}

func RemoveSymlink(symlinkPath string) error {
	return os.Remove(symlinkPath)
}
