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

	return os.Symlink(sourcePath, targetPath)
}
