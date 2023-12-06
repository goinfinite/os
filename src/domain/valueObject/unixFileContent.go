package valueObject

import (
	"encoding/base64"
	"errors"
	"regexp"
)

const unixFileContentRegexExpression = `^(?:[A-Za-z0-9+\/]{4})*(?:[A-Za-z0-9+\/]{4}|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{2}={2})$`

type UnixFileContent string

func NewUnixFileContent(value string) (UnixFileContent, error) {
	unixFileContent := UnixFileContent(value)
	if !unixFileContent.isValid() {
		return "", errors.New("InvalidUnixFileContent")
	}
	return unixFileContent, nil
}

func NewUnixFileContentPanic(value string) UnixFileContent {
	unixFileContent, err := NewUnixFileContent(value)
	if err != nil {
		panic(err)
	}
	return unixFileContent
}

func (unixFileContent UnixFileContent) isValid() bool {
	isEmpty := false
	if len(unixFileContent) < 1 {
		isEmpty = true
	}

	unixFileContentRegex := regexp.MustCompile(unixFileContentRegexExpression)
	return unixFileContentRegex.MatchString(string(unixFileContent)) && !isEmpty
}

func (unixFileContent UnixFileContent) GetDecodedContent() (string, error) {
	decodedContent, err := base64.StdEncoding.DecodeString(string(unixFileContent))
	if err != nil {
		return "", err
	}

	return string(decodedContent), nil
}

func (unixFileContent UnixFileContent) String() string {
	return string(unixFileContent)
}
