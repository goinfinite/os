package valueObject

import (
	"encoding/hex"
	"errors"
	"regexp"

	"golang.org/x/crypto/sha3"
)

const sslHashIdExpression = "^[a-fA-F0-9]{64}$"

type SslHashId string

func NewSslHashId(value string) (SslHashId, error) {
	sslHashId := SslHashId(value)
	if !sslHashId.isValid() {
		return "", errors.New("InvalidSslHashId")
	}

	return sslHashId, nil
}

func NewSslHashIdPanic(value string) SslHashId {
	sslHashId, err := NewSslHashId(value)
	if err != nil {
		panic(err)
	}

	return sslHashId
}

func (sslHashId SslHashId) isValid() bool {
	sslHashIdRegex := regexp.MustCompile(sslHashIdExpression)
	return sslHashIdRegex.MatchString(string(sslHashId))
}

func sslHashIdFactory(value string) (SslHashId, error) {
	hash := sha3.New256()
	_, err := hash.Write([]byte(value))
	if err != nil {
		return "", errors.New("InvalidSslHashId")
	}
	sslHashIdBytes := hash.Sum(nil)
	sslHashIdStr := hex.EncodeToString(sslHashIdBytes)

	return NewSslHashId(sslHashIdStr)
}

func NewSslHashIdFromSslPairContent(
	sslCertificate SslCertificateStr,
	sslChainCertificates []SslCertificateStr,
	sslPrivateKey SslPrivateKey,
) (SslHashId, error) {
	var sslChainCertificatesMerged string
	for _, sslChainCertificate := range sslChainCertificates {
		sslChainCertificatesMerged += sslChainCertificate.String() + "\n"
	}

	contentToEncode := sslCertificate.String() + "\n" + sslChainCertificatesMerged + "\n" + sslPrivateKey.String()

	sslHashId, err := sslHashIdFactory(contentToEncode)
	if err != nil {
		return "", errors.New("InvalidSslHashIdFromSslPairContent")
	}

	return sslHashId, nil
}

func NewSslHashIdFromSslCertificateContent(
	sslCertificate SslCertificateStr,
) (SslHashId, error) {
	sslHashId, err := sslHashIdFactory(sslCertificate.String())
	if err != nil {
		return "", errors.New("InvalidSslHashIdFromSslCertificateContent")
	}

	return sslHashId, nil
}

func (sslHashId SslHashId) String() string {
	return string(sslHashId)
}
