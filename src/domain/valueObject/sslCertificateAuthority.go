package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var sslCertificateAuthorityRegex = regexp.MustCompile(`^\w{1,3}[\w\.\,\'\(\)\ \-]{0,100}$`)

type SslCertificateAuthority string

func NewSslCertificateAuthority(value any) (
	certificateAuthority SslCertificateAuthority, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return certificateAuthority, errors.New("SslCertificateAuthorityMustBeString")
	}

	if !sslCertificateAuthorityRegex.MatchString(stringValue) {
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
