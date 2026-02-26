package infraHelper

func CopyFile(srcPath string, dstPath string) error {
	return fileClerk.CopyFile(srcPath, dstPath)
}
