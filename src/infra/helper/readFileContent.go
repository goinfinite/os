package infraHelper

func ReadFileContent(filePath string) (string, error) {
	return fileClerk.ReadFileContent(filePath, nil)
}
