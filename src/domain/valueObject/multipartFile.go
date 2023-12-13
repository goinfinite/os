package valueObject

import (
	"errors"
	"mime/multipart"
)

type MultipartFile struct {
	name UnixFileName
	size Byte
	Open func() (multipart.File, error)
}

func NewMultipartFile(value *multipart.FileHeader) (MultipartFile, error) {
	fileName, err := NewUnixFileName(value.Filename)
	if err != nil {
		return MultipartFile{}, errors.New("InvalidFileName")
	}

	fileSize := Byte(value.Size)
	if isTooBig(fileSize) {
		return MultipartFile{}, errors.New("FileIsTooBig")
	}

	return MultipartFile{
		name: fileName,
		size: fileSize,
		Open: value.Open,
	}, nil
}

func NewMultipartFilePanic(value *multipart.FileHeader) MultipartFile {
	multipartFile, err := NewMultipartFile(value)
	if err != nil {
		panic(err)
	}

	return multipartFile
}

func isTooBig(fileSize Byte) bool {
	return fileSize.ToGiB() > 5
}

func (multipartFile MultipartFile) GetFileName() UnixFileName {
	return multipartFile.name
}

func (multipartFile MultipartFile) GetFileSize() Byte {
	return multipartFile.size
}
