package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CompressUnixFiles struct {
	Paths           []valueObject.UnixFilePath      `json:"paths"`
	DestinationPath valueObject.UnixFilePath        `json:"destinationPath"`
	CompressionType valueObject.UnixCompressionType `json:"compressionType"`
}

func NewCompressUnixFiles(
	paths []valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
	compressionType valueObject.UnixCompressionType,
) CompressUnixFiles {
	return CompressUnixFiles{
		Paths:           paths,
		DestinationPath: destinationPath,
		CompressionType: compressionType,
	}
}
