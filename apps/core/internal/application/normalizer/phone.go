package normalizer

import "unicode"

func Phone(value string) string {
	digits := make([]rune, 0, len(value))
	for _, r := range value {
		if unicode.IsDigit(r) {
			digits = append(digits, r)
		}
	}

	return string(digits)
}
