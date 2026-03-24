package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

const sslCertificateIdExpression = "^[a-fA-F0-9]{64}$"

type SslCertificateId string

func NewSslCertificateId(value interface{}) (sslCertificateId SslCertificateId, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return sslCertificateId, errors.New("SslCertificateIdMustBeString")
	}

	re := regexp.MustCompile(sslCertificateIdExpression)
	if !re.MatchString(stringValue) {
		return sslCertificateId, errors.New("InvalidSslCertificateId")
	}

	return SslCertificateId(stringValue), nil
}

func NewSslCertificateIdFromSslCertificateContent(
	sslCertificate SslCertificateContent,
) (sslCertificateId SslCertificateId, err error) {
	sslCertificateIdContent, err := voHelper.StrongStringHasher(
		sslCertificate.String(),
	)
	if err != nil {
		return sslCertificateId, errors.New(
			"InvalidSslCertificateIdFromSslCertificateContent",
		)
	}
	return NewSslCertificateId(sslCertificateIdContent)
}

func (vo SslCertificateId) String() string {
	return string(vo)
}
