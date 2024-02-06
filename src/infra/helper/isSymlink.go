package infraHelper

import (
	"os"
	"path/filepath"
)

func IsSymlink(linkPath string) bool {
	linkInfo, err := os.Lstat(linkPath)
	if err != nil {
		return false
	}

	isSymlink := linkInfo.Mode()&os.ModeSymlink == os.ModeSymlink
	return isSymlink
}

func IsSymlinkTo(linkPath string, targetPath string) bool {
	isSymlink := IsSymlink(linkPath)
	if !isSymlink {
		return false
	}

	if !isSymlink {
		return false
	}

	linkTarget, err := os.Readlink(linkPath)
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
