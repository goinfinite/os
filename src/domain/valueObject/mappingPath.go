package valueObject

import (
	"errors"
	"regexp"
)

const mappingPathRegex string = `^[^\s<>;'":#{}?\[\]]{1,512}$`

type MappingPath string

func NewMappingPath(value string) (MappingPath, error) {
	mp := MappingPath(value)
	if !mp.isValid() {
		return "", errors.New("InvalidMappingPath")
	}
	return mp, nil
}

func NewMappingPathPanic(value string) MappingPath {
	mp, err := NewMappingPath(value)
	if err != nil {
		panic(err)
	}
	return mp
}

func (mp MappingPath) isValid() bool {
	re := regexp.MustCompile(mappingPathRegex)
	return re.MatchString(string(mp))
}

func (mp MappingPath) String() string {
	return string(mp)
}
