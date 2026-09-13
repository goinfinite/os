package voHelper

import "testing"

func TestStrongStringHasher(t *testing.T) {
	t.Run("MatchesFips202Sha3_256Vector", func(t *testing.T) {
		hashed := StrongStringHasher("abc")

		expectedHash := "3a985da74fe225b2045c172d6bd390bd855f086e3e9d525b46bfe24511431532"
		if hashed != expectedHash {
			t.Errorf("DigestMismatch: expected '%s', got '%s'", expectedHash, hashed)
		}
	})

	t.Run("MatchesFips202Sha3_256EmptyStringVector", func(t *testing.T) {
		hashed := StrongStringHasher("")

		expectedHash := "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"
		if hashed != expectedHash {
			t.Errorf("DigestMismatch: expected '%s', got '%s'", expectedHash, hashed)
		}
	})
}

func TestStrongStringShortHasher(t *testing.T) {
	t.Run("TruncatesDigestToShortHashLength", func(t *testing.T) {
		shortHash := StrongStringShortHasher("abc")

		expectedShortHash := "3a985da74fe2"
		if shortHash != expectedShortHash {
			t.Errorf(
				"ShortDigestMismatch: expected '%s', got '%s'",
				expectedShortHash,
				shortHash,
			)
		}

		if len(shortHash) != shortHashLength {
			t.Errorf(
				"ShortDigestLengthMismatch: expected '%d', got '%d'",
				shortHashLength,
				len(shortHash),
			)
		}
	})

	t.Run("TruncatesEmptyStringDigest", func(t *testing.T) {
		shortHash := StrongStringShortHasher("")

		expectedShortHash := "a7ffc6f8bf1e"
		if shortHash != expectedShortHash {
			t.Errorf(
				"ShortDigestMismatch: expected '%s', got '%s'",
				expectedShortHash,
				shortHash,
			)
		}
	})
}
