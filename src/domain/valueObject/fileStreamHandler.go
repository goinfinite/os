package valueObject

import (
	"errors"
	"mime/multipart"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type FileStreamHandler struct {
	Name tkValueObject.UnixFileName
	Size tkValueObject.Byte
	Open func() (multipart.File, error)
}

func NewFileStreamHandler(value *multipart.FileHeader) (
	fileStreamHandler FileStreamHandler, err error,
) {
	fileName, err := tkValueObject.NewUnixFileName(value.Filename, false)
	if err != nil {
		return fileStreamHandler, err
	}

	if value.Size > 5*1024*1024*1024 {
		return fileStreamHandler, errors.New("FileIsTooBig")
	}

	fileSize, err := tkValueObject.NewByte(value.Size)
	if err != nil {
		return fileStreamHandler, errors.New("InvalidFileSize")
	}

	return FileStreamHandler{
		Name: fileName,
		Size: fileSize,
		Open: value.Open,
	}, nil
}
