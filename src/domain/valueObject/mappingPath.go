package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const mappingPathRegex string = `^[^\s<>;'":#{}?\[\]]{1,512}$`

type MappingPath string

func NewMappingPath(value string) (MappingPath, error) {
	if len(value) == 0 {
		value = "/"
	}

	startsWithTrailingSlash := strings.HasPrefix(value, "/")
	if !startsWithTrailingSlash {
		value = "/" + value
	}

	mappingPath := MappingPath(value)
	if !mappingPath.isValid(value) {
		return "", errors.New("InvalidMappingPath")
	}

	return mappingPath, nil
}

func NewMappingPathPanic(value string) MappingPath {
	mappingPath, err := NewMappingPath(value)
	if err != nil {
		panic(err)
	}

	return mappingPath
}

func (MappingPath) isValid(value string) bool {
	re := regexp.MustCompile(mappingPathRegex)
	return re.MatchString(value)
}

func (mappingPath MappingPath) String() string {
	return string(mappingPath)
}
