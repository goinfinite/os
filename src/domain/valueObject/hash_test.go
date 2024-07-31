package valueObject

import "testing"

func TestHash(t *testing.T) {
	t.Run("ValidHash (NTLM)", func(t *testing.T) {
		ntlmValidHash := "84412CB56723EE2B08D680D78D0D8A46"
		_, err := NewHash(ntlmValidHash)
		if err != nil {
			t.Errorf("Expected no error for NTLM hash '%s', got '%s'", ntlmValidHash, err.Error())
		}
	})

	t.Run("ValidHash (MD5)", func(t *testing.T) {
		md5ValidHash := "0340b99427817c02f343941e984c9bb2"
		_, err := NewHash(md5ValidHash)
		if err != nil {
			t.Errorf("Expected no error for MD5 hash '%s', got '%s'", md5ValidHash, err.Error())
		}
	})

	t.Run("ValidHash (SHA1)", func(t *testing.T) {
		sha1ValidHash := "d8b5fdce85438fb7cb4b343b6d9c812e351e253d"
		_, err := NewHash(sha1ValidHash)
		if err != nil {
			t.Errorf("Expected no error for SHA1 hash '%s', got '%s'", sha1ValidHash, err.Error())
		}
	})

	t.Run("ValidHash (SHA256)", func(t *testing.T) {
		sha256ValidHash := "7df07c76c7dbf395c07c224959bc19a62e84cc5f9db24398523f4961beedb673"
		_, err := NewHash(sha256ValidHash)
		if err != nil {
			t.Errorf("Expected no error for SHA256 hash '%s', got '%s'", sha256ValidHash, err.Error())
		}
	})

	t.Run("ValidHash (SHA512)", func(t *testing.T) {
		sha512ValidHash := "341666654e2340d961374f0fcfb67e06e752b7b2ca34186163745fc62f82c0470b507e087a4ffee8dabc576b779ab43efada95b88024fdddb393cef09cf7e1fb"
		_, err := NewHash(sha512ValidHash)
		if err != nil {
			t.Errorf("Expected no error for SHA512 hash '%s', got '%s'", sha512ValidHash, err.Error())
		}
	})

	t.Run("InvalidHash", func(t *testing.T) {
		invalidHashes := []interface{}{
			"",
			"     ",
			"Kf81h",
			"TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4gQ3JhcyBhbGlxdWV0IGRpYW0gaWQgcGxhY2VyYXQgZGFwaWJ1cy4gQ3VyYWJpdHVyIGVsZWlmZW5kIG1hdHRpcyB1cm5hIG5vbiB2dWxwdXRhdGUuIFN1c3BlbmRpc3NlIHBvdGVudGkuIE51bmMgZGlnbmlzc2ltIG5pc2wgdml0YWUgbnVsb",
		}
		for _, invalidHash := range invalidHashes {
			_, err := NewHash(invalidHash)
			if err == nil {
				t.Errorf("Expected error for '%s', got nil", invalidHash)
			}
		}
	})
}
