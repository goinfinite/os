package valueObject

import (
	"errors"
	"regexp"
)

const sslCertificateAuthorityRegex string = `^\w{1,3}[\w\.\,\'\(\)\ ]{0,100}$`

type SslCertificateAuthority string

func NewSslCertificateAuthority(value string) (SslCertificateAuthority, error) {
	sca := SslCertificateAuthority(value)

	if !sca.isValid() {
		return "", errors.New("InvalidSslCertificateAuthority")
	}

	return SslCertificateAuthority(sca), nil
}

func NewSslCertificateAuthorityPanic(value string) SslCertificateAuthority {
	sn, err := NewSslCertificateAuthority(value)
	if err != nil {
		panic(err)
	}
	return sn
}

func (sca SslCertificateAuthority) isValid() bool {
	compiledScaRegex := regexp.MustCompile(sslCertificateAuthorityRegex)
	return compiledScaRegex.MatchString(string(sca))
}

func (sca SslCertificateAuthority) String() string {
	return string(sca)
}
