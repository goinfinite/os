package infraHelper

import (
	"os"
)

func CreateSymlink(
	sourcePath string,
	targetPath string,
	shouldOverwrite bool,
) error {
	if shouldOverwrite {
		err := os.Remove(targetPath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	err := os.Symlink(sourcePath, targetPath)
	if err != nil {
		return err
	}

	return nil
}
