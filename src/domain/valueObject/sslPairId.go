package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

const sslPairIdExpression = "^[a-fA-F0-9]{64}$"

type SslPairId string

func NewSslPairId(value any) (sslPairId SslPairId, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
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
	var sslChainCertificatesMerged strings.Builder
	for _, sslChainCertificate := range sslChainCertificates {
		sslChainCertificatesMerged.WriteString(sslChainCertificate.String())
		sslChainCertificatesMerged.WriteString("\n")
	}
	contentToEncode := sslCertificate.String() + "\n" + sslChainCertificatesMerged.String() + "\n" + sslPrivateKey.String()

	sslPairIdContent := voHelper.StrongStringHasher(contentToEncode)
	return NewSslPairId(sslPairIdContent)
}

func (vo SslPairId) String() string {
	return string(vo)
}
