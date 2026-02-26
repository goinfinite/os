package infraHelper

func UpdateFile(filePath string, content string, shouldOverwrite bool) error {
	return fileClerk.UpdateFileContent(filePath, content, shouldOverwrite)
}
