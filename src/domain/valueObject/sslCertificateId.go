package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var sslCertificateIdRegex = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)

type SslCertificateId string

func NewSslCertificateId(value any) (sslCertificateId SslCertificateId, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return sslCertificateId, errors.New("SslCertificateIdMustBeString")
	}

	if !sslCertificateIdRegex.MatchString(stringValue) {
		return sslCertificateId, errors.New("InvalidSslCertificateId")
	}

	return SslCertificateId(stringValue), nil
}

func NewSslCertificateIdFromSslCertificateContent(
	sslCertificate SslCertificateContent,
) (sslCertificateId SslCertificateId, err error) {
	sslCertificateIdContent := voHelper.StrongStringHasher(sslCertificate.String())
	return NewSslCertificateId(sslCertificateIdContent)
}

func (vo SslCertificateId) String() string {
	return string(vo)
}
