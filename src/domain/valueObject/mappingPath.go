package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const mappingPathRegex string = `^[^\s<>;'":#{}?\[\]]{1,512}$`

type MappingPath string

func NewMappingPath(value string) (MappingPath, error) {
	hasLeadingSlash := strings.HasPrefix(value, "/")
	if !hasLeadingSlash {
		value = "/" + value
	}

	mappingPath := MappingPath(value)
	if !mappingPath.isValid() {
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

func (mappingPath MappingPath) isValid() bool {
	re := regexp.MustCompile(mappingPathRegex)
	return re.MatchString(string(mappingPath))
}

func (mappingPath MappingPath) String() string {
	return string(mappingPath)
}
