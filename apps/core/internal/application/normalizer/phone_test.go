package normalizer

import "testing"

func TestPhoneKeepsOnlyDigits(t *testing.T) {
	tests := map[string]string{
		"+55 (11) 99999-9999": "5511999999999",
		"5511999999999":       "5511999999999",
		"55 11 99999-9999":    "5511999999999",
	}

	for input, expected := range tests {
		if got := Phone(input); got != expected {
			t.Fatalf("expected %q, got %q", expected, got)
		}
	}
}
