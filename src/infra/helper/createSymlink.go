package infraHelper

import (
	"errors"
	"os"
)

func CreateSymlink(
	sourcePath string,
	targetPath string,
	shouldOverwrite bool,
) error {
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
