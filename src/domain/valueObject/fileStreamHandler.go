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
	fileName, err := NewUnixFileName(value.Filename)
	if err != nil {
		return fileStreamHandler, err
	}

	fileSize, err := tkValueObject.NewByte(value.Size)
	if err != nil {
		return fileStreamHandler, errors.New("InvalidFileSize")
	}

	if fileSize.ToGiB() > 5 {
		return fileStreamHandler, errors.New("FileIsTooBig")
	}

	return FileStreamHandler{
		Name: fileName,
		Size: fileSize,
		Open: value.Open,
	}, nil
}
