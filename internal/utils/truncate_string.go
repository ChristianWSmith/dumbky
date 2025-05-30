package utils

func TruncateStringWithEllipsis(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	r := []rune(s)
	if len(r) <= maxLength {
		return s
	}

	return string(r[:maxLength]) + "..."
}
