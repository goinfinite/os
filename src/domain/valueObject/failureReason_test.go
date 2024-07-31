package valueObject

import (
	"testing"
)

func TestFailureReason(t *testing.T) {
	t.Run("ValidFailureReason", func(t *testing.T) {
		validFailureReasons := []interface{}{
			"InvalidRecordId",
			"Container must be primary",
			"Your currently vhost is not able to get a alias",
			"This user should not be able to update API required policies",
		}

		for _, reason := range validFailureReasons {
			_, err := NewFailureReason(reason)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", reason, err)
			}
		}
	})

	t.Run("InvalidFailureReason", func(t *testing.T) {
		invalidFailureReasons := []interface{}{
			"",
		}

		for _, reason := range invalidFailureReasons {
			_, err := NewFailureReason(reason)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", reason)
			}
		}
	})
}
