package valueObject

import "testing"

func TestEncodedContent(t *testing.T) {
	t.Run("ValidEncodedContent", func(t *testing.T) {
		validEncodedContents := []interface{}{
			"TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4gUHJvaW4gdG9ydG9yIG1hZ25hLCBiaWJlbmR1bSBpbiBtYWduYSB2aXRhZSwgc2FnaXR0aXMgZmVybWVudHVtIGp1c3RvLiBGdXNjZSBldCBuaWJoIHZ1bHB1dGF0ZSwgY29uZ3VlIGlwc3VtIGF0LCBjb252YWxsaXMgYW50ZS4=",
			"U2VkIGhlbmRyZXJpdCBuZWMgbnVsbGEgdmVsIGFjY3Vtc2FuLiBOdW5jIGxlY3R1cyBkdWksIHNvZGFsZXMgdXQgb3JuYXJlIHNlZCwgcG9zdWVyZSBub24gZWxpdC4gVml2YW11cyBzZWQgcHVydXMgc3VzY2lwaXQsIHBoYXJldHJhIG9kaW8gZXUsIHBvc3VlcmUgbWkuIFBoYXNlbGx1cyB1bGxhbWNvcnBlciBtYWxlc3VhZGEgcmlzdXMsIHV0IGFjY3Vtc2FuIHNhcGllbiBsYW9yZWV0IHZpdGFlLiBOdWxsYSB1dCBsaWd1bGEganVzdG8u",
			"VXQgZmluaWJ1cyBmZWxpcyBlZ2V0IG51bmMgY29tbW9kbywgZXUgY3Vyc3VzIHF1YW0gdWxsYW1jb3JwZXIuIFV0IHNlbXBlciBsZW8gaWQgb2RpbyBpbnRlcmR1bSwgY29uc2VjdGV0dXIgcHVsdmluYXIgbmVxdWUgdHJpc3RpcXVlLiBEdWlzIGZhdWNpYnVzIG1hZ25hIGV1IHF1YW0gc2FnaXR0aXMgcG9zdWVyZS4gQ2xhc3MgYXB0ZW50IHRhY2l0aSBzb2Npb3NxdSBhZCBsaXRvcmEgdG9ycXVlbnQgcGVyIGNvbnViaWEgbm9zdHJhLCBwZXIgaW5jZXB0b3MgaGltZW5hZW9zLiBQcm9pbiBpZCBlbmltIGVnZXQgbG9yZW0gdmVzdGlidWx1bSBpYWN1bGlzLiBTZWQgYWMgdG9ydG9yIHRvcnRvci4gU2VkIG5lYyBzYXBpZW4gc2l0IGFtZXQgb3JjaSBhY2N1bXNhbiBoZW5kcmVyaXQgbmVjIHZpdGFlIGVzdC4gUGVsbGVudGVzcXVlIHB1bHZpbmFyIGlwc3VtIHRvcnRvciwgbm9uIGN1cnN1cyBlbmltIHB1bHZpbmFyIHV0LiBGdXNjZSBldSBqdXN0byB0aW5jaWR1bnQsIGNvbnZhbGxpcyB1cm5hIG5vbiwgZXVpc21vZCBhcmN1Lg==",
		}

		for _, encodedContent := range validEncodedContents {
			_, err := NewEncodedContent(encodedContent)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", encodedContent, err.Error())
			}
		}
	})

	t.Run("InvalidEncodedContent", func(t *testing.T) {
		invalidEncodedContents := []interface{}{
			"", "ab123123sadasdbbb", "asjklfjskldgnsdmfnmxncsahoidjwqiejqelk",
		}

		for _, encodedContent := range invalidEncodedContents {
			_, err := NewEncodedContent(encodedContent)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", encodedContent)
			}
		}
	})
}
