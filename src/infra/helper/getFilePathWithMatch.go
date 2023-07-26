package infraHelper

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func GetFilePathWithMatch(dir string, partialMatch string) (string, error) {
	var result string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(info.Name(), partialMatch) {
			result = path
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return "", errors.New("FileNotFound")
	}

	return result, nil
}
