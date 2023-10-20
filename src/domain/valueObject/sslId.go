package valueObject

import (
	"encoding/hex"
	"errors"
	"regexp"

	"golang.org/x/crypto/sha3"
)

const sslIdExpression = "^[a-fA-F0-9]{64}$"

type SslId string

func NewSslId(value string) (SslId, error) {
	sslId := SslId(value)
	if !sslId.isValid() {
		return "", errors.New("InvalidSslId")
	}

	return sslId, nil
}

func NewSslIdPanic(value string) SslId {
	sslId, err := NewSslId(value)
	if err != nil {
		panic(err)
	}

	return sslId
}

func (sslId SslId) isValid() bool {
	sslIdRegex := regexp.MustCompile(sslIdExpression)
	return sslIdRegex.MatchString(string(sslId))
}

func sslIdFactory(value string) (SslId, error) {
	hash := sha3.New256()
	_, err := hash.Write([]byte(value))
	if err != nil {
		return "", errors.New("InvalidSslId")
	}
	sslIdBytes := hash.Sum(nil)
	sslIdStr := hex.EncodeToString(sslIdBytes)

	return NewSslId(sslIdStr)
}

func NewSslIdFromSslPairContent(
	sslCertificate SslCertificateContent,
	sslChainCertificates []SslCertificateContent,
	sslPrivateKey SslPrivateKey,
) (SslId, error) {
	var sslChainCertificatesMerged string
	for _, sslChainCertificate := range sslChainCertificates {
		sslChainCertificatesMerged += sslChainCertificate.String() + "\n"
	}

	contentToEncode := sslCertificate.String() + "\n" + sslChainCertificatesMerged + "\n" + sslPrivateKey.String()

	sslId, err := sslIdFactory(contentToEncode)
	if err != nil {
		return "", errors.New("InvalidSslIdFromSslPairContent")
	}

	return sslId, nil
}

func NewSslIdFromSslCertificateContent(
	sslCertificate SslCertificateContent,
) (SslId, error) {
	sslId, err := sslIdFactory(sslCertificate.String())
	if err != nil {
		return "", errors.New("InvalidSslIdFromSslCertificateContent")
	}

	return sslId, nil
}

func (sslId SslId) String() string {
	return string(sslId)
}
