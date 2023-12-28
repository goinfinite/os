package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CompressUnixFiles struct {
	SourcePaths     []valueObject.UnixFilePath       `json:"sourcePaths"`
	DestinationPath valueObject.UnixFilePath         `json:"destinationPath"`
	CompressionType *valueObject.UnixCompressionType `json:"compressionType"`
}

func NewCompressUnixFiles(
	sourcePaths []valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
	compressionType *valueObject.UnixCompressionType,
) CompressUnixFiles {
	return CompressUnixFiles{
		SourcePaths:     sourcePaths,
		DestinationPath: destinationPath,
		CompressionType: compressionType,
	}
}
