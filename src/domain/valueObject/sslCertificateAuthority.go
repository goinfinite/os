package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const sslCertificateAuthorityRegex string = `^\w{1,3}[\w\.\,\'\(\)\ \-]{0,100}$`

type SslCertificateAuthority string

func NewSslCertificateAuthority(value interface{}) (
	certificateAuthority SslCertificateAuthority, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return certificateAuthority, errors.New("SslCertificateAuthorityMustBeString")
	}

	re := regexp.MustCompile(sslCertificateAuthorityRegex)
	if !re.MatchString(stringValue) {
		return certificateAuthority, errors.New("InvalidSslCertificateAuthority")
	}

	return SslCertificateAuthority(stringValue), nil
}

func (vo SslCertificateAuthority) String() string {
	return string(vo)
}

func (vo SslCertificateAuthority) IsSelfSigned() bool {
	return string(vo) == "Self-signed"
}
