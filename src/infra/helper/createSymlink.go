package infraHelper

import (
	"os"
)

func CreateSymlink(
	pkiSourcePath string,
	pkiTargetPath string,
	shouldOverwrite bool,
) error {
	if shouldOverwrite {
		err := os.Remove(pkiTargetPath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	err := os.Symlink(pkiSourcePath, pkiTargetPath)
	if err != nil {
		return err
	}

	return nil
}
