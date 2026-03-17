package valueObject

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type SslPrivateKey tkValueObject.EnvelopedPrivateKey

func NewSslPrivateKey(
	value interface{},
) (privateKey SslPrivateKey, err error) {
	envelopedKey, err := tkValueObject.NewEnvelopedPrivateKey(
		value,
	)
	if err != nil {
		return privateKey, errors.New(
			"InvalidSslPrivateKey",
		)
	}

	return SslPrivateKey(envelopedKey), nil
}

func (vo SslPrivateKey) String() string {
	return string(vo)
}

// UnmarshalJSON is intentionally absent. All input paths construct SslPrivateKey
// via NewSslPrivateKey() explicitly — no code path unmarshals it from JSON.
func (vo SslPrivateKey) MarshalJSON() ([]byte, error) {
	encodedContent := base64.StdEncoding.EncodeToString(
		[]byte(string(vo)),
	)
	return json.Marshal(encodedContent)
}
