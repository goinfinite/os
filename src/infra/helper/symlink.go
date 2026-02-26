package infraHelper

func IsSymlink(sourcePath string) bool {
	return fileClerk.IsSymlink(sourcePath)
}

func IsSymlinkTo(sourcePath string, targetPath string) bool {
	return fileClerk.IsSymlinkTo(sourcePath, targetPath)
}

func CreateSymlink(sourcePath string, targetPath string, shouldOverwrite bool) error {
	return fileClerk.CreateSymlink(sourcePath, targetPath, shouldOverwrite)
}

func RemoveSymlink(symlinkPath string) error {
	return fileClerk.RemoveSymlink(symlinkPath)
}
