package valueObject

import (
	"errors"
	"mime/multipart"
)

type FileStreamHandler struct {
	name UnixFileName
	size Byte
	Open func() (multipart.File, error)
}

func NewFileStreamHandler(value *multipart.FileHeader) (FileStreamHandler, error) {
	fileName, err := NewUnixFileName(value.Filename)
	if err != nil {
		return FileStreamHandler{}, errors.New("InvalidFileName")
	}

	fileSize := Byte(value.Size)
	if isTooBig(fileSize) {
		return FileStreamHandler{}, errors.New("FileIsTooBig")
	}

	return FileStreamHandler{
		name: fileName,
		size: fileSize,
		Open: value.Open,
	}, nil
}

func NewFileStreamHandlerPanic(value *multipart.FileHeader) FileStreamHandler {
	fileStreamHandler, err := NewFileStreamHandler(value)
	if err != nil {
		panic(err)
	}

	return fileStreamHandler
}

func isTooBig(fileSize Byte) bool {
	return fileSize.ToGiB() > 5
}

func (fileStreamHandler FileStreamHandler) GetFileName() UnixFileName {
	return fileStreamHandler.name
}

func (fileStreamHandler FileStreamHandler) GetFileSize() Byte {
	return fileStreamHandler.size
}
