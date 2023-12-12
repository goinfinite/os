package valueObject

import "mime/multipart"

type MultipartFile struct {
	CustomName          string
	MultipartFileHeader *multipart.FileHeader
}

func NewMultipartFile(
	customName string,
	multipartFileHeader *multipart.FileHeader,
) MultipartFile {
	return MultipartFile{
		CustomName:          customName,
		MultipartFileHeader: multipartFileHeader,
	}
}
