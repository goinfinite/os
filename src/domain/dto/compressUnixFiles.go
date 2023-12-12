package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CompressUnixFiles struct {
	Paths           []valueObject.UnixFilePath      `json:"paths"`
	DestinationPath valueObject.UnixFilePath        `json:"destinationPath"`
	CompressionType valueObject.UnixCompressionType `json:"compressionType"`
}

func NewCompressUnixFiles(
	Paths []valueObject.UnixFilePath,
	DestinationPath valueObject.UnixFilePath,
	CompressionType valueObject.UnixCompressionType,
) CompressUnixFiles {
	return CompressUnixFiles{
		Paths:           Paths,
		DestinationPath: DestinationPath,
		CompressionType: CompressionType,
	}
}
