package valueObject

import "errors"

type UnixCompressionType string

func NewUnixCompressionType(value string) (UnixCompressionType, error) {
	unixCompressionType := UnixCompressionType(value)
	if !unixCompressionType.isValid() {
		return "", errors.New("InvalidUnixCompressionType")
	}
	return unixCompressionType, nil
}

func NewUnixCompressionTypePanic(value string) UnixCompressionType {
	unixCompressionType, err := NewUnixCompressionType(value)
	if err != nil {
		panic(err)
	}
	return unixCompressionType
}

func (unixCompressionType UnixCompressionType) isValid() bool {
	unixCompressionTypeStr := string(unixCompressionType)
	return unixCompressionTypeStr == "gzip" || unixCompressionTypeStr == "zip"
}

func (unixCompressionType UnixCompressionType) String() string {
	return string(unixCompressionType)
}
