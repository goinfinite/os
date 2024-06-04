package voTestHelpers

import (
	"encoding/base64"
)

func GenerateString(desiredSize int) string {
	desiredSizeBytesLength := float64(desiredSize) * 3
	desiredSizeStringLength := desiredSizeBytesLength / 4
	randomBytes := make([]byte, uint(desiredSizeStringLength))
	return base64.StdEncoding.EncodeToString(randomBytes)
}
