package valueObject

import "testing"

func TestMappingId(t *testing.T) {
	t.Run("ValidMappingId", func(t *testing.T) {
		validMappingIds := []interface{}{
			0, 1, 10000000000000, "455", 40.5,
		}

		for _, mappingId := range validMappingIds {
			_, err := NewMappingId(mappingId)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", mappingId, err.Error())
			}
		}
	})

	t.Run("InvalidMappingId", func(t *testing.T) {
		invalidMappingIds := []interface{}{
			-1, -10000000000000, "-455", -40.5,
		}

		for _, mappingId := range invalidMappingIds {
			_, err := NewMappingId(mappingId)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", mappingId)
			}
		}
	})
}
