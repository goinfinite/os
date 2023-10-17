package valueObject

import (
	"encoding/hex"
	"errors"
	"regexp"

	"golang.org/x/crypto/sha3"
)

type SslHashId string

func newSslHashId(contentToEncode string) (SslHashId, error) {
	hash := sha3.New256()
	_, err := hash.Write([]byte(contentToEncode))
	if err != nil {
		return "", nil
	}
	sslHashIdBytes := hash.Sum(nil)
	sslHashIdStr := hex.EncodeToString(sslHashIdBytes)

	return SslHashId(sslHashIdStr), nil
}

func (sslHashId SslHashId) isValid() bool {
	sha256RegexExpression := "^[a-fA-F0-9]{64}$"
	sha256RegexRegex := regexp.MustCompile(sha256RegexExpression)
	return sha256RegexRegex.MatchString(string(sslHashId))
}

func NewSslHashIdFromSslPair(
	sslCertificate SslCertificateStr,
	sslChainCertificates []SslCertificateStr,
	sslPrivateKey SslPrivateKey,
) (SslHashId, error) {
	contentToEncode := sslCertificate.String()
	for _, sslChainCertificate := range sslChainCertificates {
		contentToEncode += sslChainCertificate.String()
	}
	contentToEncode += "\n" + sslPrivateKey.String()
	sslHashId, err := newSslHashId(contentToEncode)
	if err != nil || !sslHashId.isValid() {
		return "", errors.New("InvalidSslPairId")
	}
	return sslHashId, nil
}

func NewSslHashIdFromSslCertificate(
	sslCertificate SslCertificateStr,
) (SslHashId, error) {
	sslHashId, err := newSslHashId(sslCertificate.String())
	if err != nil || !sslHashId.isValid() {
		return "", errors.New("InvalidSslPairId")
	}
	return sslHashId, nil
}

func NewSslHashIdFromSslCertificatePanic(
	sslCertificate SslCertificateStr,
) SslHashId {
	sslHashId, err := NewSslHashIdFromSslCertificate(sslCertificate)
	if err != nil {
		panic(err)
	}
	return sslHashId
}

func (sslHashId SslHashId) String() string {
	return string(sslHashId)
}

func NewSslHashIdFromString(sslHashIdStr string) (SslHashId, error) {
	sslHashId := SslHashId(sslHashIdStr)
	if !sslHashId.isValid() {
		return "", errors.New("InvalidSslPairId")
	}
	return sslHashId, nil
}

func NewSslHashIdFromStringPanic(sslHashIdStr string) SslHashId {
	sslHashId, err := NewSslHashIdFromString(sslHashIdStr)
	if err != nil {
		panic(err)
	}
	return sslHashId
}
