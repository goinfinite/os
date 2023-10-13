package valueObject

import (
	"math/big"
	"math/rand"
	"testing"
)

func TestNewSslSerialNumber(t *testing.T) {
	t.Run("ValidSerialNumber", func(t *testing.T) {
		validSerialNumbers := []interface{}{
			genDummyBigInt(false),
			genDummyBigInt(false).String(),
		}

		for _, sslSerialNumber := range validSerialNumbers {
			_, err := NewSslSerialNumber(sslSerialNumber)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", sslSerialNumber, err)
			}
		}
	})

	t.Run("InvalidSerialNumber", func(t *testing.T) {
		invalidSerialNumbers := []interface{}{
			genDummyBigInt(true),
			genDummyBigInt(true).String(),
		}

		for _, sslSerialNumber := range invalidSerialNumbers {
			_, err := NewSslSerialNumber(sslSerialNumber)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", sslSerialNumber)
			}
		}
	})
}

func genDummyBigInt(negative bool) *big.Int {
	randomIntNumber := rand.Int63n(9999999)
	if negative {
		randomIntNumber -= randomIntNumber * 2
	}
	return new(big.Int).SetInt64(randomIntNumber)
}
