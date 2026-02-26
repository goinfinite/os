package infraHelper

func MakeDir(dirPath string) error {
	return fileClerk.CreateDir(dirPath)
}
