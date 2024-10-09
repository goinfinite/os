package valueObject

import (
	"encoding/hex"
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	"golang.org/x/crypto/sha3"
)

const sslIdExpression = "^[a-fA-F0-9]{64}$"

type SslId string

func NewSslId(value interface{}) (sslId SslId, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return sslId, errors.New("SslIdMustBeString")
	}

	re := regexp.MustCompile(sslIdExpression)
	if !re.MatchString(stringValue) {
		return sslId, errors.New("InvalidSslId")
	}

	return SslId(stringValue), nil
}

func sslIdFactory(sslCertContent *string) (sslId SslId, err error) {
	hash := sha3.New256()
	_, err = hash.Write([]byte(*sslCertContent))
	if err != nil {
		return sslId, errors.New("InvalidSslId")
	}
	sslIdBytes := hash.Sum(nil)
	sslIdStr := hex.EncodeToString(sslIdBytes)

	return NewSslId(sslIdStr)
}

func NewSslIdFromSslPairContent(
	sslCertificate SslCertificateContent,
	sslChainCertificates []SslCertificateContent,
	sslPrivateKey SslPrivateKey,
) (sslId SslId, err error) {
	sslChainCertificatesMerged := ""
	for _, sslChainCertificate := range sslChainCertificates {
		sslChainCertificatesMerged += sslChainCertificate.String() + "\n"
	}
	contentToEncode := sslCertificate.String() + "\n" + sslChainCertificatesMerged + "\n" + sslPrivateKey.String()

	sslId, err = sslIdFactory(&contentToEncode)
	if err != nil {
		return sslId, errors.New("InvalidSslIdFromSslPairContent")
	}

	return sslId, nil
}

func NewSslIdFromSslCertificateContent(
	sslCertificate SslCertificateContent,
) (sslId SslId, err error) {
	certContentStr := sslCertificate.String()
	sslId, err = sslIdFactory(&certContentStr)
	if err != nil {
		return sslId, errors.New("InvalidSslIdFromSslCertificateContent")
	}

	return sslId, nil
}

func (vo SslId) String() string {
	return string(vo)
}
