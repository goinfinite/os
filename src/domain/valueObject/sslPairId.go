package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const sslPairIdExpression = "^[a-fA-F0-9]{64}$"

type SslPairId string

func NewSslPairId(value interface{}) (sslPairId SslPairId, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return sslPairId, errors.New("SslPairIdMustBeString")
	}

	re := regexp.MustCompile(sslPairIdExpression)
	if !re.MatchString(stringValue) {
		return sslPairId, errors.New("InvalidSslPairId")
	}

	return SslPairId(stringValue), nil
}

func NewSslPairIdFromSslPairContent(
	sslCertificate SslCertificateContent,
	sslChainCertificates []SslCertificateContent,
	sslPrivateKey SslPrivateKey,
) (sslPairId SslPairId, err error) {
	sslChainCertificatesMerged := ""
	for _, sslChainCertificate := range sslChainCertificates {
		sslChainCertificatesMerged += sslChainCertificate.String() + "\n"
	}
	contentToEncode := sslCertificate.String() + "\n" + sslChainCertificatesMerged + "\n" + sslPrivateKey.String()

	sslPairIdContent, err := voHelper.TransformContentIntoStrongHash(contentToEncode)
	if err != nil {
		return sslPairId, errors.New("InvalidSslPairIdFromSslPairContent")
	}
	return NewSslPairId(sslPairIdContent)
}

func (vo SslPairId) String() string {
	return string(vo)
}
