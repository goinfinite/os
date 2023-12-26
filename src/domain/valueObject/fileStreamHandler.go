package valueObject

import (
	"errors"
	"mime/multipart"
)

type FileStreamHandler struct {
	Name UnixFileName
	Size Byte
	Open func() (multipart.File, error)
}

func NewFileStreamHandler(value *multipart.FileHeader) (FileStreamHandler, error) {
	fileName, err := NewUnixFileName(value.Filename)
	if err != nil {
		return FileStreamHandler{}, errors.New("InvalidFileName")
	}

	fileSize := Byte(value.Size)
	isTooBig := fileSize.ToGiB() > 5
	if isTooBig {
		return FileStreamHandler{}, errors.New("FileIsTooBig")
	}

	return FileStreamHandler{
		Name: fileName,
		Size: fileSize,
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
