package infraHelper

import (
	"bufio"
	"os"
)

func UpdateFile(filePath string, content string, shouldOverwrite bool) error {
	fileFlags := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	if shouldOverwrite {
		fileFlags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}

	file, err := os.OpenFile(filePath, fileFlags, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
